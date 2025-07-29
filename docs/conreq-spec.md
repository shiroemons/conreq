# conreq - HTTP並行リクエストテストツール仕様書

## 概要

conreqは、API動作確認を目的とした同時HTTPリクエストを送信するGo製CLIツールです。同一エンドポイントに対して複数のリクエストを同時または遅延付きで送信し、その結果を詳細に記録・表示します。

## 基本要件

### 目的
- API動作確認における同時リクエストのテスト
- 同時実行時のAPI挙動の検証
- RequestIDベースのトラッキング

### 制約
- 同時リクエスト数: 1〜5
- 負荷テスト用途ではない
- Go言語で実装

## 機能仕様

### コアフィーチャー

#### 1. HTTPリクエスト送信
- 対応メソッド: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- カスタムヘッダー対応
- リクエストボディ対応（文字列またはファイル）
- タイムアウト設定（デフォルト: 30秒）

#### 2. 同時実行制御
- 並行数: 1〜5（デフォルト: 1）
- リクエスト間遅延: オプション（デフォルト: 0）
- 最初のリクエストは即座に実行、2番目以降に遅延適用

#### 3. RequestID管理
- デフォルトヘッダー名: `X-Request-ID`
- デフォルト値: UUID v4（リクエストごとに異なる）
- オプション:
  - 同一ID使用モード
  - カスタムID指定
  - ヘッダー名変更

#### 4. 結果出力
- 通常出力: 人間が読みやすい形式
- JSON出力: 機械処理可能な形式
- レスポンスボディ表示制御（通常出力時のみ有効）
- ファイル出力対応

## コマンドライン仕様

### 基本構文
```bash
conreq [OPTIONS] URL
```

### オプション

#### 基本オプション
| オプション | 短縮形 | 説明 | デフォルト |
|-----------|--------|------|-----------|
| --method | -X | HTTPメソッド | GET |
| --header | -H | リクエストヘッダー（複数指定可） | なし |
| --data | -d | リクエストボディ（@でファイル指定可） | なし |
| --concurrent | -c | 同時リクエスト数（1-5） | 1 |

#### RequestID関連
| オプション | 説明 | デフォルト |
|-----------|------|-----------|
| --same-request-id | 全リクエストで同一のRequest IDを使用 | false |
| --request-id | カスタムRequest ID値を指定 | UUID v4自動生成 |
| --request-id-header | Request IDヘッダー名 | X-Request-ID |

#### タイミング制御
| オプション | 説明 | デフォルト |
|-----------|------|-----------|
| --delay | リクエスト間の遅延（例: 100ms, 1s） | 0 |
| --timeout | タイムアウト時間（例: 10s, 30s） | 30s |

#### 出力制御
| オプション | 短縮形 | 説明 | デフォルト |
|-----------|--------|------|-----------|
| --no-body | なし | レスポンスボディを非表示（JSON出力時は無視） | false |
| --json | なし | JSON形式で出力 | false |
| --output | -o | 結果をファイルに出力 | 標準出力 |

### 使用例

```bash
# 基本的な使用
conreq https://api.example.com/users

# 3つの同時POSTリクエスト
conreq -c 3 -X POST -d '{"name":"test"}' https://api.example.com/users

# ヘッダー付き、100ms間隔
conreq -c 3 -H "Authorization: Bearer token" --delay 100ms https://api.example.com

# 同一RequestIDで送信
conreq -c 3 --same-request-id https://api.example.com

# JSON出力をファイルに保存
conreq -c 5 --json -o result.json https://api.example.com
```

## 出力仕様

### 通常出力形式

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

### 通常出力形式（--no-body使用時）

```
=== Results ===
[1] 2024-01-20 15:30:45.123456 | Status: 200 | Time: 145ms | X-Request-ID: 550e8400-e29b-41d4-a716-446655440001
[Body omitted]

[2] 2024-01-20 15:30:45.234567 | Status: 200 | Time: 132ms | X-Request-ID: 550e8400-e29b-41d4-a716-446655440002
[Body omitted]
```

