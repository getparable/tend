# Installed lifecycle verification

Use fictional fixtures in a signed-in macOS 26+ Apple Silicon session. Do not
run an updater, inspect unrelated LaunchAgents, or copy account credentials.
This slice has no saved tasks. Public signing/notarization is not part of it.

## Latest verified acceptance, October 30–31, 2026

Native comment `cmt_4de56182ac274aa395bc8fead21e3823` and its retained
attachments supersede the historical signing, missing-artifact and attended
blockers recorded below. The report records actual GUI checks by Codex in the
signed-in MacBook Pro session on macOS 27.0.1 arm64, build 26A434. Michael
explicitly decided: "Proceed with macOS 27; defer macOS 26 verification".
This applies to this first development slice. macOS 26 remains untested; the
deployment target stays 26.0.

The tested source was base `f88b072b9a9e5cdebd7b5aa1637ea0c85f55193c` plus
staged tree `ea8ef1f801e16eff419ac96232b8f4f6b5f60925`, retained in checkpoint
`c5b251548e1644d4d6db806a0b47b42e739f1fe2`. The recovery for this handoff
matched all 40 checkpoint paths and reproduced that exact index tree before
updating these documents. Executable source is unchanged. The installed GUI
artifact and the separately rebuilt signed `dist` bundle have distinct identities.
The tested installation at `~/Applications/Tend Lifecycle.app` had SHA-256s:

- UI: `2ce5b86aaab48b72da4c8639889d7dd38a9b357aca72b75055d277fc92222b7d`.
- Helper: `b2dd8649c27f8798c7f39cfb190a09f3d91bc9d0b650d2d73c6e2556595ca053`.

The retained results are:

| Check | Actual outcome |
| --- | --- |
| Explicit registration | Enabled/Available without an initial approval prompt. |
| Disable in System Settings | Approval-required/disabled message, Unavailable, registration and run controls disabled. |
| Quit and reopen while disabled | Stayed disabled; CLI registration refused to re-register. |
| Explicit re-enable in Settings | Enabled/Available again. |
| Output selection | Actual fixture output selected; accessibility observation retained. |
| Keyboard navigation | Tab focused controls; Space activated Run Fixture. Original Keyboard Navigation setting restored. |
| Native appearance | Actual window captures in light and dark appearance; original Dark appearance restored. |
| Pinned `make check` | Formatting, vet, all Go race tests including socket `TestLifecycle`, both Swift tests and signed app build passed; no skips. |
| `make verify-launchd` | Test-label reconnect, KeepAlive restart, cancellation and teardown passed. |
| `make verify-lifecycle` | Actual SMAppService, UI quit/reconnect, result, idle restart, cancellation and unregister/child teardown passed. |
| Final cleanup | Helper status 0; scoped launchctl lookup exited 113 (job absent); no owned runtime state remained. |

The first `make check` failed before tests because the selected developer tools
resolved an old SDK path. The successful rerun selected already-installed tools
per command, without changing the global developer path:

```sh
DEVELOPER_DIR=/Library/Developer/CommandLineTools \
SDKROOT=/Library/Developer/CommandLineTools/SDKs/MacOSX27.0.sdk \
GOTOOLCHAIN=go1.26.5 CODE_SIGN_IDENTITY='<operator-supplied-identity>' make check
GOTOOLCHAIN=go1.26.5 make verify-launchd
GOTOOLCHAIN=go1.26.5 make verify-lifecycle
```

The successful build emitted nonfatal CLT linker search-path warnings. The
earlier IPC timeout did not recur; its cause remains undetermined. A disabled
state was observed, not a fabricated initial approval prompt. These retained
proofs do not claim a remote CI pass, finalized commit, published Change or merge.
Review and publication remain host-owned after this staged source handoff.

Original native evidence references are preserved here. The text attachments
were read completely during recovery and their SHA-256s checked. Screenshots
remain in the original numbered parts; this recovery does not claim fresh GUI
observations or acceptance-inspection receipts.

