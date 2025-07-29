# conreq

[![Test](https://github.com/shiroemons/conreq/actions/workflows/test.yml/badge.svg)](https://github.com/shiroemons/conreq/actions/workflows/test.yml)
[![Lint](https://github.com/shiroemons/conreq/actions/workflows/lint.yml/badge.svg)](https://github.com/shiroemons/conreq/actions/workflows/lint.yml)
[![Build](https://github.com/shiroemons/conreq/actions/workflows/build.yml/badge.svg)](https://github.com/shiroemons/conreq/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/shiroemons/conreq)](https://goreportcard.com/report/github.com/shiroemons/conreq)
[![Go Reference](https://pkg.go.dev/badge/github.com/shiroemons/conreq.svg)](https://pkg.go.dev/github.com/shiroemons/conreq)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

conreqは、同一のAPIエンドポイントに対して複数の並行HTTPリクエストを送信し、APIの挙動を検証するためのCLIツールです。

## 特徴

- 1〜5個の並行HTTPリクエストを送信
- Request IDヘッダーの自動生成またはカスタム値の設定
- 同一Request IDでの複数リクエスト送信
- Request IDヘッダー名のカスタマイズ
- リクエスト間の遅延時間設定
- タイムアウト制御
- テキストまたはJSON形式での結果出力
- ファイルからのリクエストボディ読み込み（@記法対応）
- 全HTTPメソッドのサポート（GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS）

## インストール

### Homebrew (macOS/Linux)

```bash
brew tap shiroemons/tap
brew install conreq
```

### バイナリダウンロード

[Releases](https://github.com/shiroemons/conreq/releases)ページから最新のバイナリをダウンロードしてください。

```bash
# macOS (Apple Silicon)
curl -L https://github.com/shiroemons/conreq/releases/latest/download/conreq_Darwin_arm64.tar.gz | tar xz
sudo mv conreq /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/shiroemons/conreq/releases/latest/download/conreq_Darwin_x86_64.tar.gz | tar xz
sudo mv conreq /usr/local/bin/
```

### Go install

Go 1.24以降がインストールされている場合：

```bash
go install github.com/shiroemons/conreq/cmd/conreq@latest
```

### ソースからビルド

```bash
git clone https://github.com/shiroemons/conreq.git
cd conreq
go build -o conreq cmd/conreq/main.go
```

## 使い方

### 基本的な使用方法

```bash
# 単一のGETリクエスト
conreq https://api.example.com/endpoint

# 3つの並行GETリクエスト
conreq https://api.example.com/endpoint -c 3

# POSTリクエストでJSONボディを送信
conreq https://api.example.com/endpoint -X POST -d '{"key": "value"}' -H "Content-Type: application/json"

# ファイルからボディを読み込んでPOSTリクエスト
conreq https://api.example.com/endpoint -X POST -d @request_body.json -H "Content-Type: application/json"

# 同一Request IDで複数リクエストを送信
conreq https://api.example.com/endpoint -c 3 --same-request-id
```

### コマンドラインオプション

| オプション | 短縮形 | 説明 | デフォルト |
|-----------|--------|------|------------|
| `--method` | `-X` | HTTPメソッド (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS) | GET |
| `--concurrent` | `-c` | 同時リクエスト数 (1-5) | 1 |
| `--header` | `-H` | カスタムヘッダー（複数指定可） | なし |
| `--data` | `-d` | リクエストボディ（@でファイル指定可） | なし |
| `--same-request-id` | | 全リクエストで同一のRequest IDを使用 | false |
| `--request-id` | | カスタムRequest ID値を指定 | UUID v4自動生成 |
| `--request-id-header` | | Request IDヘッダー名 | X-Request-ID |
| `--delay` | | リクエスト間の遅延時間 | 0s |
| `--timeout` | | タイムアウト時間 | 30s |
| `--no-body` | | レスポンスボディを非表示（JSON出力時は無視） | false |
| `--json` | | JSON形式で出力 | false |
| `--output` | `-o` | 結果をファイルに出力 | 標準出力 |
| `--version` | `-v` | バージョン情報を表示 | - |
| `--help` | `-h` | ヘルプを表示 | - |

### 出力例

#### テキスト形式（デフォルト）

```
=== Request Summary ===
URL: https://api.example.com/users
Method: GET
Concurrent: 3
Total Requests: 3

=== Results ===
[1] 2024-01-20 15:30:45.123456 | Status: 200 | Time: 145ms | X-Request-ID: 550e8400-e29b-41d4-a716-446655440001
{"status": "ok", "data": {...}}

[2] 2024-01-20 15:30:45.234567 | Status: 200 | Time: 132ms | X-Request-ID: 550e8400-e29b-41d4-a716-446655440002
{"status": "ok", "data": {...}}

[3] 2024-01-20 15:30:45.345678 | Status: 500 | Time: 89ms | X-Request-ID: 550e8400-e29b-41d4-a716-446655440003
{"error": "internal server error"}

=== Summary ===
Success: 2/3 (66.7%)
Failed: 1/3 (33.3%)
Average Response Time: 122ms
```

#### JSON形式

```json
{
  "metadata": {
    "url": "https://api.example.com/users",
    "method": "POST",
    "concurrent": 3,
    "total_requests": 3,
    "started_at": "2024-01-20T15:30:45.123456Z",
    "completed_at": "2024-01-20T15:30:45.456789Z",
    "total_duration_ms": 333
  },
  "results": [
    {
      "index": 1,
      "request_id": "550e8400-e29b-41d4-a716-446655440001",
      "started_at": "2024-01-20T15:30:45.123456Z",
      "completed_at": "2024-01-20T15:30:45.268456Z",
      "duration_ms": 145,
      "request": {
        "method": "POST",
        "url": "https://api.example.com/users",
        "headers": {
          "Content-Type": "application/json",
          "X-Request-ID": "550e8400-e29b-41d4-a716-446655440001"
        },
        "body": "{\"name\":\"test user\"}"
      },
      "response": {
        "status_code": 200,
        "status_text": "OK",
        "headers": {
          "Content-Type": "application/json"
        },
        "body": "{\"status\":\"ok\",\"data\":{\"id\":1,\"name\":\"test user\"}}"
      },
      "error": null
    }
  ],
  "summary": {
    "total": 3,
    "successful": 2,
    "failed": 1,
    "success_rate": 66.7,
    "average_duration_ms": 122,
    "min_duration_ms": 89,
    "max_duration_ms": 145,
    "status_codes": {
      "200": 2,
      "500": 1
    }
  }
}
```

## 使用例

### API負荷テスト

```bash
# 5つの並行リクエストで負荷をかける
conreq https://api.example.com/heavy-endpoint -c 5 --timeout 60s
```

### レート制限の検証

```bash
# 100ms間隔で5つのリクエストを送信
conreq https://api.example.com/rate-limited -c 5 --delay 100ms
```

### 冪等性の確認

```bash
# 同じX-Request-IDで複数のリクエストを送信
conreq https://api.example.com/idempotent -c 3 --request-id "fixed-request-id"
```

## 開発

### 必要なもの

- Go 1.24以降

### ビルド

```bash
go build -o conreq cmd/conreq/main.go
```

### テスト

```bash
go test ./...
```

### プロジェクト構造

```
conreq/
├── cmd/
│   └── conreq/
│       └── main.go      # CLIエントリーポイント
├── internal/
│   ├── client/          # HTTPクライアント実装
│   ├── config/          # 設定管理
│   ├── output/          # 出力フォーマッター
│   └── runner/          # 並行実行ロジック
└── pkg/
    └── requestid/       # RequestID生成
```

## ライセンス

MITライセンスに基づいて公開されています。
