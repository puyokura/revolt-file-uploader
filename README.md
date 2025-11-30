# revolt-file-uploader

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Revolt のファイルサイズ制限（20MB）を超える大きなファイルを、分割して送信・復元するための CLI ツールです。

## 特徴

- **自動分割アップロード**: 20MB を超えるファイルを自動的に 15MB 単位で分割し、Revolt にアップロードします。
- **メタデータ生成**: 分割されたファイルの情報をまとめた JSON ファイルを生成し、復元時に使用します。
- **簡単復元**: 生成された JSON ファイルを指定するだけで、分割されたファイルを自動的にダウンロード・結合し、元のファイルを復元します。

## インストール

### 最新ビルド（nightly.link）

mainブランチの最新ビルドを以下からダウンロードできます:

- **Linux (AMD64)**: [rev-up-linux-amd64.zip](https://nightly.link/puyokura/revolt-file-uploader/workflows/ci_release/main/rev-up-linux-amd64.zip)
- **Linux (ARM64)**: [rev-up-linux-arm64.zip](https://nightly.link/puyokura/revolt-file-uploader/workflows/ci_release/main/rev-up-linux-arm64.zip)
- **Windows (AMD64)**: [rev-up-windows-amd64.zip](https://nightly.link/puyokura/revolt-file-uploader/workflows/ci_release/main/rev-up-windows-amd64.zip)
- **macOS (Intel)**: [rev-up-darwin-amd64.zip](https://nightly.link/puyokura/revolt-file-uploader/workflows/ci_release/main/rev-up-darwin-amd64.zip)
- **macOS (Apple Silicon)**: [rev-up-darwin-arm64.zip](https://nightly.link/puyokura/revolt-file-uploader/workflows/ci_release/main/rev-up-darwin-arm64.zip)

### Go Install

```bash
go install github.com/puyokura/revolt-file-uploader@latest
```

または、[Releases](https://github.com/puyokura/revolt-file-uploader/releases) からバイナリをダウンロードしてください。

## 使い方

### 1. ファイルの送信 (Send)

指定したファイルを Revolt のサーバー/チャンネルにアップロードします。

```bash
rev-up send <ファイルパス> [フラグ]
```

**フラグ:**

- `-s, --server <ServerID/ServerName>`: 送信先のサーバーを指定します。
- `-c, --channel <ChannelID/ChannelName>`: 送信先のチャンネルを指定します。

**例:**

```bash
rev-up send ./large-video.mp4 -s "My Server" -c "general"
```

20MB 以上のファイルの場合、自動的に分割され、各パーツと復元用 JSON ファイルがアップロードされます。

### 2. ファイルの復元 (Repair)

分割アップロードされたファイルを復元します。Revolt 上にある復元用 JSON ファイルをダウンロードし、そのパスを指定してください。

```bash
rev-up repair <JSONファイルパス>
```

**例:**

```bash
rev-up repair ./large-video.mp4.json
```

このコマンドを実行すると、JSON に記載された情報を元に分割ファイルをダウンロードし、元のファイルをカレントディレクトリに復元します。

## 仕様

- **分割サイズ**: 15MB (Revolt の制限 20MB に対して安全マージンを確保)
- **メタデータ**: 分割ファイルの ID や順序を記録した JSON 形式

## ライセンス

[MIT License](LICENSE)