- [Verification report][report]: operator, commands, source and artifact identities, outcomes and deferral.
- [Evidence index][index]: original screenshot and source-patch part identities, byte ranges, lengths, hashes and reconstruction instructions.
- [Window/state index][states]: actual accessibility observations, selected output, disabled/reopened and keyboard states.
- [Pinned check log][check-log], [launchd log][launchd-log], [installed lifecycle log][lifecycle-log] and [final artifact/status log][final-log]: complete original command output.

[report]: https://cloud.detent.build/organizations/org_4c0075d611a60e617ce5/api/v2/projects/prj_05081fee186b49c98e12c28395a85b45/attachments/att_f7c8aeb7b99c9148be97279cea352af1
[index]: https://cloud.detent.build/organizations/org_4c0075d611a60e617ce5/api/v2/projects/prj_05081fee186b49c98e12c28395a85b45/attachments/att_677b9e16cb20d83783405fb3d6c8d427
[states]: https://cloud.detent.build/organizations/org_4c0075d611a60e617ce5/api/v2/projects/prj_05081fee186b49c98e12c28395a85b45/attachments/att_e9a8936318cc7349dd1c70f361120f7a
[check-log]: https://cloud.detent.build/organizations/org_4c0075d611a60e617ce5/api/v2/projects/prj_05081fee186b49c98e12c28395a85b45/attachments/att_d97bc59a884dcdddfc9b7e916225d63e
[launchd-log]: https://cloud.detent.build/organizations/org_4c0075d611a60e617ce5/api/v2/projects/prj_05081fee186b49c98e12c28395a85b45/attachments/att_9be73cc7bd71011ce29987fad48b48ae
[lifecycle-log]: https://cloud.detent.build/organizations/org_4c0075d611a60e617ce5/api/v2/projects/prj_05081fee186b49c98e12c28395a85b45/attachments/att_6f75ca81a40136ee2399184e7fb228f9
[final-log]: https://cloud.detent.build/organizations/org_4c0075d611a60e617ce5/api/v2/projects/prj_05081fee186b49c98e12c28395a85b45/attachments/att_79b70e4234b4e28a798e1065845390ba

### Recovery check, October 31, 2026

The current runner repeated `make check` with the per-command CLT/SDK selection
and `GOTOOLCHAIN=go1.26.5` shown above, using default ad-hoc signing. It passed:
formatting, vet, Go race tests, both Swift tests and the app build with strict
deep signature verification. This run's `TestLifecycle` explicitly skipped
because the supplied TMPDIR exceeds Darwin's socket-path limit. The earlier
MacBook Pro run above passed that test without skips. Scratch output stayed
under supplied TMPDIR; the Swift build directory linked there for this run.

This recovery changed only README, ADR and verification prose relative to the
tested checkpoint. Installed lifecycle and GUI checks were not repeated on
this runner; their original evidence and artifact identities remain authoritative.
No installed bundle, registration, personal task or unrelated job was changed.
Full current command output is published as native evidence with this handoff.

### UI control review repair, October 31, 2026

Independent review of head `fb25151351f8a33784d640c20736be16a0771abb`
(tree `416bad1ab3d9237f90d37c64e722cbe0a9185482`) found that duplicate
start controls stayed enabled during running/stopping and an `already_running`
response incorrectly changed the service label to Unavailable. Both start
buttons and the start-request guard now use the same active-run policy.
Application-level errors preserve the last run and an Available connection;
transport, unsafe-peer, invalid-response and incompatible-version failures
still block execution. Stop and Refresh retain their connection/busy checks.

The response/control state lives in `TendCore/LifecycleState.swift` so Swift
regression tests exercise running, stopping, rejected starts, cancellation,
terminal results and transport failures without registering a helper or running
user work. `make check` passed, including all five Swift tests and the app build
with strict signature verification. The Go socket `TestLifecycle` still skips
for the supplied long TMPDIR. Complete current output and the staged tree are
published as native evidence. Installed/attended checks were not repeated for
this repair; the earlier GUI evidence above belongs to its recorded source and
artifact identities and does not claim a fresh GUI pass for the repaired app.

