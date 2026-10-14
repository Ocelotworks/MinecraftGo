// Command mcclient launches a real Minecraft client against this server and
// captures everything it says.
//
// It joins on startup, so there is no main menu to click through, and its
// output — the client's log plus anything the JVM writes over the top of it —
// goes to the terminal and to a file.  That log is the point of the tool: a
// protocol mistake shows up there as a disconnect reason or a stack trace,
// where on the server side it is just a connection that went away.
//
// The client, its libraries and its JRE all come from the Minecraft launcher's
// install, the same place cmd/mcdata gets them from.  Run the launcher once for
// the version you want first, so it is on disc.
//
// Usage:
//
//	go run ./cmd/mcclient
//	go run ./cmd/mcclient -server 127.0.0.1:25565 -username Dev
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Ocelotworks/MinecraftGo/cmd/internal/launcher"
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/gofrs/uuid"
)

type options struct {
	server        string
	username      string
	versionID     string
	minecraftPath string
	javaPath      string
	gameDir       string
	logPath       string
	width         int
	height        int
	maxHeap       string
	wait          time.Duration
	noQuickPlay   bool
	dryRun        bool
}

func main() {
	opts := options{}

	flag.StringVar(&opts.server, "server", "127.0.0.1:25565", "server to join on startup")
	flag.StringVar(&opts.username, "username", "Dev", "username to log in with")
	flag.StringVar(&opts.versionID, "version", constants.GameVersion, "version to launch, as the launcher names it")
	flag.StringVar(&opts.minecraftPath, "minecraft", "", "launcher directory holding versions/ and libraries/ (default: the standard .minecraft location)")
	flag.StringVar(&opts.javaPath, "java", "", "java executable to launch with (default: the launcher's runtime for this version)")
	flag.StringVar(&opts.gameDir, "game-dir", filepath.Join("run", "client"), "directory for the client's own files: its logs, options and screenshots")
	flag.StringVar(&opts.logPath, "log", "", "file to copy the client's output into (default: mcclient.log in the game directory)")
	flag.IntVar(&opts.width, "width", 854, "window width")
	flag.IntVar(&opts.height, "height", 480, "window height")
	flag.StringVar(&opts.maxHeap, "xmx", "2G", "maximum heap size, as -Xmx takes it; empty to leave it to the JVM")
	flag.DurationVar(&opts.wait, "wait", 10*time.Second, "how long to wait for the server to accept connections before launching anyway; 0 to launch straight away")
	flag.BoolVar(&opts.noQuickPlay, "no-quick-play", false, "start at the main menu instead of joining")
	flag.BoolVar(&opts.dryRun, "dry-run", false, "print the command that would run, and exit")

	flag.Parse()

	code, err := run(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mcclient:", err)
		os.Exit(1)
	}
	os.Exit(code)
}

// run launches the client and returns its exit code.
func run(opts options) (int, error) {
	minecraftDir := opts.minecraftPath
	if minecraftDir == "" {
		found, err := launcher.FindMinecraftDir()
		if err != nil {
			return 0, err
		}
		minecraftDir = found
	}

	manifest, err := launcher.ReadVersionManifest(minecraftDir, opts.versionID)
	if err != nil {
		return 0, err
	}
	if manifest.MainClass == "" || len(manifest.Arguments.Game) == 0 {
		return 0, fmt.Errorf("the manifest for %s has no main class or argument list; only 1.13 and later are supported", manifest.ID)
	}

	javaPath := opts.javaPath
	if javaPath == "" {
		javaPath, err = launcher.FindJava(minecraftDir, manifest.JavaVersion.Component)
		if err != nil {
			return 0, err
		}
	}

	fmt.Printf("Launching %s (protocol %d) as %s\n", manifest.ID, constants.ProtocolVersion, opts.username)
	fmt.Printf("Using java %s (version wants %d)\n", javaPath, manifest.JavaVersion.MajorVersion)

	command, logPath, err := buildCommand(opts, manifest, minecraftDir, javaPath)
	if err != nil {
		return 0, err
	}

	if opts.dryRun {
		fmt.Println()
		for _, argument := range command.Args {
			fmt.Println("   ", argument)
		}
		return 0, nil
	}

	if !opts.noQuickPlay && opts.wait > 0 {
		waitForServer(opts.server, opts.wait)
	}

	return launch(command, logPath)
}

