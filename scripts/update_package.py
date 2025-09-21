#!/usr/bin/env python3
import hashlib
import os
import re
import sys
import urllib.request
from pathlib import Path

def calculate_sha256(file_path_or_url):
    """计算文件或URL的SHA256"""
    if file_path_or_url.startswith('http'):
        # 下载文件并计算校验和
        with urllib.request.urlopen(file_path_or_url) as response:
            data = response.read()
            return hashlib.sha256(data).hexdigest()
    else:
        # 本地文件
        with open(file_path_or_url, 'rb') as f:
            return hashlib.sha256(f.read()).hexdigest()

def update_package_swift(version, url, checksum):
    """更新 Package.swift"""
    package_file = Path(__file__).parent.parent / "Package.swift"
    
    with open(package_file, 'r') as f:
        content = f.read()
    
    # 更新 URL
    content = re.sub(
        r'url: "https://github\.com/birdmichael/MihomoKit/releases/download/[^"]+/[^"]+\.xcframework\.zip"',
        f'url: "{url}"',
        content
    )
    
    # 更新 checksum
    content = re.sub(
        r'checksum: "[^"]+"',
        f'checksum: "{checksum}"',
        content
    )
    
    with open(package_file, 'w') as f:
        f.write(content)
    
    print(f"Updated Package.swift with version {version}, checksum: {checksum[:8]}...")

def main():
    if len(sys.argv) != 2:
        print("Usage: python3 update_package.py <version>")
        print("Example: python3 update_package.py v1.19.13")
        sys.exit(1)
    
    version = sys.argv[1]
    
    # 构建下载 URL
    url = f"https://github.com/birdmichael/MihomoKit/releases/download/{version}/MihomoKit-{version}.xcframework.zip"
    
    print(f"Calculating checksum for {url}...")
    
    try:
        checksum = calculate_sha256(url)
        update_package_swift(version, url, checksum)
        print(f"✅ Package.swift updated successfully!")
        print(f"   Version: {version}")
        print(f"   URL: {url}")
        print(f"   Checksum: {checksum}")
    except Exception as e:
        print(f"❌ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
