// Copyright (c) 2026 Parable, LLC

// Shared by the native model and its response/control regression tests.
public struct LifecycleState: Sendable {
  public var connection = "Unavailable"
  public var run: Run?
  public var message = "Register the development helper to begin."

  public init() {}

  public var canStartFixture: Bool {
    connection == "Available" && run?.state != "running" && run?.state != "stopping"
  }

  public mutating func receive(_ response: Response) {
    if let error = response.error {
      switch error.code {
      case "incompatible_version": connection = "Incompatible"
      case "service_unavailable", "unsafe_peer", "invalid_response": connection = "Unavailable"
      default:
        // A helper rejecting an operation is still reachable. Preserve the
        // latest run so an already_running response cannot hide Stop.
        connection = "Available"
      }
      message = error.message
      return
    }
    connection = "Available"
    run = response.result?.run
    message = "Connected to helper PID \(response.result?.pid ?? 0). Results live in helper memory."
  }
}
