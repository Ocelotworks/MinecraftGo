package main

import (
	"encoding/json"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Ocelotworks/MinecraftGo/cmd/internal/launcher"
	"github.com/Ocelotworks/MinecraftGo/entity"
)

// reportBlock is one block in the generator's blocks.json report.  The report
// lists every state with the id the client will use for it, which is the whole
// reason for running the generator: those ids come from registry order at
// runtime and appear nowhere in the jar.
type reportBlock struct {
	Properties map[string][]string `json:"properties"`
	States     []struct {
		ID         int               `json:"id"`
		Default    bool              `json:"default"`
		Properties map[string]string `json:"properties"`
	} `json:"states"`
}

// registriesReport is the generator's registries.json, used here for the block
// registry's protocol ids.
type registriesReport map[string]struct {
	Entries map[string]struct {
		ProtocolID int `json:"protocol_id"`
	} `json:"entries"`
}

// generateFromReports runs the game's data generator once and converts the two
// reports the server depends on: the block state table and the packet ids.
func generateFromReports(jar *jarFile, version jarVersion, opts options) error {
	minecraftDir := opts.minecraftPath
	if minecraftDir == "" {
		found, err := launcher.FindMinecraftDir()
		if err != nil {
			return err
		}
		minecraftDir = found
	}

	versionID := opts.versionID
	if versionID == "" {
		versionID = version.ID
	}

	manifest, err := launcher.ReadVersionManifest(minecraftDir, versionID)
	if err != nil {
		return err
	}

	javaPath := opts.javaPath
	if javaPath == "" {
		javaPath, err = launcher.FindJava(minecraftDir, manifest.JavaVersion.Component)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Using java %s (needs %d)\n", javaPath, manifest.JavaVersion.MajorVersion)

	libraries, skipped := manifest.Classpath(minecraftDir, nil)
	fmt.Printf("Classpath: %d libraries from %s (%d not installed for this platform)\n", len(libraries), versionID, skipped)

	reportsDir, cleanup, err := runDataGenerator(javaPath, jar.path, libraries, opts.keepReports)
	if err != nil {
		return err
	}
	defer cleanup()

	if !opts.skipPackets {
		if err := generatePackets(reportsDir, opts.outputPath); err != nil {
			return err
		}
	}

	if !opts.skipRegistryIDs {
		if err := generateRegistryIDs(reportsDir, opts.outputPath); err != nil {
			return err
		}
	}

	if opts.skipBlocks {
		return nil
	}

	blocks, err := convertBlockReports(reportsDir, jar, filepath.Join(opts.outputPath, "data", "blocks.json"))
	if err != nil {
		return err
	}

	target := filepath.Join(opts.outputPath, "data", "blocks.json")
	encoded, err := json.MarshalIndent(blocks, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(target, append(encoded, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", target, err)
	}

	maxStateID := 0
	for _, block := range blocks {
		if block.MaxStateId > maxStateID {
			maxStateID = block.MaxStateId
		}
	}
	fmt.Printf("Wrote %s: %d blocks, %d block states, %d bits for the global palette\n",
		target, len(blocks), maxStateID+1, bits.Len32(uint32(maxStateID)))

	return nil
}

// runDataGenerator invokes net.minecraft.data.Main --reports and returns the
// directory its reports landed in.
func runDataGenerator(javaPath string, jarPath string, libraries []string, keepDir string) (string, func(), error) {
	outputDir := keepDir
	cleanup := func() {}

	if outputDir == "" {
		temp, err := os.MkdirTemp("", "mcdata-*")
		if err != nil {
			return "", cleanup, err
		}
		outputDir = temp
		cleanup = func() { os.RemoveAll(temp) }
	} else if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", cleanup, err
	}

	absJar, err := filepath.Abs(jarPath)
	if err != nil {
		return "", cleanup, err
	}

	classpath := strings.Join(append(append([]string{}, libraries...), absJar), string(os.PathListSeparator))

	fmt.Println("Running the data generator, this takes a few seconds...")
	command := exec.Command(javaPath, "-cp", classpath, "net.minecraft.data.Main", "--reports", "--output", outputDir)
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("running the data generator: %w", err)
	}

	reportsDir := filepath.Join(outputDir, "reports")
	if _, err := os.Stat(filepath.Join(reportsDir, "blocks.json")); err != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("the data generator produced no blocks.json in %s", reportsDir)
	}

	return reportsDir, cleanup, nil
}

