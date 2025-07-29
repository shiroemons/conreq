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
- リアルタイムでの進行状況表示（--streamオプション）

## インストール

### Homebrew (macOS/Linux)

```bash
# ワンライナーインストール
brew install shiroemons/tap/conreq

# または、tapを追加してからインストール
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
conreq https://httpbin.org/get

# 3つの並行GETリクエスト
conreq https://httpbin.org/get -c 3

# リアルタイムで進行状況を表示しながら実行
conreq https://httpbin.org/delay/1 -c 5 --stream

# POSTリクエストでJSONボディを送信
conreq https://httpbin.org/post -X POST -d '{"key": "value"}' -H "Content-Type: application/json"

# ファイルからボディを読み込んでPOSTリクエスト
conreq https://httpbin.org/post -X POST -d @request_body.json -H "Content-Type: application/json"

# 同一Request IDで複数リクエストを送信
conreq https://httpbin.org/anything -c 3 --same-request-id

# カスタムヘッダーで認証
conreq https://httpbin.org/bearer -H "Authorization: Bearer YOUR_TOKEN"

# PUTリクエストでデータ更新
conreq https://httpbin.org/put -X PUT -d '{"name": "updated"}' -H "Content-Type: application/json"

# DELETEリクエスト
conreq https://httpbin.org/delete -X DELETE

# ステータスコードのテスト
conreq https://httpbin.org/status/500 -c 3
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
| `--stream` | | リアルタイムで進行状況を表示 | false |
| `--output` | `-o` | 結果をファイルに出力 | 標準出力 |
| `--version` | `-v` | バージョン情報を表示 | - |
| `--help` | `-h` | ヘルプを表示 | - |

### 出力例

#### ストリーミング出力（--stream）

```
🚀 Starting 3 concurrent requests at 2025-07-29 17:38:58

+Time      Time         | Request    Status   Code   Duration  Request-ID
──────────────────────────────────────────────────────────────────────────────────────────────────────────────
[   258µs] 17:38:58.238 | Request 0  ⏳ PENDING     -          -  1baa21bf-589e-4188-a805-96213490eb14
[   278µs] 17:38:58.238 | Request 1  ⏳ PENDING     -          -  0705c6a8-e70b-4faa-b2a2-897eb2cca2c7
[   280µs] 17:38:58.238 | Request 0  🔄 RUNNING     -          -  1baa21bf-589e-4188-a805-96213490eb14
[   281µs] 17:38:58.238 | Request 2  ⏳ PENDING     -          -  f4961313-64cd-433b-9fca-2ab23bf4bb5a
[   283µs] 17:38:58.238 | Request 1  🔄 RUNNING     -          -  0705c6a8-e70b-4faa-b2a2-897eb2cca2c7
[   283µs] 17:38:58.238 | Request 2  🔄 RUNNING     -          -  f4961313-64cd-433b-9fca-2ab23bf4bb5a
[   1.77s] 17:39:00.009 | Request 2  ✅ DONE      200      1.77s  f4961313-64cd-433b-9fca-2ab23bf4bb5a
[   1.80s] 17:39:00.037 | Request 0  ✅ DONE      200      1.80s  1baa21bf-589e-4188-a805-96213490eb14
[   2.03s] 17:39:00.264 | Request 1  ✅ DONE      200      2.03s  0705c6a8-e70b-4faa-b2a2-897eb2cca2c7

🎉 All requests completed in 2.03s at 2025-07-29 17:39:00
==============================================================================================================

Final Results:

[... 続けて通常の結果出力 ...]
```

**ストリーミング出力の絵文字説明：**
- ⏳ PENDING: リクエスト待機中
- 🔄 RUNNING: リクエスト実行中
- ✅ DONE (2xx): 成功レスポンス
- ⚠️ DONE (4xx): クライアントエラー
- ❌ DONE (5xx): サーバーエラー
- ❌ FAILED: ネットワークエラーなど

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
# 5つの並行リクエストで負荷をかける（進行状況を表示）
conreq https://httpbin.org/delay/1 -c 5 --timeout 60s --stream

# 大きなレスポンスを並行取得
conreq https://httpbin.org/bytes/10240 -c 3
```

### レート制限の検証

```bash
# 100ms間隔で5つのリクエストを送信
conreq https://httpbin.org/anything -c 5 --delay 100ms

# カスタムレスポンスヘッダーを設定
conreq "https://httpbin.org/response-headers?X-RateLimit-Limit=10&X-RateLimit-Remaining=5" -c 3
```

### 冪等性の確認

```bash
# 同じX-Request-IDで複数のリクエストを送信
conreq https://httpbin.org/anything -c 3 --request-id "fixed-request-id"

# 同一Request IDで全リクエストを送信
conreq https://httpbin.org/uuid -c 3 --same-request-id
```

### 認証のテスト

```bash
# Basic認証
conreq https://httpbin.org/basic-auth/user/pass -H "Authorization: Basic dXNlcjpwYXNz"

# Bearer Token認証
conreq https://httpbin.org/bearer -H "Authorization: Bearer YOUR_TOKEN"

# APIキー認証
conreq https://httpbin.org/anything -H "X-API-Key: YOUR_API_KEY"
```

### エラーハンドリングの検証

```bash
# 特定のステータスコードをテスト
conreq https://httpbin.org/status/503 -c 5

# ランダムなステータスコード (200, 201, 400, 500のいずれかを返す)
conreq "https://httpbin.org/status/200,201,400,500" -c 5

# タイムアウトのテスト
conreq https://httpbin.org/delay/10 --timeout 3s
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
