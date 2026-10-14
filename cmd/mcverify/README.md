# mcverify

Checks the packets this server sends against the game's own decoders.

```bash
go run ./cmd/mcverify
```

Run it from the repository root, since the server reads `data/` relative to where
it is started. It exits 0 when every packet decodes cleanly, 1 when one does not,
and 2 when it could not run at all.

```
Captured 333 packets, checking 47 of them
  OK    configuration_001_0d         ClientboundUpdateEnabledFeaturesPacket
  OK    configuration_002_0f         ClientboundSelectKnownPacks
  ...
  FAIL  play_044_2e                  at byte 2267 of 55579: ... ByteArray with size 2097151 is bigger than allowed 2048

43 ok, 0 with bytes left over, 4 failed, 1 skipped
```

## What it does

1. Starts the server on a loopback port, in this process.
2. Joins it with a bare protocol client (`helpers/joinprobe`): handshake, login,
   configuration, then it listens until the server has been quiet for two seconds.
3. Saves what the server sent, one file per packet.
4. Runs `DecodePackets.java` over those files with the real client's codecs, from
   the client jar and libraries in the Minecraft launcher's install.

A packet that decodes but leaves bytes unread (`LEFT`) is one the server filled
with more than the client expects. One that throws (`FAIL`) is one the client
would have disconnected over, and the message says how far into it the client got.
That byte offset is usually the quickest way to find the field that is wrong:
count along the packet to it.

## Why bother

The client's disconnect report names the packet, and sometimes the field, but
getting it means launching the game, logging in, and reading a log. Every
protocol mismatch while moving to 26.3 was found that way, one packet per round
trip. This checks every packet the server sends on joining in a few seconds, and
reads the game's own code for the answer instead of the protocol documentation,
which lags behind it. The 26.3 change to how a `BitSet` is encoded is the example:
the wire format changed in the codec, and nothing written down said so.

`cmd/mcclient` is the end-to-end version, with a real client and a real window.
Use this after a change and that one before a release.

## What it covers

- **Login and configuration**: every packet, since each `registry_data` packet is
  a different registry.
- **Play**: the first few packets of each type (`-per-id`, default 4) that the
  server sends as the player joins. It sends a few hundred chunk packets that all
  have the same shape, so more than a few add time and nothing else.
- **Not** serverbound packets, or anything the server sends later than the join.
  A packet the server only sends in response to something, such as a chat message,
  is not in the capture. `-quiet` makes it wait longer for late packets.

Some packets cannot be checked here and are reported as `SKIP`. The login packet
in play carries the dimension as a registry id, which the client resolves against
the registries the server sent in configuration. The check does not build those,
so it cannot decode that one. It still lays out as the client's codec says, and a
real client is the way to confirm it.

## Where the game comes from

The same place `cmd/mcdata` and `cmd/mcclient` get it:

- the client jar is `versions/<id>/<id>.jar` and the libraries are the ones its
  manifest lists, all from the launcher's install;
- the JRE is the one the manifest names in `javaVersion.component`. The decoder is
  a single Java source file, compiled as it is run, and the launcher's JRE includes
  the compiler for that.

The version has to have been launched from the launcher at least once. The default
is `constants.GameVersion`, the version the server's data was last generated from.

## Flags

| Flag | Meaning |
| --- | --- |
| `-version` | version to check against, as the launcher names it |
| `-minecraft` | launcher directory holding `versions/` and `libraries/` |
| `-java` | java executable to run the decoder with; it must include the compiler |
| `-keep` | directory to keep the captured packets in, for inspection |
| `-per-id` | how many packets of each type to check in play |
| `-quiet` | how long the server has to stay silent in play before it is finished |
| `-server-output` | show the server's own console output, including a hex dump of every packet |

`-keep` leaves one file per packet, named `<phase>_<n>_<id>.bin`, holding what the
client's codec reads: the packet id as a VarInt, then the body.

## Checking it still works

The check is only worth having if it fails when the server is wrong. To see it do
that, change `dataTypes.WriteBitSet` back to a VarInt count of big-endian longs,
the encoding before 26.3, and run it: every chunk packet fails with the same
`ByteArray with size 2097151 is bigger than allowed 2048` the real client gave.