## Build and automated gates

Required development tools: Go 1.26.5, Apple Swift 6.4, macOS SDK 27.0 and
`swift-format`, `codesign`, `plutil`, `launchctl` from the selected Apple tools.
The deployment target is macOS 26.0. Every Make target selects Go 1.26.5, overriding
an inherited `GOTOOLCHAIN` environment value. A first run may download it; offline
runs require the pinned toolchain already installed. `make check` prints selected
versions and refuses mismatched tools, including an unsupported command-line
Make override. Apple tool selection remains per command through `DEVELOPER_DIR`
and `SDKROOT`; no global developer settings are changed.

```sh
make check
make app
make install-dev
make verify-launchd
make verify-lifecycle
```

### Finalization gate repair, October 31, 2026

Recovery verified retained head `56a38c0c6b1e2e64267b2639e075a91bd6e43e07`
and tree `693679bb7c468a91dab49a739c8ee0f4601d7153`. Receipt
`timing_45488df73cb5147fce98d2a5b78d734f` records `make check` exiting 2 in
about 50 ms with Go 1.26.6, emitting only `make: *** [check] Error 1`.
On the unchanged source, a noninteractive `make check` passed with this host's
Go 1.26.5, Apple Swift 6.4 and SDK 27.0. Selecting Go 1.26.6 reproduced the
immediate exit 2 and identical diagnostic before any Swift or test command.
The old exact-version test rejected Go 1.26.6; `go.mod` did not force a downgrade.
This demonstrates a build-environment selection defect. It does not establish
that the historical host had no further prerequisite problems.

The Makefile now selects Go 1.26.5 for every target and reports tool versions
with actionable mismatch messages. All existing formatting, vet, race, Swift
test and bundle/signature gates are preserved. The regression check uses a real
Go 1.26.6 host executable and inherited `GOTOOLCHAIN=go1.26.6`; the repaired
`make check` selects Go 1.26.5 and passes, including its recursive app build.
An explicit unsupported Make command-line override is still rejected.
Complete reproduction and passing output, environment selections and the exact
resulting staged tree are published as native diagnostic evidence for this
handoff. No installed bundle or registration was changed. The prior installed
and GUI evidence above is retained; those checks were not repeated for this
build-only repair. The supplied long TMPDIR still causes the explicit socket
`TestLifecycle` skip; the prior no-skip and installed IPC proofs remain separate.

`make app` produces `dist/Tend Lifecycle.app`; `install-dev` refuses replacement
at `~/Applications/Tend Lifecycle.app`. Unregister and remove only this isolated
development bundle before installing a changed build. Do not replace a live
registered helper. This is not an upgrade workflow.

On the initial implementation host, ad-hoc signing passed `codesign --verify` but the
SMAppService job fails launch constraints (`OS_REASON_CODESIGNING`, exit 78).
Registration still reports Enabled: that is not proof of an available helper.
For the installed proof, an operator must provide an authorized development
signing identity, rebuild with `make app CODE_SIGN_IDENTITY='…'`, then reinstall
after unregistering/removing only the previous isolated development bundle.
The tooling never searches Keychain or chooses a personal signing identity.

### Resuming with the signed recovery bundle

The operator handoff in native comment `cmt_eaf292ce42384a61995a4000a583bcea`
identifies a development-signed copy of the existing artifact on the owning
Mac Mini at:

```text
/Users/michaelhvisser/detent-runner-bootstrap/operator-handoffs/tend-signed-20261030/Tend Lifecycle.app
```

