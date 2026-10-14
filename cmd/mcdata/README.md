# mcdata

Regenerates the game data this server reads from a Minecraft client jar.

```bash
go run ./cmd/mcdata -jar decomp/26.1.2.jar
```

## What it produces

| Output | Source | Why it matters |
| --- | --- | --- |
| `constants/Protocol.go` | the jar's `version.json` | the protocol number sent in the handshake |
| `data/minecraft/**` | the jar's `data/minecraft` tree | the registry entries and tags sent during configuration; the client numbers them by the order it receives them, so a stale set shifts every id |
| `data/blocks.json` | the game's own data generator | block state ids, which decide what the client actually renders for each block in a chunk |
| `constants/Packets.go` | the game's own data generator | packet ids, which decide whether the client can decode anything the server sends at all |

## Why the data generator

Neither block state ids nor packet ids are written down anywhere in the jar.

Packet ids are handed out per state and flow in alphabetical order of the
packet's registry name, so one packet added upstream renumbers every packet
after it. 26.1.2 added `game_rule_values` and `low_disk_space_warning` to
play/clientbound and `attack`, `set_game_rule` and `spectate_entity` to
play/serverbound, which is enough to move most of the ids either side of them.
Worse, the two directions are alphabetical by *different* name sets — the
protocol wiki calls `spectate_entity` "Teleport To Entity" — so a
hand-transcribed table drifts in ways a simple offset does not explain.

Block state ids are not written down anywhere in the jar either. They are handed out at
runtime: `Blocks`'s static initialiser walks the block registry in
registration order and adds every state of every block to
`Block.BLOCK_STATE_REGISTRY` in turn. The only way to read them out is to let
the game do it, so `mcdata` runs `net.minecraft.data.Main --reports` and
converts `reports/blocks.json` and `reports/registries.json`.

That needs a JRE new enough for the jar plus the libraries that version depends
on. Both come from the Minecraft launcher's install:

- the library list comes from `versions/<id>/<id>.json`, and the jars from
  `libraries/`;
- the JRE is the one the manifest names in `javaVersion.component`, found under
  the launcher's `runtime` directory.

Libraries the launcher has not downloaded are skipped — a manifest lists the
natives for every platform and only the current one's are on disc.

If the launcher is somewhere unusual, or you would rather use your own JDK,
pass `-minecraft` and `-java`.

## Flags

| Flag | Meaning |
| --- | --- |
| `-jar` | path to the Minecraft client jar (required) |
| `-out` | repository root to write into (default `.`) |
| `-minecraft` | launcher directory holding `versions/` and `libraries/` |
| `-java` | java executable to run the data generator with |
| `-version` | version manifest to take the library list from, if it does not match the jar |
| `-keep-reports` | keep the raw generator output in a directory for inspection |
| `-prune` | delete files under `data/minecraft` that the jar no longer contains |
| `-skip-protocol`, `-skip-datapack`, `-skip-blocks`, `-skip-packets` | run only part of it |

## Limits

Some fields in `data/blocks.json` — hardness, resistance, harvest tools, light
levels, drops, bounding box — are gameplay values the reports do not describe.
They are carried over from the existing file for blocks that still exist, and
left at defaults for blocks that are new. Drops in particular hold item ids,
which can move between versions, so treat them as approximate until something
actually depends on them.

Stale files under `data/minecraft` are reported but not deleted unless you pass
`-prune`, since the directory is checked in and deleting a few hundred files is
worth doing deliberately.

Constant names in `constants/Packets.go` come from the registry names, so a
packet renamed upstream renames its constant and the call sites stop compiling.
That is deliberate — it is the only signal that the packet you were sending is
not the packet you are sending now.