### JSON出力形式（--no-bodyオプションは無視される）

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
          "X-Request-ID": "550e8400-e29b-41d4-a716-446655440001",
          "Authorization": "Bearer token"
        },
        "body": "{\"name\":\"test user\"}"
      },
      "response": {
        "status_code": 200,
        "status_text": "OK",
        "headers": {
          "Content-Type": "application/json",
          "X-Response-ID": "resp-123"
        },
        "body": "{\"status\":\"ok\",\"data\":{\"id\":1,\"name\":\"test user\"}}"
      },
      "error": null
    },
    {
      "index": 2,
      "request_id": "550e8400-e29b-41d4-a716-446655440002",
      "started_at": "2024-01-20T15:30:45.234567Z",
      "completed_at": "2024-01-20T15:30:45.366567Z",
      "duration_ms": 132,
      "request": {
        "method": "POST",
        "url": "https://api.example.com/users",
        "headers": {
          "Content-Type": "application/json",
          "X-Request-ID": "550e8400-e29b-41d4-a716-446655440002",
          "Authorization": "Bearer token"
        },
        "body": "{\"name\":\"test user\"}"
      },
      "response": {
        "status_code": 200,
        "status_text": "OK",
        "headers": {
          "Content-Type": "application/json",
          "X-Response-ID": "resp-124"
        },
        "body": "{\"status\":\"ok\",\"data\":{\"id\":2,\"name\":\"test user\"}}"
      },
      "error": null
    },
    {
      "index": 3,
      "request_id": "550e8400-e29b-41d4-a716-446655440003",
      "started_at": "2024-01-20T15:30:45.345678Z",
      "completed_at": "2024-01-20T15:30:45.434678Z",
      "duration_ms": 89,
      "request": {
        "method": "POST",
        "url": "https://api.example.com/users",
        "headers": {
          "Content-Type": "application/json",
          "X-Request-ID": "550e8400-e29b-41d4-a716-446655440003",
          "Authorization": "Bearer token"
        },
        "body": "{\"name\":\"test user\"}"
      },
      "response": {
        "status_code": 500,
        "status_text": "Internal Server Error",
        "headers": {
          "Content-Type": "application/json"
        },
        "body": "{\"error\":\"internal server error\"}"
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

### エラー時のJSON出力例

接続エラーやタイムアウトが発生した場合：

```json
{
  "results": [
    {
      "index": 1,
      "request_id": "550e8400-e29b-41d4-a716-446655440001",
      "started_at": "2024-01-20T15:30:45.123456Z",
      "completed_at": "2024-01-20T15:31:15.123456Z",
      "duration_ms": 30000,
      "request": {
        "method": "GET",
        "url": "https://api.example.com/users",
        "headers": {
          "X-Request-ID": "550e8400-e29b-41d4-a716-446655440001"
        },
        "body": null
      },
      "response": null,
      "error": "request timeout: context deadline exceeded"
    }
  ]
}
```

## エラーハンドリング

### バリデーション
- URL形式の検証
- 同時実行数は1-5の範囲内
- タイムアウト・遅延の時間形式検証
- ファイル読み込みエラー

### 実行時エラー
- 接続エラー時も他のリクエストは継続
- 各リクエストのエラーは個別に記録
- タイムアウトエラーの明示

## 実装指針

### プロジェクト構成
```
conreq/
├── cmd/
│   └── main.go          # CLIエントリーポイント（cobraコマンド定義含む）
├── go.mod
├── go.sum
├── internal/
│   ├── client/
│   │   └── http.go      # HTTPクライアント実装
│   ├── runner/
│   │   └── concurrent.go # 並行実行ロジック
│   ├── output/
│   │   ├── formatter.go # 通常出力フォーマッター
│   │   └── json.go      # JSON出力フォーマッター
│   └── config/
│       └── config.go    # 設定構造体
├── pkg/
│   └── requestid/
│       └── generator.go # RequestID生成ロジック
└── README.md
```

### 推奨ライブラリ
- CLIフレームワーク: `github.com/spf13/cobra`
- UUID生成: `github.com/google/uuid`
- 時間パース: 標準ライブラリ `time.ParseDuration`
- HTTP Client: 標準ライブラリ `net/http`

### 実装の優先順位
1. 基本的なHTTPリクエスト機能
2. 同時実行機能
3. RequestID管理
4. 出力フォーマット
5. エラーハンドリングの改善

### テスト方針
- 単体テスト: 各パッケージごと
- 統合テスト: HTTPモックサーバーを使用
- E2Eテスト: 実際のテストAPIエンドポイントを用意

## 今後の拡張可能性（スコープ外）
- リクエスト数の上限拡張
- メトリクス詳細化（パーセンタイル等）
- リトライ機能
- 負荷テストモード
- WebSocket対応