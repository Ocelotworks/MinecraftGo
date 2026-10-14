// Command mcverify checks the packets this server sends against the game's own
// decoders.
//
// It starts the server on a loopback port, joins it with a bare protocol client
// (helpers/joinprobe), records everything the server says through login,
// configuration and the start of play, and then decodes each packet with the
// real client's codecs.  A packet that decodes with bytes left over, or does not
// decode at all, is a packet the client would have disconnected over.
//
// That is the same check a real client makes, without a window, a launcher login
// or a round trip through someone reading a log.  cmd/mcclient is the end-to-end
// version of it; this one is quick enough to run after every protocol change, and
// it says which packet is wrong and how far the client got into it.
//
// The game jar, its libraries and the JRE come from the Minecraft launcher's
// install, as they do for cmd/mcdata and cmd/mcclient, and the decoder itself is
// DecodePackets.java, which the launcher's JRE runs as a single-file program.
//
// It reads data/ relative to the working directory, so run it from the repository
// root:
//
//	go run ./cmd/mcverify
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Ocelotworks/MinecraftGo/cmd/internal/launcher"
	"github.com/Ocelotworks/MinecraftGo/constants"
	"github.com/Ocelotworks/MinecraftGo/controller"
	"github.com/Ocelotworks/MinecraftGo/helpers/joinprobe"
)

//go:embed DecodePackets.java
var decoderSource string

type options struct {
	versionID     string
	minecraftPath string
	javaPath      string
	keepDir       string
	perID         int
	quiet         time.Duration
	serverOutput  bool
}

func main() {
	opts := options{}

	flag.StringVar(&opts.versionID, "version", constants.GameVersion, "version to check against, as the launcher names it")
	flag.StringVar(&opts.minecraftPath, "minecraft", "", "launcher directory holding versions/ and libraries/ (default: found automatically)")
	flag.StringVar(&opts.javaPath, "java", "", "java executable to run the decoder with; it must include the compiler, since the decoder is built on the fly (default: the version's own runtime)")
	flag.StringVar(&opts.keepDir, "keep", "", "directory to keep the captured packets in, for inspection")
	flag.IntVar(&opts.perID, "per-id", 4, "how many packets of each type to check in play, where the server repeats itself (login and configuration are always checked in full)")
	flag.DurationVar(&opts.quiet, "quiet", 2*time.Second, "how long the server has to stay silent in play before it counts as finished")
	flag.BoolVar(&opts.serverOutput, "server-output", false, "show the server's own console output, which includes a hex dump of every packet")
	flag.Parse()

	failed, err := run(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mcverify:", err)
		os.Exit(2)
	}
	if failed {
		os.Exit(1)
	}
}

// run reports whether any packet failed to decode.
func run(opts options) (bool, error) {
	if _, err := os.Stat(filepath.Join("data", "minecraft")); err != nil {
		return false, fmt.Errorf("no data/minecraft here: run from the repository root")
	}

	minecraftDir := opts.minecraftPath
	if minecraftDir == "" {
		found, err := launcher.FindMinecraftDir()
		if err != nil {
			return false, err
		}
		minecraftDir = found
	}

	manifest, err := launcher.ReadVersionManifest(minecraftDir, opts.versionID)
	if err != nil {
		return false, err
	}
	clientJar := manifest.ClientJar(minecraftDir)
	if _, err := os.Stat(clientJar); err != nil {
		return false, fmt.Errorf("the client jar for %s is not installed (%s): launch that version from the launcher once", opts.versionID, clientJar)
	}

	javaPath := opts.javaPath
	if javaPath == "" {
		javaPath, err = launcher.FindJava(minecraftDir, manifest.JavaVersion.Component)
		if err != nil {
			return false, err
		}
	}

	if opts.versionID != constants.GameVersion {
		fmt.Printf("Warning: the server speaks %s but is being checked against %s, so mismatches are expected\n", constants.GameVersion, opts.versionID)
	}

	workDir, err := os.MkdirTemp("", "mcverify-*")
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(workDir)

	packetDir := opts.keepDir
	if packetDir == "" {
		packetDir = filepath.Join(workDir, "packets")
	}

	packets, err := capture(opts)
	if err != nil {
		return false, err
	}
	written, err := writePackets(packetDir, packets, opts.perID)
	if err != nil {
		return false, err
	}
	fmt.Printf("Captured %d packets, checking %d of them\n", len(packets), written)

	libraries, _ := manifest.Classpath(minecraftDir, nil)
	classpath := strings.Join(append(libraries, clientJar), string(os.PathListSeparator))

	results, err := decode(javaPath, classpath, workDir, packetDir)
	if err != nil {
		return false, err
	}

	return report(results), nil
}

