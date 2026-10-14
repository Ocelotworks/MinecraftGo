// Package launcher reads the Minecraft launcher's install: the version
// manifests under versions/, the library jars they name, and the JREs they ask
// for.  It is what lets the tools under cmd/ run the game's own code — the data
// generator for mcdata, the client itself for mcclient — without shipping a
// copy of any of it.
package launcher

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// Rule gates a library or a launch argument on the platform, or on a feature
// the launcher says it is using.
//
// os.version is a regex against the operating system version.  Nothing in a
// current manifest uses it, and reading the Windows version needs a syscall, so
// it is ignored here: a rule naming one is judged on its other fields.
type Rule struct {
	Action string `json:"action"`
	OS     struct {
		Name    string `json:"name"`
		Arch    string `json:"arch"`
		Version string `json:"version"`
	} `json:"os"`
	Features map[string]bool `json:"features"`
}

// Features names the launcher features a rule can ask about, such as
// is_quick_play_multiplayer.  A feature missing from the map counts as off.
type Features map[string]bool

// currentOS is the name Mojang's manifests use for this platform.
func currentOS() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "darwin":
		return "osx"
	default:
		return runtime.GOOS
	}
}

// currentArch is the name Mojang's manifests use for this architecture.
func currentArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "386":
		return "x86"
	default:
		return runtime.GOARCH
	}
}

// matches reports whether a rule applies to this platform and feature set.
func (r Rule) matches(features Features) bool {
	if r.OS.Name != "" && r.OS.Name != currentOS() {
		return false
	}
	if r.OS.Arch != "" && r.OS.Arch != currentArch() {
		return false
	}
	for feature, wanted := range r.Features {
		if features[feature] != wanted {
			return false
		}
	}
	return true
}

// allowed applies a rule list the way the launcher does: nothing is allowed
// unless a rule says so, and a later matching rule overrides an earlier one.
// An empty list allows.
func allowed(rules []Rule, features Features) bool {
	if len(rules) == 0 {
		return true
	}

	result := false
	for _, rule := range rules {
		if rule.matches(features) {
			result = rule.Action == "allow"
		}
	}
	return result
}

// Argument is one entry in a manifest's argument list.  It is either a bare
// string or an object carrying rules and one or more values.
type Argument struct {
	Rules  []Rule
	Values []string
}

func (a *Argument) UnmarshalJSON(raw []byte) error {
	plain := ""
	if err := json.Unmarshal(raw, &plain); err == nil {
		a.Rules, a.Values = nil, []string{plain}
		return nil
	}

	conditional := struct {
		Rules []Rule          `json:"rules"`
		Value json.RawMessage `json:"value"`
	}{}
	if err := json.Unmarshal(raw, &conditional); err != nil {
		return fmt.Errorf("parsing launch argument %s: %w", raw, err)
	}

	values := []string{}
	if err := json.Unmarshal(conditional.Value, &values); err != nil {
		single := ""
		if err := json.Unmarshal(conditional.Value, &single); err != nil {
			return fmt.Errorf("parsing launch argument value %s: %w", conditional.Value, err)
		}
		values = []string{single}
	}

	a.Rules, a.Values = conditional.Rules, values
	return nil
}

// Library is one entry in a manifest's library list.
type Library struct {
	Name      string `json:"name"`
	Rules     []Rule `json:"rules"`
	Downloads struct {
		Artifact struct {
			Path string `json:"path"`
		} `json:"artifact"`
	} `json:"downloads"`
}

// VersionManifest is the subset of a launcher version manifest these tools
// need: what to put on the classpath, which JRE to use, and — for the client —
// the main class and the argument templates to fill in.
type VersionManifest struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	MainClass   string `json:"mainClass"`
	JavaVersion struct {
		Component    string `json:"component"`
		MajorVersion int    `json:"majorVersion"`
	} `json:"javaVersion"`
	AssetIndex struct {
		ID string `json:"id"`
	} `json:"assetIndex"`
	Arguments struct {
		JVM  []Argument `json:"jvm"`
		Game []Argument `json:"game"`
	} `json:"arguments"`
	Libraries []Library `json:"libraries"`
}

