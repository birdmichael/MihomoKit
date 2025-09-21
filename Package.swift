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
            checksum: "8f71c304a350d67dd61553cc67059795f20d6610de38687a2b8a9a0915e42e89"
        ),
    ]
)
