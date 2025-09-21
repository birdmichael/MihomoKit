// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "MihomoKit",
    platforms: [
        .iOS(.v12)
    ],
    products: [
        .library(
            name: "MihomoKit",
            targets: ["MihomoKit"]
        ),
    ],
    targets: [
        .binaryTarget(
            name: "MihomoKit",
            url: "https://github.com/birdmichael/MihomoKit/releases/download/v1.19.13/MihomoKit-v1.19.13.xcframework.zip",
            checksum: "PLACEHOLDER_CHECKSUM"
        ),
    ]
)
