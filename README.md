# Tend

Scripts and agents, on your schedule.

Tend is an open-source Mac app in development for adding, scheduling, and managing scripts and AI tasks. The stack is a native SwiftUI interface, a Go task service, and macOS launchd supervising that service.

**Status: installed lifecycle proof.** A development app and Go helper can run a
harmless in-memory fixture. Saved tasks, schedules and releases are not yet
implemented.

## Development lifecycle proof

Initial scope is macOS 26+ on Apple Silicon. The build gate pins Go 1.26.5,
Apple Swift 6.4 and macOS SDK 27.0, targeting macOS 26.0.
Installed and GUI verification passed on macOS 27.0.1. macOS 26 runtime
verification is explicitly deferred for this first development slice.

```sh
make check
make app
make install-dev
make verify-launchd
make verify-lifecycle
```

The development bundle installs only at `~/Applications/Tend Lifecycle.app`,
with a separate development bundle ID and runtime directory. Open it to inspect
service state, explicitly register/unregister its helper, and run/stop harmless
fixtures. Quitting the UI leaves the service running. Results disappear when
the helper restarts. A disabled helper is never silently re-enabled.
The default build uses ad-hoc signing. For installed SMAppService verification,
build with an explicitly supplied authorized development identity:
`make app CODE_SIGN_IDENTITY='<operator-supplied-identity>'`. See the verification
recipe before installing or replacing a development bundle.

See the [architecture decision](docs/adr/0001-installed-lifecycle.md) and
[verification recipe](docs/verification/lifecycle.md) for protocol, approval,
setup, cleanup and the distinction between automated and attended checks.

## First milestone

Complete one useful workflow: add a local script, run it, inspect its output, schedule it daily, and pause its schedule.

- Add a named task with a script path, interpreter, arguments, and working directory.
- List tasks with their last result and next scheduled run.
- Run now, stop an active run, edit a task, or delete it.
- Choose a daily schedule and enable or pause it.
- Keep run history with timestamps, exit status, stdout, and stderr.
- Let macOS trigger scheduled execution even when the app is closed.

The initial app manages tasks created in Tend. Inspecting existing LaunchAgents and cron entries comes later, starting with a read-only view.

## AI tasks

The next milestone adds a prompt task with a working directory, runner, and optional model selection. Installed AI command-line tools and custom commands are the first integration path. Each supported runner needs an adapter for noninteractive execution, arguments, cancellation, and output.

The goal is to support different providers and local models. No model integrations are implemented yet.

## Proposed architecture

- **SwiftUI app:** task list, editor, run details, and menu bar access.
- **Go runner:** process execution, output capture, cancellation, and durable run records.
- **launchd:** one user helper supervised independently of the app window; the Go service will own scheduling.

A task is a saved instruction. A run is one execution of that instruction. Scripts and AI prompts share schedules, controls, and run history.

The lifecycle proof establishes a versioned JSON CLI/socket protocol. Durable local storage is the next slice. The first version requires no hosted service or account. Tasks run with the current user's permissions. Credentials and task output stay out of source control.

## Later

- Additional scheduling options and failure notifications.
- AI runner adapters and model selection.
- Read-only discovery of existing macOS scheduled jobs.
- Import and export of task definitions.

## Inspiration

- [AIU](https://github.com/getparable/aiu) for native Mac styling, compact menu bar presentation, system typography, and restrained material surfaces.
- [mac-dash](https://github.com/talhaorak/mac-dash) for future visibility into existing macOS jobs and services.

The [original Tend mock-up](docs/design/tend-main-window.svg) and [UI direction](docs/ui-direction.md) define the first screens and interaction rules.

Tend is an independent implementation. This repository contains no code copied from these projects.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Early contributions should focus on the first milestone and include clear setup and verification instructions.

## License

[MIT](LICENSE).