Preserve that artifact. Before installing it, verify its strict signature with
`codesign --verify --deep --strict`, record its executable hashes and signing
identity, and correlate its provenance with checkpoint
`981fc25611c3ede7688b4d9eb4f6700f559a00f2` (tree
`7276c488e0cf5206eadbcd266b50cbca0bdefb87`). The earlier artifact manifest is
native attachment `att_fe35c610c91536987cdb0f9bb046da30`, SHA-256
`e50adf611a0fe136212c3baae0079e07d333077c5edf06302cabe55b2f45ce4f`.
Its executable hashes describe the earlier ad-hoc signatures; re-signing changes
those bytes, so the hashes alone cannot establish the signed copy's identity.

Check for an existing isolated installation and its registration before copying.
Do not replace a registered bundle. After resolving any existing installation
explicitly, copy the verified signed bundle into `dist/Tend Lifecycle.app` and
use `make install-dev`, then `make verify-lifecycle`. Run build gates before
placing the signed bundle in `dist`: both `make app` and `make check` rebuild
that path and replace its signature. A source rebuild needs a fresh authorized
development signature. If the signed handoff is absent on the current runner,
resume on its owning host or have the operator supply it; an ad-hoc replacement
does not resolve this prerequisite.

The launchd harness uses a unique `io.parable.tend.lifecycle.test.*` label,
the installed Go binary, and the private `TendLifecycle-proof` runtime directory.
It verifies real socket requests, completion with disconnected clients,
KeepAlive restart to a different PID with empty memory, cancellation, and
unregister/active-child teardown. It removes only its label and its newly
created runtime state. A pre-existing proof directory is refused.

## Actual installed app proof

`make verify-lifecycle` checks actual SMAppService state. If unregistered or not
found it explicitly registers the development helper. If macOS requires
approval or reports disabled, it stops with an error and directs the operator to
Login Items & Extensions. It never re-registers a disabled helper. Approve it
there only when intended, then rerun. A missing GUI session or missing approval
is blocked, never passed.

With the helper enabled, the harness:

1. Opens the installed SwiftUI app, waits for its model to report Available,
   and starts the harmless three-second script through the bundled CLI.
2. Observes Running in that client, asks only its owned UI process to quit, and
   checks that the helper still reports Running after the UI exits.
3. Waits for successful completion while the UI is absent, reopens the installed
   app, and records that client's Succeeded state and fixture output.
4. Signals only Tend's idle registered job and waits for KeepAlive to provide a
   different PID with an empty in-memory result.
5. Starts the TERM-resistant fixture, requests cancellation and waits for
   Cancelled. Starts one more fixture for unregister/teardown coverage.
6. Quits its UI, unregisters the helper and waits up to eight seconds for socket
   removal. Removes only the known lock and an empty development data directory.

The UI model observation uses an opt-in private scratch file. Normal app
launches do not export data. No synthetic success state is used. Inspect the
actual window as well; retain only window-scoped screenshots with fictional
output, plus a screen/state index. Do not capture other apps or the desktop.

## Attended state checks

Record who performed each check, OS, build/artifact hashes and actual outcome:

- Open the installed app and explicitly Register Helper. Record whether this
  installation becomes enabled immediately or requires approval. Do not invent
  an approval prompt if none occurs.
- If approval is required, observe the corresponding message, approve in
  System Settings, and verify Available after refresh.
- Disable only this helper in System Settings. Observe the approval/disabled
  message and unavailable execution, then quit/reopen and confirm it stays
  disabled. Re-enable explicitly in Settings if continuing the proof.
- Verify output selection, keyboard access and native light/dark appearance.
  The small proof window intentionally does not implement the full task UI.
- Run the recipe on macOS 26 before claiming minimum-OS runtime acceptance;
  this run is explicitly deferred for the first lifecycle slice as recorded above.

Unavailable attended observations remain listed as missing evidence. Automated
test-label coverage and UI model snapshots do not substitute for these steps.

## Cleanup and evidence

The automated harness unregisters after taking ownership of an idle fixture.
If registration stops at approval, registration remains available for the
operator to approve or explicitly remove. To unregister manually:

```sh
"$HOME/Applications/Tend Lifecycle.app/Contents/MacOS/Tend" --helper unregister
"$HOME/Applications/Tend Lifecycle.app/Contents/MacOS/Tend" --helper status
```

