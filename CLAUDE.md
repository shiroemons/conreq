# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

**conreq** は、同一APIエンドポイントに対して複数の並行HTTPリクエストを送信し、APIの挙動を検証するためのCLIツールです。

## プロジェクトの現状

このプロジェクトは初期開発段階にあり、詳細な仕様書（`docs/conreq-spec.md`）は存在しますが、実装はまだ開始されていません。

## アーキテクチャとディレクトリ構造

```
conreq/
├── cmd/
│   └── conreq/
│       └── main.go      # CLIエントリーポイント（Cobra使用）
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
└── pkg/
    └── requestid/
        └── generator.go # RequestID生成ロジック
```

## 主要機能

1. **並行リクエスト実行**: 1-5個の同時HTTPリクエスト送信
2. **X-Request-IDヘッダー管理**: UUID v4自動生成またはカスタム値
3. **出力形式**: 人間が読みやすい形式とJSON形式
4. **タイミング制御**: リクエスト間遅延とタイムアウト
5. **ファイルサポート**: ファイルからのリクエストボディ読み込み

## 開発用コマンド

```bash
# ビルド
go build -o conreq cmd/conreq/main.go

# テスト実行
go test ./...

# 単一パッケージのテスト
go test ./internal/client

# lintチェック（golangci-lintが必要）
golangci-lint run

# コードフォーマット
go fmt ./...

# 依存関係の整理
go mod tidy
```

## 実装時の重要事項

1. **HTTPクライアント**: 
   - デフォルトタイムアウト: 30秒
   - リダイレクト無効化（CheckRedirect使用）
   - カスタムヘッダーサポート

2. **並行実行**:
   - sync.WaitGroupを使用した同期
   - エラーチャネルによるエラー収集
   - 実行順序の記録（タイムスタンプ付き）

3. **出力処理**:
   - 通常出力: tabwriterによる整形
   - JSON出力: 構造化されたレスポンス情報

4. **エラーハンドリング**:
   - HTTPエラーとタイムアウトの区別
   - 並行実行時の部分的失敗への対応

## 依存関係

```go
// 外部依存
github.com/spf13/cobra     // CLIフレームワーク
github.com/google/uuid      // UUID生成

// 標準ライブラリ
net/http                    // HTTPクライアント
sync                        // 並行処理
time                        // タイミング制御
encoding/json               // JSON処理
```

## 仕様書

詳細な仕様は `docs/conreq-spec.md` に記載されています。このファイルには以下が含まれます：
- 全コマンドラインオプションの詳細
- 各HTTPメソッドの動作仕様
- 出力フォーマットの例
- エラーケースの処理方法