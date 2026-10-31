// Copyright (c) 2026 Parable, LLC
import Foundation
import Testing

@testable import TendCore

@Test func sharedGolden() throws {
  let root = URL(filePath: #filePath).deletingLastPathComponent()
    .deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent()
  let data = try Data(contentsOf: root.appending(path: "tests/fixtures/protocol/status.json"))
  let response = try Response.decode(data, id: "fixture-1")
  #expect(response.result?.run?.state == "succeeded")
  #expect(response.result?.run?.exitCode == 0)
}

@Test func rejectsMismatch() throws {
  for json in [
    #"{"version":2,"id":"a","result":{"pid":1}}"#,
    #"{"version":1,"id":"b","result":{"pid":1}}"#,
    #"{"version":1,"id":"a"}"#,
    #"{"version":1,"id":"a","result":{"pid":1},"error":{"code":"x","message":"x"}}"#,
  ] {
    #expect(throws: (any Error).self) { try Response.decode(Data(json.utf8), id: "a") }
  }
}