Wait for only Tend's socket to disappear. Remove only the known lifecycle lock
and empty `~/Library/Application Support/TendLifecycle-dev` directory. The
development bundle can then be removed explicitly. Never delete a directory
containing unexpected data or signal a stored PID. Do not use broad job cleanup.

Capture full command output to TMPDIR and publish UTF-8 `.log`/`.txt` attachments
in parts of at most 32768 bytes, preserving exact bytes and part identities.
Keep runtime evidence out of Git. Record the base SHA plus that the build used
uncommitted/staged source until Detent finalizes it; record executable SHA-256s
and signature verification. A local proof cannot claim a future finalized SHA.

## Historical implementation observations (superseded)

These observations preserve earlier attempts and their limitations. Use the
latest verified acceptance section above for the current evidence and deferral.

See the native issue's diagnostic attachments for exact commands, outputs,
artifact hashes and cleanup. The development host is macOS 27.0.1 arm64 with
Go 1.26.5, Swift 6.4 and SDK 27.0. Initial SMAppService status was Not Found;
explicit registration returned Enabled without an approval prompt. The source
helper subsequently failed a code-signing launch constraint; the installed
UI quit/reconnect proof therefore remains blocked on trusted signing. The source
handoff must preserve any remaining attended disabled/approval and macOS 26
checks; this document is a recipe, not a claim those checks passed.

`make check` and `make app` passed locally, including signature verification,
Go race tests and the two Swift protocol tests. The socket test explicitly
skipped because the supplied TMPDIR path exceeds Darwin's 103-byte socket
limit. `make verify-launchd` separately passed real socket IPC, idle restart,
cancellation and active-child teardown at its short isolated runtime path.
`make verify-lifecycle` failed waiting for the actual registered helper and
unregistered it during cleanup. No installed UI reconnect, disabled-state,
approval-required, macOS 26 runtime or remote CI pass is claimed.

### Recovery observations, October 30, 2026

On the subsequent runner, all 40 checkpoint file contents matched the retained
checkpoint before this documentation update. The checkpoint and base were
readable locally; no implementation was recreated. The supplied signed handoff,
`~/Applications/Tend Lifecycle.app`, and both isolated runtime directories were
absent. The scoped `launchctl print gui/501/io.parable.tend.lifecycle.agent`
returned exit 113 (service not found). With no installed app executable,
SMAppService registration status could not be queried. No registration,
approval change, installation, or job cleanup was performed.

This runner has macOS 27.0.1 arm64, Swift 6.4 and SDK 27.0. Its default Go is
1.27.1, so `make check` stopped at the pinned Go version check (exit 2).
With `GOTOOLCHAIN=go1.26.5` and the Makefile's deployment/CGO settings, focused
`gofmt -l cmd internal tools tests`, Swift strict formatting lint, `go vet ./...`
and `go test -race -v ./...` passed. The socket integration test explicitly
skipped because the supplied TMPDIR exceeds Darwin's socket path limit.
`swift test --package-path macos --scratch-path "$TMPDIR/tend-swift-build"`
with the selected toolchain's Testing macro plugin built the Swift targets and
passed both protocol tests. Full command arguments and output are retained in
the recovery diagnostic attachments. All scratch output used supplied TMPDIR.

`make verify-lifecycle` stopped before registration because the installed app
was absent (exit 2). `make app`, `make install-dev` and `make verify-launchd`
were not repeated without the signed bundle. Artifact identity/signature
verification, actual installed lifecycle, approval-required/disabled and
incompatible UI state evidence, attended appearance/keyboard/output checks,
and macOS 26 runtime evidence remain outstanding. The next prerequisite is
access to the supplied signed artifact and an attended session on its host;
the prior signing-input handoff is not a claim that installed verification ran.

### Signed installed recovery verification, October 30, 2026