// FindMinecraftDir locates the launcher directory holding versions/ and
// libraries/.
func FindMinecraftDir() (string, error) {
	candidates := make([]string, 0, 3)

	if appData := os.Getenv("APPDATA"); appData != "" {
		candidates = append(candidates, filepath.Join(appData, ".minecraft"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			filepath.Join(home, ".minecraft"),
			filepath.Join(home, "Library", "Application Support", "minecraft"))
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(filepath.Join(candidate, "versions")); err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("could not find a Minecraft install (looked in %v); pass -minecraft", candidates)
}

// ReadVersionManifest loads the manifest for a version id, falling back to a
// search when the directory is not named after the version.
func ReadVersionManifest(minecraftDir string, versionID string) (VersionManifest, error) {
	manifest := VersionManifest{}

	direct := filepath.Join(minecraftDir, "versions", versionID, versionID+".json")
	if raw, err := os.ReadFile(direct); err == nil {
		if err := json.Unmarshal(raw, &manifest); err != nil {
			return manifest, fmt.Errorf("parsing %s: %w", direct, err)
		}
		return manifest, nil
	}

	entries, err := os.ReadDir(filepath.Join(minecraftDir, "versions"))
	if err != nil {
		return manifest, fmt.Errorf("listing versions in %s: %w", minecraftDir, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(minecraftDir, "versions", entry.Name(), entry.Name()+".json")
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		candidate := VersionManifest{}
		if json.Unmarshal(raw, &candidate) == nil && candidate.ID == versionID {
			return candidate, nil
		}
	}

	return manifest, fmt.Errorf("no manifest for version %q under %s; run the launcher once for it, or pass -version", versionID, minecraftDir)
}

// ClientJar is the path the launcher keeps a version's client jar at.
func (m VersionManifest) ClientJar(minecraftDir string) string {
	return filepath.Join(minecraftDir, "versions", m.ID, m.ID+".jar")
}

// Classpath collects the library jars a version needs on this platform.
//
// Libraries the launcher has not downloaded are skipped rather than treated as
// missing: a manifest lists the natives for every platform, and only this one's
// are on disc.
//
// Since 1.19 those natives are ordinary classpath entries rather than a
// directory to unpack, which is why nothing here extracts anything.  LWJGL and
// Netty pull their own libraries out of these jars at startup, into whichever
// directory the manifest's -D arguments point them at.
func (m VersionManifest) Classpath(minecraftDir string, features Features) (paths []string, skipped int) {
	for _, library := range m.Libraries {
		if library.Downloads.Artifact.Path == "" || !allowed(library.Rules, features) {
			continue
		}
		path := filepath.Join(minecraftDir, "libraries", filepath.FromSlash(library.Downloads.Artifact.Path))
		if _, err := os.Stat(path); err != nil {
			skipped++
			continue
		}
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, skipped
}

// ResolveArguments fills in a manifest argument list: it drops the entries this
// platform and feature set do not allow, and expands the ${...} placeholders
// from values.  A placeholder with no value is an error rather than something
// passed through, since the game would read the literal text as the setting.
func ResolveArguments(arguments []Argument, features Features, values map[string]string) ([]string, error) {
	resolved := make([]string, 0, len(arguments))

	for _, argument := range arguments {
		if !allowed(argument.Rules, features) {
			continue
		}
		for _, value := range argument.Values {
			expanded, err := expand(value, values)
			if err != nil {
				return nil, err
			}
			resolved = append(resolved, expanded)
		}
	}

	return resolved, nil
}

// expand replaces every ${name} in an argument with its value.
func expand(argument string, values map[string]string) (string, error) {
	for cursor := 0; ; {
		start := strings.Index(argument[cursor:], "${")
		if start < 0 {
			return argument, nil
		}
		start += cursor

		end := strings.Index(argument[start:], "}")
		if end < 0 {
			return argument, nil
		}
		end += start

		name := argument[start+2 : end]
		value, ok := values[name]
		if !ok {
			return "", fmt.Errorf("launch argument %q wants ${%s}, which this tool does not set", argument, name)
		}

		argument = argument[:start] + value + argument[end+1:]
		cursor = start + len(value)
	}
}

// FindJava locates a JRE able to run a version.  The launcher installs one per
// component under a runtime directory, laid out as
// <runtime>/<component>/<platform>/<component>/bin/java.
func FindJava(minecraftDir string, component string) (string, error) {
	executable := "java"
	if runtime.GOOS == "windows" {
		executable = "java.exe"
	}

	roots := []string{filepath.Join(minecraftDir, "runtime")}
	if runtime.GOOS == "windows" {
		for _, programFiles := range []string{os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles")} {
			if programFiles != "" {
				roots = append(roots, filepath.Join(programFiles, "Minecraft Launcher", "runtime"))
			}
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			roots = append(roots, filepath.Join(localAppData,
				"Packages", "Microsoft.4297127D64EC6_8wekyb3d8bbwe", "LocalCache", "Local", "runtime"))
		}
	}

	for _, root := range roots {
		matches, err := filepath.Glob(filepath.Join(root, component, "*", component, "bin", executable))
		if err != nil || len(matches) == 0 {
			continue
		}
		return matches[0], nil
	}

	return "", fmt.Errorf("could not find the %s runtime (looked under %v); pass -java with a JDK new enough for this version", component, roots)
}
