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
            url: "https://github.com/birdmichael/MihomoKit/releases/download/v2.0/MihomoKit-v2.0.xcframework.zip",
            checksum: "564f892234b442b1f192dd62f50512c49e6f942788b35fdb8e7e2e94aa5d2d8c"
        ),
    ]
)