// capture starts the server on a loopback port and joins it.
func capture(opts options) ([]joinprobe.Packet, error) {
	if !opts.serverOutput {
		// The server prints every packet it sends, which would bury the report.
		devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
		defer devNull.Close()

		stdout := os.Stdout
		os.Stdout = devNull
		defer func() { os.Stdout = stdout }()
	}

	server := controller.CreateMinecraft()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		controller.Init(conn, nil, server)
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return joinprobe.Join(conn, "mcverify", 30*time.Second, opts.quiet)
}

// writePackets saves every login and configuration packet, and up to perID
// packets of each type in play, one file each, named <phase>_<n>_<id>.bin so the
// decoder can tell which protocol to use.
//
// Configuration packets are all kept because they differ: each registry_data packet
// is a different registry.  Play is where the server repeats itself, with hundreds of
// chunk packets that all have the same shape.
func writePackets(dir string, packets []joinprobe.Packet, perID int) (int, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return 0, err
	}

	type key struct {
		phase string
		id    int
	}
	taken := map[key]int{}
	written := 0

	for n, packet := range packets {
		k := key{packet.Phase, packet.ID}
		if packet.Phase == joinprobe.PhasePlay && taken[k] >= perID {
			continue
		}
		taken[k]++

		name := fmt.Sprintf("%s_%03d_%02x.bin", packet.Phase, n, packet.ID)
		if err := os.WriteFile(filepath.Join(dir, name), packet.Frame(), 0o644); err != nil {
			return written, err
		}
		written++
	}

	return written, nil
}

type result struct {
	status string
	file   string
	detail string
}

// decode runs DecodePackets.java over a directory of packets.
func decode(javaPath string, classpath string, workDir string, packetDir string) ([]result, error) {
	source := filepath.Join(workDir, "DecodePackets.java")
	if err := os.WriteFile(source, []byte(decoderSource), 0o644); err != nil {
		return nil, err
	}
	resultsPath := filepath.Join(workDir, "results.txt")

	// The game's bootstrap is noisy: JVM warnings on stderr and its own log lines on stdout.
	command := exec.Command(javaPath, "-cp", classpath, source, packetDir, resultsPath)
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("running the decoder with %s: %w\n%s", javaPath, err, output)
	}

	raw, err := os.ReadFile(resultsPath)
	if err != nil {
		return nil, fmt.Errorf("the decoder wrote no results: %w\n%s", err, output)
	}

	results := []result{}
	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		fields := strings.SplitN(strings.TrimRight(line, "\r"), "\t", 3)
		if len(fields) == 3 {
			results = append(results, result{fields[0], fields[1], fields[2]})
		}
	}
	return results, nil
}

// report prints the results and says whether any packet was wrong.
func report(results []result) bool {
	counts := map[string]int{}
	for _, r := range results {
		counts[r.status]++
		fmt.Printf("  %-5s %-28s %s\n", r.status, strings.TrimSuffix(r.file, ".bin"), r.detail)
	}

	fmt.Printf("\n%d ok, %d with bytes left over, %d failed, %d skipped\n", counts["OK"], counts["LEFT"], counts["FAIL"], counts["SKIP"])
	if counts["SKIP"] > 0 {
		fmt.Println("Skipped packets need registries this check does not build; they were not checked.")
	}

	return counts["LEFT"] > 0 || counts["FAIL"] > 0
}
