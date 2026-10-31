// Copyright (c) 2026 Parable, LLC
import Foundation

public struct Request: Encodable, Sendable {
  public let version = 1
  public let id: String
  public let method: String
  public let fixture: String?

  public init(method: String, fixture: String? = nil) {
    id = UUID().uuidString
    self.method = method
    self.fixture = fixture
  }
}

public struct ServiceError: Decodable, Sendable {
  public let code: String
  public let message: String
}

public struct Run: Decodable, Sendable {
  public let state: String
  public let output: String
  public let exitCode: Int?
  enum CodingKeys: String, CodingKey {
    case state, output
    case exitCode = "exit_code"
  }
}

public struct Result: Decodable, Sendable {
  public let pid: Int
  public let run: Run?
}

public struct Response: Decodable, Sendable {
  public let version: Int
  public let id: String
  public let result: Result?
  public let error: ServiceError?

  public static func decode(_ data: Data, id: String? = nil) throws -> Response {
    guard data.count <= 65536 else { throw BridgeError.invalidResponse }
    let response = try JSONDecoder().decode(Response.self, from: data)
    guard response.version == 1 else { throw BridgeError.incompatible }
    guard id == nil || response.id == id,
      (response.result == nil) != (response.error == nil)
    else { throw BridgeError.invalidResponse }
    return response
  }
}

public enum BridgeError: Error, LocalizedError, Equatable {
  case unavailable, incompatible, invalidResponse, timeout
  public var errorDescription: String? {
    switch self {
    case .unavailable: "The bundled helper client is unavailable."
    case .incompatible:
      "Incompatible helper protocol. Unregister the development helper before replacing the app."
    case .invalidResponse: "The helper returned an invalid response."
    case .timeout: "The helper client timed out."
    }
  }
}

public enum Bridge {
  // Work runs off the main actor. One request per process, no shell strings.
  // The trusted bundled CLI enforces socket framing and a five-second deadline.
  public static func call(_ request: Request, cli: URL) async throws -> Response {
    try await Task.detached {
      let process = Process()
      process.executableURL = cli
      process.arguments = ["request"]
      let input = Pipe()
      let output = Pipe()
      process.standardInput = input
      process.standardOutput = output
      process.standardError = FileHandle.nullDevice
      do { try process.run() } catch { throw BridgeError.unavailable }
      let watchdog = DispatchWorkItem { if process.isRunning { process.terminate() } }
      DispatchQueue.global().asyncAfter(deadline: .now() + 7, execute: watchdog)
      defer { watchdog.cancel() }
      var data = try JSONEncoder().encode(request)
      data.append(10)
      try input.fileHandleForWriting.write(contentsOf: data)
      try input.fileHandleForWriting.close()
      let response = output.fileHandleForReading.readDataToEndOfFile()
      process.waitUntilExit()
      return try Response.decode(response, id: request.id)
    }.value
  }
}