// buildCommand assembles the java invocation from the manifest's own argument
// templates, so the client is started the way the launcher would start it.
func buildCommand(opts options, manifest launcher.VersionManifest, minecraftDir string, javaPath string) (*exec.Cmd, string, error) {
	gameDir, err := filepath.Abs(opts.gameDir)
	if err != nil {
		return nil, "", err
	}
	// The client writes its logs, options.txt and screenshots into the game
	// directory, so it gets one of its own: the repository root already has a
	// logs/ directory holding mcdata's generator output.
	if err := os.MkdirAll(gameDir, 0o755); err != nil {
		return nil, "", err
	}

	// Netty and LWJGL unpack the native libraries they find on the classpath
	// into here.  Keeping it inside the game directory means a stale unpack can
	// be cleared by deleting one directory.
	nativesDir := filepath.Join(gameDir, "natives")
	if err := os.MkdirAll(nativesDir, 0o755); err != nil {
		return nil, "", err
	}

	features := launcher.Features{
		"has_custom_resolution":     true,
		"is_quick_play_multiplayer": !opts.noQuickPlay,
	}

	libraries, skipped := manifest.Classpath(minecraftDir, features)
	clientJar := manifest.ClientJar(minecraftDir)
	if _, err := os.Stat(clientJar); err != nil {
		return nil, "", fmt.Errorf("no client jar at %s; run the launcher once for %s", clientJar, manifest.ID)
	}
	classpath := strings.Join(append(libraries, clientJar), string(os.PathListSeparator))
	fmt.Printf("Classpath: %d libraries from %s (%d not installed for this platform)\n", len(libraries), manifest.ID, skipped)

	values := map[string]string{
		"natives_directory": nativesDir,
		"classpath":         classpath,
		"launcher_name":     "mcclient",
		"launcher_version":  "1",
		"auth_player_name":  opts.username,
		"version_name":      manifest.ID,
		"version_type":      manifest.Type,
		"game_directory":    gameDir,
		"assets_root":       filepath.Join(minecraftDir, "assets"),
		"assets_index_name": manifest.AssetIndex.ID,
		"auth_uuid":         offlineUUID(opts.username),
		// The server only checks these when encryption is on, and it is off by
		// default, so an offline session is enough and no real account is
		// involved.  A client launched like this cannot join an online server.
		"auth_access_token": "0",
		"clientid":          "",
		"auth_xuid":         "",
		"resolution_width":  strconv.Itoa(opts.width),
		"resolution_height": strconv.Itoa(opts.height),
	}
	if !opts.noQuickPlay {
		values["quickPlayMultiplayer"] = opts.server
	}

	// Nothing here adds the manifest's logging.client.argument.  That switch
	// points log4j at a config that emits XML for the launcher to parse; left
	// off, the client's bundled config writes plain text, which is the whole
	// point of running it this way.
	jvmArguments, err := launcher.ResolveArguments(manifest.Arguments.JVM, features, values)
	if err != nil {
		return nil, "", err
	}
	gameArguments, err := launcher.ResolveArguments(manifest.Arguments.Game, features, values)
	if err != nil {
		return nil, "", err
	}

	arguments := make([]string, 0, len(jvmArguments)+len(gameArguments)+2)
	if opts.maxHeap != "" {
		arguments = append(arguments, "-Xmx"+opts.maxHeap)
	}
	arguments = append(arguments, jvmArguments...)
	arguments = append(arguments, manifest.MainClass)
	arguments = append(arguments, gameArguments...)

	logPath := opts.logPath
	if logPath == "" {
		logPath = filepath.Join(gameDir, "mcclient.log")
	}

	command := exec.Command(javaPath, arguments...)
	command.Dir = gameDir

	return command, logPath, nil
}

// offlineUUID is the id the server will give this username, from
// controller.LoginStart's offline branch.  Deriving it the same way here keeps
// the client's own idea of who it is in step with the server's.
func offlineUUID(username string) string {
	return strings.ReplaceAll(uuid.NewV3(uuid.Nil, "OfflinePlayer:"+username).String(), "-", "")
}

// waitForServer holds off launching until something is listening, so a client
// told to join on startup does not arrive before the server is up.  Failing to
// connect is not an error: the client will just show the multiplayer screen.
func waitForServer(address string, limit time.Duration) {
	deadline := time.Now().Add(limit)
	announced := false

	for {
		connection, err := net.DialTimeout("tcp", address, time.Second)
		if err == nil {
			connection.Close()
			return
		}
		if !time.Now().Before(deadline) {
			fmt.Printf("Nothing listening on %s after %s, launching anyway\n", address, limit)
			return
		}
		if !announced {
			fmt.Printf("Waiting for %s\n", address)
			announced = true
		}
		time.Sleep(250 * time.Millisecond)
	}
}

// launch runs the client, copying its output to the terminal and to a log file,
// and returns its exit code.
func launch(command *exec.Cmd, logPath string) (int, error) {
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return 0, err
	}
	logFile, err := os.Create(logPath)
	if err != nil {
		return 0, fmt.Errorf("opening %s: %w", logPath, err)
	}
	defer logFile.Close()

	fmt.Fprintf(logFile, "%s\n\n", strings.Join(command.Args, "\n"))

	// The client's log4j config writes to stdout; the JVM's own failures — a
	// missing native, a crash before logging starts — come out on stderr.  Both
	// are wanted, in the order they happened.
	output := io.MultiWriter(os.Stdout, logFile)
	command.Stdout = output
	command.Stderr = output

	if err := command.Start(); err != nil {
		return 0, fmt.Errorf("starting %s: %w", command.Path, err)
	}
	fmt.Printf("Client is pid %d, logging to %s\n", command.Process.Pid, logPath)

	// Ctrl-C should close the client rather than orphan it holding the window.
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	defer signal.Stop(interrupts)
	go func() {
		if _, ok := <-interrupts; ok {
			command.Process.Kill()
		}
	}()

	err = command.Wait()
	exit := &exec.ExitError{}
	if errors.As(err, &exit) {
		fmt.Printf("Client exited with status %d\n", exit.ExitCode())
		return exit.ExitCode(), nil
	}
	if err != nil {
		return 0, err
	}

	fmt.Println("Client exited cleanly")
	return 0, nil
}