// convertBlockReports turns the generator's reports into the block table the
// server reads.
//
// The reports carry the registry ids and every state, and the jar carries the
// display names.  Gameplay values the reports do not describe — hardness,
// drops, harvest tools and the like — are carried over from the existing file
// where the block still exists, since nothing in the jar states them directly.
func convertBlockReports(reportsDir string, jar *jarFile, existingPath string) ([]entity.BlockData, error) {
	blockReport := map[string]reportBlock{}
	if err := readJSONFile(filepath.Join(reportsDir, "blocks.json"), &blockReport); err != nil {
		return nil, err
	}

	registries := registriesReport{}
	if err := readJSONFile(filepath.Join(reportsDir, "registries.json"), &registries); err != nil {
		return nil, err
	}
	blockRegistry := registries["minecraft:block"].Entries
	if len(blockRegistry) == 0 {
		return nil, fmt.Errorf("registries.json has no minecraft:block registry")
	}

	displayNames, err := jar.displayNames()
	if err != nil {
		return nil, err
	}

	existing, err := readExistingBlocks(existingPath)
	if err != nil {
		return nil, err
	}

	blocks := make([]entity.BlockData, 0, len(blockReport))
	carried, fresh := 0, 0

	for name, report := range blockReport {
		shortName := strings.TrimPrefix(name, "minecraft:")

		block, ok := existing[shortName]
		if ok {
			carried++
		} else {
			fresh++
			block = entity.BlockData{StackSize: 64, BoundingBox: "block"}
		}

		block.Name = shortName
		block.ID = blockRegistry[name].ProtocolID
		if display, ok := displayNames[name]; ok {
			block.DisplayName = display
		}
		block.States = convertProperties(report.Properties)

		block.MinStateId, block.MaxStateId, block.DefaultState = stateRange(report)

		if block.Drops == nil {
			block.Drops = []int{}
		}
		blocks = append(blocks, block)
	}

	// Sort by registry id so the file is stable between runs and diffs cleanly.
	sort.Slice(blocks, func(i, j int) bool { return blocks[i].ID < blocks[j].ID })

	fmt.Printf("Blocks: %d carried gameplay values over from the previous file, %d are new to this version\n", carried, fresh)
	if missing := len(existing) - carried; missing > 0 {
		fmt.Printf("Blocks: %d blocks in the previous file no longer exist and were dropped\n", missing)
	}
	if carried > 0 {
		fmt.Println("Blocks: hardness, resistance, harvest tools, light and drops are not in the reports and were")
		fmt.Println("        carried over as-is — drops in particular hold item ids that may have moved.")
	}

	return blocks, nil
}

// stateRange returns the id range a block's states occupy and its default.
func stateRange(report reportBlock) (minID int, maxID int, defaultID int) {
	minID, maxID, defaultID = -1, -1, -1

	for _, state := range report.States {
		if minID == -1 || state.ID < minID {
			minID = state.ID
		}
		if state.ID > maxID {
			maxID = state.ID
		}
		if state.Default {
			defaultID = state.ID
		}
	}

	if defaultID == -1 {
		defaultID = minID
	}
	return minID, maxID, defaultID
}

// convertProperties describes a block's properties in the order its states are
// numbered in.
//
// The game sorts a block's properties by name and counts through them with the
// last varying fastest, so the list has to come out sorted here — the report
// itself does not always emit them that way.
func convertProperties(properties map[string][]string) []entity.BlockState {
	names := make([]string, 0, len(properties))
	for name := range properties {
		names = append(names, name)
	}
	sort.Strings(names)

	states := make([]entity.BlockState, 0, len(names))
	for _, name := range names {
		values := properties[name]
		states = append(states, entity.BlockState{
			Name:      name,
			Type:      propertyType(values),
			NumValues: len(values),
			Values:    values,
		})
	}

	return states
}

func propertyType(values []string) string {
	if len(values) == 2 && values[0] == "true" && values[1] == "false" {
		return "bool"
	}
	for _, value := range values {
		if _, err := strconv.Atoi(value); err != nil {
			return "enum"
		}
	}
	return "int"
}

// readExistingBlocks loads the current block table, keyed by name.  A missing
// file is fine: everything is then generated from scratch, without the
// gameplay values the reports do not cover.
func readExistingBlocks(path string) (map[string]entity.BlockData, error) {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]entity.BlockData{}, nil
	}
	if err != nil {
		return nil, err
	}

	blocks := make([]entity.BlockData, 0)
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return nil, fmt.Errorf("parsing the existing %s: %w", path, err)
	}

	byName := make(map[string]entity.BlockData, len(blocks))
	for _, block := range blocks {
		byName[block.Name] = block
	}
	return byName, nil
}

func readJSONFile(path string, into any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, into); err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}
	return nil
}

// generateRegistryIDs writes data/registry_ids.json: the protocol id of every
// entry in every registry the generator reports.
//
// The client numbers registry entries the same way it numbers block states —
// at runtime, from registration order — so the ids appear nowhere in the jar.
// The server needs them to send tags: a tag binds a name to a list of entry ids,
// and the static registries (block, item, entity_type, ...) are never sent to
// the client, so nothing else on the server knows their numbering.
func generateRegistryIDs(reportsDir string, outputPath string) error {
	registries := registriesReport{}
	if err := readJSONFile(filepath.Join(reportsDir, "registries.json"), &registries); err != nil {
		return err
	}

	ids := make(map[string]map[string]int, len(registries))
	entries := 0
	for registryName, registry := range registries {
		ids[registryName] = make(map[string]int, len(registry.Entries))
		for entryName, entry := range registry.Entries {
			ids[registryName][entryName] = entry.ProtocolID
			entries++
		}
	}

	target := filepath.Join(outputPath, "data", "registry_ids.json")
	encoded, err := json.Marshal(ids) // map keys are sorted, so reruns produce no diff noise
	if err != nil {
		return err
	}
	if err := os.WriteFile(target, append(encoded, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", target, err)
	}

	fmt.Printf("Wrote %s: %d registries, %d entries\n", target, len(ids), entries)
	return nil
}
