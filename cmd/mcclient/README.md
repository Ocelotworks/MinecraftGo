# mcclient

Launches a real Minecraft client against this server and captures what it says.

```bash
go run . &
go run ./cmd/mcclient
```

The client joins on startup, so there is no main menu to click through, and its
output goes to the terminal and to `run/client/mcclient.log`.

## Why bother

The server only sees a connection go away. The client knows why it went away,
and says so: a disconnect reason, a registry it could not read, a stack trace
from the chunk it failed to decode. That log is the useful half of a protocol
bug, and it is easier to read when it is in a file next to the code rather than
in a window that closed.

Two details make it readable. The client is launched without the launcher's
`-Dlog4j.configurationFile` argument, so log4j uses the config bundled in the
jar and writes plain text instead of the XML the launcher parses. And it gets
its own game directory, since the client writes `logs/latest.log` where it is
started and the repository root already has a `logs/` from mcdata.

## Where the game comes from

The client jar, the libraries it needs and the JRE to run them all come from the
Minecraft launcher's install, found the same way `cmd/mcdata` finds it:

- the jar is `versions/<id>/<id>.jar`, and the library list is the manifest
  beside it;
- the JRE is the one the manifest names in `javaVersion.component`, under the
  launcher's `runtime` directory;
- the assets are the launcher's `assets/`, with the index the manifest names.

So the version has to have been launched from the launcher at least once. The
default is `constants.GameVersion`, which is the version mcdata last generated
the server's data from — the one the server can actually talk to.

Nothing is unpacked. Since 1.19 the native libraries are ordinary classpath
entries, and the manifest's `-D` arguments tell LWJGL, JNA and Netty where to
extract their own; they get `run/client/natives`.

## Logging in

The server leaves encryption off (`Minecraft.EnableEncryption`), so it never
asks Mojang whether the session is real, and an offline login is enough: the
access token is `0` and the UUID is derived the same way the server's offline
branch derives it, so both sides agree on who joined. No account is involved,
and a client started this way cannot join an online-mode server.

## Flags

| Flag | Meaning |
| --- | --- |
| `-server` | address to join (default `127.0.0.1:25565`) |
| `-username` | name to log in as (default `Dev`) |
| `-version` | version to launch (default `constants.GameVersion`) |
| `-game-dir` | directory for the client's logs, options and screenshots (default `run/client`) |
| `-log` | file to copy the output into (default `mcclient.log` in the game directory) |
| `-width`, `-height` | window size (default 854x480) |
| `-xmx` | maximum heap (default `2G`) |
| `-wait` | how long to wait for the server to start listening before launching anyway (default 10s) |
| `-no-quick-play` | start at the main menu instead of joining |
| `-dry-run` | print the command that would run, and exit |
| `-minecraft`, `-java` | override the launcher directory or the JRE |

## Limits

There is no headless Minecraft client, so a window opens. `-width`/`-height`
keep it small, but it needs a real display and a GPU: this will not run over a
plain SSH session or in CI.

Only 1.13 and later manifests are read. Older ones carry their arguments as a
single `minecraftArguments` string rather than the templated list this fills in,
and `inheritsFrom` — what modded profiles use — is not followed either.
