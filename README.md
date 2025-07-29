# conreq

conreqは、同一のAPIエンドポイントに対して複数の並行HTTPリクエストを送信し、APIの挙動を検証するためのCLIツールです。

## 特徴

- 1〜5個の並行HTTPリクエストを送信
- X-Request-IDヘッダーの自動生成またはカスタム値の設定
- リクエスト間の遅延時間設定
- タイムアウト制御
- テキストまたはJSON形式での結果出力
- ファイルからのリクエストボディ読み込み
- 全HTTPメソッドのサポート（GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS）

## インストール

### Go 1.24以降がインストールされている場合

```bash
go install github.com/shiroemons/conreq/cmd@latest
```

### ソースからビルド

```bash
git clone https://github.com/shiroemons/conreq.git
cd conreq
go build -o conreq cmd/main.go
```

## 使い方

### 基本的な使用方法

```bash
# 単一のGETリクエスト
conreq https://api.example.com/endpoint

# 3つの並行GETリクエスト
conreq https://api.example.com/endpoint -n 3

# POSTリクエストでJSONボディを送信
conreq https://api.example.com/endpoint -X POST -d '{"key": "value"}' -H "Content-Type: application/json"

# ファイルからボディを読み込んでPOSTリクエスト
conreq https://api.example.com/endpoint -X POST -f request_body.json -H "Content-Type: application/json"
```

### コマンドラインオプション

| オプション | 短縮形 | 説明 | デフォルト |
|-----------|--------|------|------------|
| `--method` | `-X` | HTTPメソッド (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS) | GET |
| `--count` | `-n` | 同時リクエスト数 (1-5) | 1 |
| `--header` | `-H` | カスタムヘッダー（複数指定可） | なし |
| `--data` | `-d` | リクエストボディ | なし |
| `--data-file` | `-f` | リクエストボディを含むファイルパス | なし |
| `--request-id` | `-r` | X-Request-IDヘッダーの値 | 自動生成 |
| `--delay` | | リクエスト間の遅延時間 | 0s |
| `--timeout` | `-t` | リクエストタイムアウト | 30s |
| `--json` | `-j` | JSON形式で出力 | false |
| `--version` | `-v` | バージョン情報を表示 | - |
| `--help` | `-h` | ヘルプを表示 | - |

### 出力例

#### テキスト形式（デフォルト）

```
=== 実行結果 ===
開始時刻: 2025-07-29T11:49:25+09:00
終了時刻: 2025-07-29T11:49:26+09:00
総実行時間: 860.596291ms
リクエスト数: 3
成功数: 3
エラー数: 0
平均レスポンス時間: 833.091403ms

=== リクエスト詳細 ===
No  Request ID                            Status  Duration      Time
1   550e8400-e29b-41d4-a716-446655440000  200     860.223292ms  11:49:25.631
2   6ba7b810-9dad-11d1-80b4-00c04fd430c8  200     819.191417ms  11:49:25.631
3   6ba7b814-9dad-11d1-80b4-00c04fd430c8  200     819.8595ms    11:49:25.632
```

#### JSON形式

```json
{
  "start_time": "2025-07-29T11:49:31+09:00",
  "end_time": "2025-07-29T11:49:32+09:00",
  "total_time": "1.124951666s",
  "request_count": 2,
  "success_count": 2,
  "error_count": 0,
  "average_duration": "965.623729ms",
  "responses": [
    {
      "request_id": "550e8400-e29b-41d4-a716-446655440000",
      "status_code": 200,
      "headers": {
        "Content-Type": "application/json",
        "Content-Length": "443"
      },
      "body": "{...}",
      "duration": "806.336542ms",
      "timestamp": "2025-07-29T11:49:31.54474+09:00",
      "request_index": 0
    }
  ]
}
```

## 使用例

### API負荷テスト

```bash
# 5つの並行リクエストで負荷をかける
conreq https://api.example.com/heavy-endpoint -n 5 -t 60s
```

### レート制限の検証

```bash
# 100ms間隔で5つのリクエストを送信
conreq https://api.example.com/rate-limited -n 5 --delay 100ms
```

### 冪等性の確認

```bash
# 同じX-Request-IDで複数のリクエストを送信
conreq https://api.example.com/idempotent -n 3 -r "fixed-request-id"
```

## 開発

### 必要なもの

- Go 1.24以降

### ビルド

```bash
go build -o conreq cmd/main.go
```

### テスト

```bash
go test ./...
```

### プロジェクト構造

```
conreq/
├── cmd/
│   └── main.go          # CLIエントリーポイント
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
