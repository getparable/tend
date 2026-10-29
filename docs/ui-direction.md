# UI direction

Tend uses a native Mac task interface with AIU as its family reference. The first version should make it easy to answer three questions: what will run, when will it run, and what happened last time?

This is design guidance for implementation. No Tend UI has been built yet.

## References inspected

Reviewed on October 29, 2026:

- [AIU SwiftUI source](https://github.com/getparable/aiu/blob/258a6c1f0da8fc5b6957b012b6d102ed5fc8ceb4/macos/AIUBar.swift): `Panel`, `AccountCard`, `cardSurface`, `Palette`, and `AIUBarApp`. AIU was inspected through source, not a live UI capture.

These are references for independently implemented layouts and interactions. Upstream screenshots remain linked references, not bundled Tend assets.

## Main window

Use a native split view. Start around 960 by 640 points, with a resizable sidebar around 240 points wide. These are proposed Tend dimensions to validate with the first build.

The sidebar contains search, an Add Task button, All / Scheduled / Paused filters, and a compact task list. Each row shows its name, script or AI task icon, schedule summary, and last result. Selection uses the system accent. Keep schedule state separate from execution state: a task can have its schedule paused while its last run succeeded.

The selected task opens in the detail area:

1. Name, task kind, and current execution state.
2. Run Now, Edit, and a schedule toggle. Show Stop for an active run. Put Delete in the overflow menu.
3. A short schedule summary, including the next run and local time zone.
4. Script path, interpreter, and working folder, with a readable preview where useful.
5. Recent runs showing result, start time, and duration. Selecting a run opens its output.

Use two detail columns only when space permits. At narrower widths, stack the sections. Keep logs in a selectable, monospaced output area with clear stdout and stderr labels. Run history belongs to the selected task, not a separate dashboard users must hunt through.

Avoid CPU graphs, aggregate statistics, template galleries, and provider dashboards in the first milestone. The central object is the user's task.

## Add and edit

Use a native sheet with a single primary Save action. For the first milestone, show Name, Script File, Interpreter, Arguments, Working Folder, and Schedule. Offer file and folder pickers. Schedule starts with Manual or Daily at a chosen local time.

Show validation beside the field that needs attention. Saving a task does not execute it. Make Run Now a separate, explicit action after saving.

When AI tasks arrive, add a Script / AI Prompt selector. AI Prompt replaces the script fields with Prompt, Runner, and Model while retaining the working folder and schedule. Hide model selection when a runner cannot support it; explain unavailable runners with an actionable message.

Keep advanced environment settings collapsed. Do not expose launchd property-list fields in the normal creation flow.

## Menu bar panel

Follow AIU's compact panel composition: a header, a bounded scrolling body, and a small footer. Start at AIU's 360-point panel width, 12-point outer padding, and a list height capped around 520 points. Allow text to wrap and verify sizing with larger accessibility text.

Use a Tend heading with Add Task and Open Tend actions. Show active runs first, then upcoming tasks, with recent failures easy to find. Each task gets a name, a short state or next-run label, and a Run Now or Stop action. Open the main window for editing and full logs.

Keep header controls visible while the task list scrolls. Closing this panel or the main window does not pause schedules. A menu bar icon should describe active runs and failures accessibly, without constant animation.

A searchable quick launcher is a later enhancement. The first version needs one reliable menu bar entry point.

## Appearance

Follow AIU's system-driven palette rather than inventing fixed brand colors:

- Use the macOS system font, with semibold titles, regular body text, and secondary captions. Reserve monospaced text for commands, paths, and output.
- Use native window and sidebar backgrounds, primary and secondary text styles, and the user's system accent color.
- Use semantic success, warning, and failure colors with an icon and a text label. Color alone never carries status.
- Use material surfaces for grouped content. AIU's `cardSurface` uses thick material, a subtle border, and 16-point corners; its compact header uses 18-point corners. Use those as starting values, not a reason to wrap every row in a card.
- Use glass effects sparingly for supported native controls and window chrome. AIU deliberately avoids stacked glass on scrolling cards because it reduces clarity and can overlap pinned headers.
- Start with a 4, 8, 12, 16, 24-point spacing scale. Align labels and content to consistent leading edges.
- Support light and dark appearance, Increase Contrast, Reduce Transparency, and Reduce Motion. Keep log backgrounds stable enough for long reading sessions.

Choose the minimum macOS version separately from the visual reference. If Tend supports versions without the chosen glass APIs, use native material or solid-background fallbacks.

## States and behavior

- Empty list: “Add your first task” with an Add Task button.
- Never run: “Not run yet.” Do not imply success before there is a result.
- Running: show elapsed time and Stop; disable duplicate manual starts for the same task in the first version.
- Success: show completion time and offer View Output.
- Failure: show a short cause, such as “Script exited with code 1,” and offer View Output and Run Again.
- Paused schedule: show “Schedule paused.” Manual execution remains available.
- Stopping: show “Stopping…” until the runner confirms the process has ended.
- Unavailable runner: explain what is missing and keep the task editable.

Use standard keyboard navigation, visible focus, accessible labels on icon buttons, and selectable output. Add Task uses Command-N; search uses Command-F. Escape dismisses a sheet without silently discarding edited fields.

## First visual acceptance pass

Build and inspect the updater-script workflow in both appearances: create, run, inspect output, schedule daily, and pause. Exercise empty, running, failed, and never-run states with fictional task data. Confirm long names and paths remain readable, the menu bar header stays visible while scrolling, and all actions work from the keyboard.

The result should feel related to AIU and keep task organization clear. Validate the actual SwiftUI implementation with screenshots before treating these dimensions and material choices as final.
