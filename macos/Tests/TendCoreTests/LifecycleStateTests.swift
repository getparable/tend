// Copyright (c) 2026 Parable, LLC
import Foundation
import Testing

@testable import TendCore

private func status(_ run: String?) throws -> Response {
  let payload = run.map { #", "run":{"state":"\#($0)","output":"fixture output"}"# } ?? ""
  return try Response.decode(
    Data(#"{"version":1,"id":"test","result":{"pid":123\#(payload)}}"#.utf8))
}

private func failure(_ code: String) throws -> Response {
  try Response.decode(
    Data(#"{"version":1,"id":"test","error":{"code":"\#(code)","message":"Rejected"}}"#.utf8))
}

@Test func startControlsFollowFixtureLifecycle() throws {
  var state = LifecycleState()
  #expect(!state.canStartFixture)
  state.receive(try status(nil))
  #expect(state.canStartFixture)
  for active in ["running", "stopping"] {
    state.receive(try status(active))
    #expect(!state.canStartFixture)
    #expect(state.connection == "Available")
  }
  for terminal in ["succeeded", "failed", "cancelled"] {
    state.receive(try status(terminal))
    #expect(state.canStartFixture)
  }
}

@Test func rejectedStartKeepsActiveRunAndConnection() throws {
  var state = LifecycleState()
  for active in ["running", "stopping"] {
    state.receive(try status(active))
    state.receive(try failure("already_running"))
    #expect(state.connection == "Available")
    #expect(state.run?.state == active)
    #expect(state.run?.output == "fixture output")
    #expect(state.message == "Rejected")
    #expect(!state.canStartFixture)
  }
  // A later cancellation result releases the start controls again.
  state.receive(try status("cancelled"))
  #expect(state.connection == "Available")
  #expect(state.canStartFixture)
}

@Test func operationErrorsDifferFromConnectionFailures() throws {
  var state = LifecycleState()
  for code in ["spawn_failed", "already_running", "invalid_request"] {
    state.receive(try failure(code))
    #expect(state.connection == "Available")
  }
  for code in ["service_unavailable", "unsafe_peer", "invalid_response", "incompatible_version"] {
    state.receive(try failure(code))
    #expect(state.connection == (code == "incompatible_version" ? "Incompatible" : "Unavailable"))
    #expect(!state.canStartFixture)
  }
}
