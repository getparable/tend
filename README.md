# Tend

Scripts and agents, on your schedule.

Tend is a planned open-source Mac app for adding, scheduling, and managing scripts and AI tasks. The intended stack is a native SwiftUI interface, a Go task runner, and macOS launchd for scheduled execution.

**Status: project bootstrap. There is no runnable app or release yet.**

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
- **launchd:** scheduled invocation of the runner, independent of the app window.

A task is a saved instruction. A run is one execution of that instruction. Scripts and AI prompts share schedules, controls, and run history.

Before implementation, settle the app-to-runner protocol and local storage format. The first version should require no hosted service or account. Tasks run with the current user's permissions. Credentials and task output stay out of source control.

## Later

- Additional scheduling options and failure notifications.
- AI runner adapters and model selection.
- Read-only discovery of existing macOS scheduled jobs.
- Import and export of task definitions.

## Inspiration

- [mac-dash](https://github.com/talhaorak/mac-dash) for visibility into macOS jobs and services.

Tend starts as an independent implementation. This repository currently contains no code copied from either project.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Early contributions should focus on the first milestone and include clear setup and verification instructions.

## License

[MIT](LICENSE).