This follow-up supersedes the earlier missing-artifact and missing-signature
prerequisites. The supplied bundle and authorized Apple Development identity
were present on the MacBook Pro. Native patch index `att_be7c2dc485e6c0ca87b66255e3b28240`
and all three parts were read completely through `read_attachment`. Each part's
byte range and SHA-256, the 76852-byte concatenated patch SHA-256
`d07ceb2947da1f84de5f3b371641fcfb7fd50142bf95c4a28de321282e908f72`, and the
resulting tree `7276c488e0cf5206eadbcd266b50cbca0bdefb87` were verified against
base `f88b072b9a9e5cdebd7b5aa1637ea0c85f55193c` using a scratch Git index.
All recovered source matched that tree except the previous verification notes.
All 33 issue-owned paths were inspected and preserved.

The historical installed log `att_777d8289aa8fd4ad64bc7d678a83002e` was also
read completely. It establishes a timeout, Enabled registration and an
`xpcproxy` launchd snapshot, not its cause. Its reported PID and never-exited
status do not prove readiness. The installed plist used the expected
`Contents/Helpers/tend` and `tend`, `serve` arguments. Helper and client both use
`~/Library/Application Support/TendLifecycle-dev/service.sock`; this host's
path is within the 103-byte limit. No argument or socket-path mismatch was found.

A fresh registration initially also showed `xpcproxy` and no runtime directory.
Scoped logs resolved the expected signed helper path. About nine seconds after
registration, `ps` and `sample` identified the actual Go executable and runtime
threads, and the private runtime directory existed. A subsequent unmodified
`GOTOOLCHAIN=go1.26.5 make verify-lifecycle` passed with the supplied signed
bundle, including UI quit/reopen and result reconnect. Thus the earlier timeout
was not reproduced and no underlying application startup defect was established.
Do not infer a signing failure or a permanent xpcproxy stall from the historical
snapshot. The harness timeout message now directs investigation of job state,
executable, socket and scoped diagnostics without prescribing signing as the cause.

The full `GOTOOLCHAIN=go1.26.5 make check CODE_SIGN_IDENTITY=<operator-supplied-identity>`
passed on macOS 27.0.1 arm64, Swift 6.4, SDK 27.0: formatting, vet, race tests,
both Swift protocol tests and the signed app build. `TestLifecycle` skipped
because Detent's supplied TMPDIR exceeds Darwin's socket limit; actual installed
IPC was covered by both lifecycle harnesses. The new helper's code-directory
hash matches the supplied helper's; their signed file hashes differ. Complete
executable hashes, signing details and commands are in the native diagnostic
attachments. These builds use the base plus staged/uncommitted issue source;
no future finalized commit or remote CI success is claimed.

After confirming SMAppService status 0 and no registered job, the supplied
installed bundle was preserved under TMPDIR, and `GOTOOLCHAIN=go1.26.5 make install-dev`
installed the fresh signed build. Strict deep signature verification passed.
`GOTOOLCHAIN=go1.26.5 make verify-launchd` and
`GOTOOLCHAIN=go1.26.5 make verify-lifecycle` both passed against that installation.
The actual helper restarted from PID 17256 to 17942 with empty memory. The UI
observed Running before quitting and Succeeded with fixture output on reopening.
Cancellation reached Cancelled, and unregister removed the socket and owned
state. Final registration status was 0. The isolated app remains installed and
unregistered; the original operator capture was never changed.

A separate window-evidence attempt observed the installed UI model in Available,
Running, Cancelled and Succeeded states. All three window-scoped `screencapture`
commands failed with `could not create image from window`; no screenshot was
produced or claimed. The attempt quit its owned UI, unregistered and removed
only its known runtime state. Actual screenshots, attended disable/reopen/
explicit re-enable behavior, keyboard/output selection and light/dark inspection
remain outstanding. No approval prompt occurred during the automatic enabled
registrations; an approval-required UI observation is not fabricated. macOS 26
runtime verification also remains outstanding. Passing automated lifecycle checks
does not complete these attended acceptance requirements.
