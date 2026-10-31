// swift-tools-version: 6.2
// Copyright (c) 2026 Parable, LLC
import PackageDescription

let package = Package(
  name: "Tend",
  platforms: [.macOS("26.0")],
  products: [.executable(name: "Tend", targets: ["Tend"])],
  targets: [
    .target(name: "TendCore"),
    .executableTarget(name: "Tend", dependencies: ["TendCore"]),
    .testTarget(name: "TendCoreTests", dependencies: ["TendCore"]),
  ]
)
