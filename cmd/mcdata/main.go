// Command mcdata regenerates the game data this server reads from a Minecraft
// client jar.
//
// Four things come out of it:
//
//   - constants/Protocol.go, from the jar's version.json.
//   - data/minecraft/**, the vanilla datapack, extracted straight from the jar.
//     These are the registry entries and tags the server sends during
//     configuration, and the client numbers them by the order it receives them,
//     so a stale set here shifts every ID.
//   - data/blocks.json, the block state table, and constants/Packets.go, the
//     packet ids.  Neither is written down anywhere in the jar: state ids are
//     assigned at runtime by registry order, and packet ids by alphabetical
//     order of the packet name within each state.  Both come from running the
//     game's own data generator (net.minecraft.data.Main --reports).
//
// Running the generator needs a JRE new enough for the jar and the libraries
// that version depends on.  Both come from the Minecraft launcher's install,
// which mcdata finds on its own; -java and -minecraft override the search.
//
// Usage:
//
//	go run ./cmd/mcdata -jar decomp/26.1.2.jar
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type options struct {
	jarPath       string
	outputPath    string
	minecraftPath string
	javaPath      string
	versionID     string
	keepReports   string
	prune         bool
	skipProtocol  bool
	skipDatapack  bool
	skipBlocks    bool
	skipPackets   bool
}

func main() {
	opts := options{}

	flag.StringVar(&opts.jarPath, "jar", "", "path to the Minecraft client jar (required)")
	flag.StringVar(&opts.outputPath, "out", ".", "repository root to write into")
	flag.StringVar(&opts.minecraftPath, "minecraft", "", "launcher directory holding versions/ and libraries/ (default: the standard .minecraft location)")
	flag.StringVar(&opts.javaPath, "java", "", "java executable to run the data generator with (default: the launcher's runtime for this version)")
	flag.StringVar(&opts.versionID, "version", "", "version manifest to take the library list from (default: the jar's own version)")
	flag.StringVar(&opts.keepReports, "keep-reports", "", "directory to keep the raw generator output in, for inspection")
	flag.BoolVar(&opts.prune, "prune", false, "delete files under data/minecraft that the jar no longer contains")
	flag.BoolVar(&opts.skipProtocol, "skip-protocol", false, "do not rewrite constants/Protocol.go")
	flag.BoolVar(&opts.skipDatapack, "skip-datapack", false, "do not extract data/minecraft")
	flag.BoolVar(&opts.skipBlocks, "skip-blocks", false, "do not rewrite data/blocks.json")
	flag.BoolVar(&opts.skipPackets, "skip-packets", false, "do not rewrite constants/Packets.go")

	flag.Parse()

	if opts.jarPath == "" {
		fmt.Fprintln(os.Stderr, "-jar is required")
		flag.Usage()
		os.Exit(2)
	}

	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "mcdata:", err)
		os.Exit(1)
	}
}

func run(opts options) error {
	jar, err := openJar(opts.jarPath)
	if err != nil {
		return err
	}
	defer jar.Close()

	version, err := jar.readVersion()
	if err != nil {
		return err
	}
	fmt.Printf("Jar is %s, protocol %d\n", version.ID, version.ProtocolVersion)

	if !opts.skipProtocol {
		path := filepath.Join(opts.outputPath, "constants", "Protocol.go")
		if err := writeProtocolConstants(path, version); err != nil {
			return err
		}
		fmt.Printf("Wrote %s\n", path)
	}

	if !opts.skipDatapack {
		result, err := jar.extractDatapack(filepath.Join(opts.outputPath, "data"), opts.prune)
		if err != nil {
			return err
		}
		fmt.Printf("Datapack: %d files written, %d unchanged\n", result.Written, result.Unchanged)
		switch {
		case opts.prune:
			fmt.Printf("Datapack: %d stale files deleted\n", len(result.Stale))
		case len(result.Stale) > 0:
			fmt.Printf("Datapack: %d files are not in this jar — rerun with -prune to delete them\n", len(result.Stale))
			for _, path := range result.Stale[:min(5, len(result.Stale))] {
				fmt.Printf("    %s\n", path)
			}
			if len(result.Stale) > 5 {
				fmt.Printf("    ... and %d more\n", len(result.Stale)-5)
			}
		}
	}

	if opts.skipBlocks && opts.skipPackets {
		return nil
	}

	return generateFromReports(jar, version, opts)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
