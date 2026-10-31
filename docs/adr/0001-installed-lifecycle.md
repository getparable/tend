# ADR 0001: one installed development service

Accepted for the lifecycle slice only, October 29, 2026.

Tend starts with macOS 26+ on Apple Silicon. The pinned development gate uses
Go 1.26.5, Apple Swift 6.4 and the macOS 27.0 SDK, compiling for macOS 26.0.
Installed lifecycle and attended GUI checks passed on macOS 27.0.1. Michael
explicitly deferred macOS 26 runtime verification for this first development
slice on October 30, 2026; the deployment target remains 26.0 and no macOS 26
runtime pass is claimed. Intel and public distribution are outside this proof.
Command Line Tools suffice; full Xcode is not needed for this build. SwiftPM's
6.4 default build system needs its bundled Testing macro loaded explicitly on
this CLT installation; `tools/swifttest` discovers it in the selected toolchain.
Every Make target exports `GOTOOLCHAIN=go1.26.5`: the `go.mod` toolchain directive
alone does not downgrade a newer host Go. This keeps worker, recursive app build
and noninteractive finalization checks on the same version. The exact Go,
Apple Swift and SDK checks remain mandatory and print their selections. Apple
tools default to `/Library/Developer/CommandLineTools` and its versioned
`SDKs/MacOSX27.0.sdk`, avoiding global Xcode selection and unversioned SDK aliases.
Explicit `DEVELOPER_DIR` and `SDKROOT` environment or Make arguments take
precedence over these defaults. Mismatches still fail with instructions to
select the supported installed tools per command; global settings are untouched.

## Installed ownership

`make app` assembles and ad-hoc signs `dist/Tend Lifecycle.app`. `make install-dev`
copies it once to **`~/Applications/Tend Lifecycle.app`** and refuses to replace
an existing bundle. This is the only supported development installation path.
Never register a bundle in a disposable build directory. Relocation, upgrade,
notarization and distribution signing belong to the packaging milestone.
On the initial host launchd rejected the ad-hoc SMAppService helper with a signing
launch-constraint violation (exit 78), despite Enabled registration. Supplying
`CODE_SIGN_IDENTITY` explicitly signs both helper and app with that identity;
the tooling never discovers or chooses a personal certificate. A trusted
development signing identity is a prerequisite for completing this host's
installed-service proof. Ad-hoc test-label coverage does not waive it. The later
MacBook Pro recovery completed both installed lifecycle harnesses with an
authorized Apple Development signature; see the verification record for artifact
identities and the later completed attended checks.

The isolated bundle ID is `io.parable.tend.lifecycle`. The SwiftUI executable
is `Contents/MacOS/Tend`; the Go CLI/helper is `Contents/Helpers/tend`. Separate
directories avoid `Tend`/`tend` collisions on case-insensitive macOS volumes.
The original MIT license is included in the bundle.

The app calls `SMAppService.agent(plistName:)` with
`io.parable.tend.lifecycle.agent.plist`, bundled in `Contents/Library/LaunchAgents`.
`BundleProgram` names `Contents/Helpers/tend`, with `serve` as an argv element.
`RunAtLoad=true`, `KeepAlive=true`, `ExitTimeOut=5` and a two-second throttle
make one user service available independently of the UI while approved and
logged in. No per-task LaunchAgents exist.

Registration is always an explicit action. Enabled, not registered, not found,
and requires approval statuses are surfaced. Apple uses requires approval for
both pending consent and revoked consent: the UI says **Approval required or
disabled in System Settings** and never re-registers that state. Protocol
incompatibility is shown separately and blocks fixture controls. Unregistering
the helper stops owned work; quitting the UI leaves it alone.

## Transport

