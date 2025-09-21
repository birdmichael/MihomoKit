<p align="center">
  <img src="https://github.com/user-attachments/assets/8b2979ca-0dc9-4419-9dfd-809714542ca1" height="150">
  <h1 align="center">MihomoKit</h1>
</p>

<pre align="center">
🧪 Working in Progress
</pre>

> I want to continue developing it, but the fact is that according to Apple's policy, Network Extension and Personal VPN permissions cannot be enabled through a **free developer account**. Maybe I won't update it until I get a paid developer account.

## Installation

### Swift Package Manager

在 Xcode 中：
1. File → Add Package Dependencies
2. 输入仓库 URL：`https://github.com/birdmichael/MihomoKit.git`
3. 选择版本并添加到项目

或在 Package.swift 中：
```swift
dependencies: [
    .package(url: "https://github.com/birdmichael/MihomoKit.git", from: "1.19.13")
]
```

## Usage

```swift
import MihomoKit

// 设置日志记录器
MihomoSetupLogger(YourLoggerImplementation())

// 初始化 Mihomo
MihomoSetup("home_directory_path", "config_content")
```

## Author

MihomoKit © MihomoX, Released under MIT. Created on Jul 17, 2024

> [Personal Website](http://wibus.ren/) · [Blog](https://blog.wibus.ren/) · GitHub [@wibus-wee](https://github.com/wibus-wee/) · Telegram [@wibus✪](https://t.me/wibus_wee)

