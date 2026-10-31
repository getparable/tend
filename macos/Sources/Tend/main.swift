// Copyright (c) 2026 Parable, LLC
import AppKit
import ServiceManagement
import SwiftUI
import TendCore

@MainActor
final class LifecycleModel: ObservableObject {
  @Published var registration = "Disabled"
  @Published var state = LifecycleState()
  @Published var busy = false
  var connection: String { state.connection }
  var run: Run? { state.run }
  var message: String { state.message }
  let helper = SMAppService.agent(plistName: "io.parable.tend.lifecycle.agent.plist")
  let cli = Bundle.main.bundleURL.appending(path: "Contents/Helpers/tend")

  func refreshRegistration() {
    switch helper.status {
    case .enabled: registration = "Enabled"
    case .requiresApproval: registration = "Approval required or disabled in System Settings"
    case .notRegistered: registration = "Disabled (not registered)"
    case .notFound: registration = "Unavailable (helper not found)"
    @unknown default: registration = "Unavailable (unknown registration state)"
    }
  }

  func register() {
    // Never register automatically, including when macOS reports a disabled item.
    guard helper.status == .notRegistered || helper.status == .notFound else {
      state.message = "Review the existing helper in Login Items & Extensions."
      return
    }
    do { try helper.register() } catch { state.message = error.localizedDescription }
    refreshRegistration()
  }

  func unregister() async {
    do { try await helper.unregister() } catch { state.message = error.localizedDescription }
    refreshRegistration()
    await request("service.status")
  }

  func request(_ method: String, fixture: String? = nil) async {
    guard !busy else { return }
    guard method != "fixture.start" || state.canStartFixture else { return }
    busy = true
    defer { busy = false }
    refreshRegistration()
    do {
      let response = try await Bridge.call(Request(method: method, fixture: fixture), cli: cli)
      state.receive(response)
      if response.error != nil { return }
      // Opt-in installed-app verification observes the same state the view uses.
      // Normal launches never export run data. The harness supplies a private
      // scratch path and removes it after collecting its diagnostic evidence.
      if let path = ProcessInfo.processInfo.environment["TEND_LIFECYCLE_PROBE"] {
        let snapshot: [String: Any] = [
          "connection": connection, "helper_pid": response.result?.pid ?? 0,
          "run_state": run?.state ?? "none", "output": run?.output ?? "",
        ]
        if let data = try? JSONSerialization.data(withJSONObject: snapshot, options: [.sortedKeys])
        {
          try? data.write(to: URL(filePath: path), options: .atomic)
          try? FileManager.default.setAttributes([.posixPermissions: 0o600], ofItemAtPath: path)
        }
      }
    } catch {
      state.connection = (error as? BridgeError) == .incompatible ? "Incompatible" : "Unavailable"
      state.message = error.localizedDescription
    }
  }
}

struct LifecycleView: View {
  @StateObject private var model = LifecycleModel()
  var body: some View {
    VStack(alignment: .leading, spacing: 16) {
      Label("Tend lifecycle proof", systemImage: "leaf").font(.title2)
      LabeledContent("Registration", value: model.registration)
      LabeledContent("Service", value: model.connection)
      Text(model.message).foregroundStyle(.secondary).textSelection(.enabled)
      HStack {
        Button("Register Helper") { model.register() }
          .disabled(model.helper.status == .enabled || model.helper.status == .requiresApproval)
        Button("Background Settings") { SMAppService.openSystemSettingsLoginItems() }
        Button("Unregister Helper") { Task { await model.unregister() } }
      }
      Divider()
      Text("Harmless fixture").font(.headline)
      Text("Run the three-second fixture, quit Tend, then reopen to inspect its result.")
      HStack {
        Button("Run Fixture") { Task { await model.request("fixture.start", fixture: "complete") } }
          .disabled(!model.state.canStartFixture)
        Button("Run Cancellation Fixture") {
          Task { await model.request("fixture.start", fixture: "cancel") }
        }
        .disabled(!model.state.canStartFixture)
        Button("Stop") { Task { await model.request("fixture.cancel") } }
        Button("Refresh") { Task { await model.request("service.status") } }
      }.disabled(model.busy || model.connection != "Available")
      Text(model.run?.state ?? "Not run yet").font(.headline)
      ScrollView {
        Text(model.run?.output ?? "")
          .font(.system(.body, design: .monospaced))
          .textSelection(.enabled).frame(maxWidth: .infinity, alignment: .leading)
      }.background(.background.secondary)
      Text("Development proof only · No saved tasks or schedules").font(.caption).foregroundStyle(
        .secondary)
    }
    .padding(24).frame(minWidth: 680, minHeight: 400)
    .task {
      while !Task.isCancelled {
        await model.request("service.status")
        try? await Task.sleep(for: .seconds(1))
      }
    }
  }
}

struct TendApp: App {
  var body: some Scene {
    WindowGroup("Tend Development") { LifecycleView() }
  }
}

// This entry point belongs to the installed app executable, so SMAppService
// resolves the calling bundle even during repeatable attended verification.
if CommandLine.arguments.count == 3, CommandLine.arguments[1] == "--quit-owned-ui",
  let pid = Int32(CommandLine.arguments[2]),
  let app = NSRunningApplication(processIdentifier: pid),
  app.bundleIdentifier == "io.parable.tend.lifecycle"
{
  exit(app.terminate() ? 0 : 1)
} else if CommandLine.arguments.count == 3, CommandLine.arguments[1] == "--helper" {
  let helper = SMAppService.agent(plistName: "io.parable.tend.lifecycle.agent.plist")
  do {
    switch CommandLine.arguments[2] {
    case "register":
      guard helper.status == .notRegistered || helper.status == .notFound else {
        print("Refusing to re-register an existing or disabled helper; use System Settings.")
        exit(2)
      }
      try helper.register()
    case "unregister": try helper.unregister()
    case "status": break
    default: exit(2)
    }
    print(
      "SMAppService status: \(helper.status.rawValue) (0 not registered, 1 enabled, 2 approval required/disabled, 3 not found)"
    )
    exit(helper.status == .requiresApproval ? 3 : 0)
  } catch {
    print(error.localizedDescription)
    exit(1)
  }
} else {
  TendApp.main()
}