SwiftUI invokes the bundled CLI with `Process`, `request` as argv, and a JSON
request on stdin. The CLI forwards the same envelope over a private Unix socket:
`~/Library/Application Support/TendLifecycle-dev/service.sock`. There is no
network port. The containing directory must be owned by the effective UID with
mode 0700, the socket with mode 0600. Reject symlink ancestors, relative/unclean
paths, unsafe existing files and paths longer than 103 UTF-8 bytes. Both peers
check `getpeereid` against their effective UID. This isolates users, not mutually
untrusted processes within one account. A no-follow lock plus `flock` prevents
two helpers replacing each other's socket. A new helper may replace only a
validated stale socket after acquiring the lock.

Protocol 1 uses one newline-terminated JSON frame per request/response, at most
65536 bytes including newline, with a 1–128-byte request ID. The response echoes
the ID and contains exactly one `result` or structured `error` (`code`, `message`).
Socket deadlines are five seconds with at most 16 concurrent handlers. The Swift
bridge also watches the trusted bundled CLI process. Unknown fields/methods,
missing newline, trailing JSON, oversized input, ID mismatch and incompatible
versions fail explicitly. CLI errors still emit JSON, with a nonzero exit code.

Example stdin for `Contents/Helpers/tend request`:

```json
{"version":1,"id":"example-1","method":"fixture.start","fixture":"complete"}
```

Methods are `service.status`, `fixture.start` (`complete` or `cancel`) and
`fixture.cancel`. Status returns helper PID and the latest optional run with
state, output and terminal exit code. One active fixture is admitted; duplicate
start returns `already_running`. The shared Go/Swift golden response lives in
`tests/fixtures/protocol/status.json`. There is no task CRUD or compatibility
negotiation beyond protocol 1 in this slice.

## Execution and boundaries

The Go service owns the child and one in-memory result. Embedded harmless
fixture scripts are delivered to `/bin/sh /dev/fd/3` through an inherited pipe;
there are no shell command strings, eval, credentials or temporary script files.
The child's environment is restricted and working directory is `/`. Combined
fixture output is capped at 4096 bytes and continues draining. General separate
stdout/stderr and durable execution are deferred to the execution milestone.

Cancellation targets the owned process group, waits one second after TERM,
then escalates to KILL. `waitid(WNOWAIT)` observes exit without reaping; a mutex
serializes group signals with the reaping handoff to prevent PID reuse. Remaining
attached descendants are killed before the leader is reaped. `exec.Cmd.WaitDelay`
bounds output-drain shutdown to one second. A stopped helper cancels its child.
Deliberately detached descendants are outside this fixture proof; this is not
a sandbox or general crash-recovery implementation.

KeepAlive restart clears the in-memory result. UI quit/reopen preserves it only
while the helper remains alive. SQLite, durable history, task CRUD, scheduling,
automatic retries and recovery after service crashes are separate milestones.

## Verification authority

`make check` runs formatting, vet, Go race tests, Swift golden/compatibility tests
and the development bundle build/signature check. The matching macOS CI job
runs the same gate and fails if its image lacks the pinned tools; no CI run is
claimed by a local execution. Socket integration tests explicitly skip when
the supplied scratch path exceeds Darwin's limit. Process, admission, peer and
protocol tests still run; `make verify-launchd` covers real socket IPC at the
short installed runtime path without redirecting scratch files outside TMPDIR.

`make verify-launchd` uses a unique test label and a distinct
`TendLifecycle-proof` runtime directory, refusing pre-existing state. Its plist
is scratch under TMPDIR. It never enumerates other jobs. This is automated
launchctl coverage, not SMAppService approval evidence.

`make verify-lifecycle` uses the actual installed SMAppService registration and
SwiftUI executable. Its opt-in `TEND_LIFECYCLE_PROBE` captures the model fields
displayed by that UI to a private scratch file. Those observations establish
client reconnect behavior; they are diagnostic evidence, not human design
approval or acceptance-inspection receipts. Approval/disabled-state interaction
and minimum-OS observations must be recorded honestly as attended work.

References: [Apple SMAppService](https://developer.apple.com/documentation/servicemanagement/smappservice),
the installed SDK's `SMAppService.h` and `launchd.plist(5)`, and
[Go os/exec](https://pkg.go.dev/os/exec). Installed behavior, not these references,
is the lifecycle evidence.
