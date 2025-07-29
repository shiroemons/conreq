package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/shiroemons/conreq/internal/config"
	"github.com/shiroemons/conreq/internal/output"
	"github.com/shiroemons/conreq/internal/runner"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		method      string
		count       int
		headers     []string
		body        string
		bodyFile    string
		requestID   string
		delay       string
		timeout     string
		outputJSON  bool
		showVersion bool
	)

	cmd := &cobra.Command{
		Use:   "conreq [URL]",
		Short: "同一エンドポイントへの並行HTTPリクエストツール",
		Long: `conreqは、同一のAPIエンドポイントに対して複数の並行HTTPリクエストを送信し、
APIの挙動を検証するためのツールです。`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Printf("conreq version %s (commit: %s, built at: %s)\n", version, commit, date)
				return nil
			}

			if len(args) == 0 {
				return cmd.Help()
			}

			// 設定を作成
			cfg := config.NewConfig()
			cfg.URL = args[0]
			cfg.Method = strings.ToUpper(method)
			cfg.Count = count
			cfg.RequestID = requestID
			cfg.OutputJSON = outputJSON

			// ヘッダーをパース
			if err := cfg.ParseHeaders(headers); err != nil {
				return err
			}

			// タイムアウトをパース
			timeoutDuration, err := config.ParseDuration(timeout)
			if err != nil {
				return fmt.Errorf("無効なタイムアウト形式: %w", err)
			}
			cfg.Timeout = timeoutDuration

			// 遅延時間をパース
			delayDuration, err := config.ParseDuration(delay)
			if err != nil {
				return fmt.Errorf("無効な遅延時間形式: %w", err)
			}
			cfg.Delay = delayDuration

			// リクエストボディの設定
			if body != "" && bodyFile != "" {
				return fmt.Errorf("--data と --data-file は同時に指定できません")
			}

			if bodyFile != "" {
				content, err := os.ReadFile(bodyFile)
				if err != nil {
					return fmt.Errorf("ファイル読み込みエラー: %w", err)
				}
				cfg.Body = string(content)
			} else {
				cfg.Body = body
			}

			// 設定を検証
			if err := cfg.Validate(); err != nil {
				return err
			}

			// リクエストを実行
			r := runner.NewRunner(cfg)
			result, err := r.Run(context.Background())
			if err != nil {
				return err
			}

			// 結果を出力
			var formatter output.Formatter
			if cfg.OutputJSON {
				formatter = output.NewJSONFormatter(os.Stdout)
			} else {
				formatter = output.NewTextFormatter(os.Stdout)
			}

			return formatter.Format(result)
		},
	}

	cmd.Flags().StringVarP(&method, "method", "X", "GET", "HTTPメソッド (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)")
	cmd.Flags().IntVarP(&count, "count", "n", 1, "同時リクエスト数 (1-5)")
	cmd.Flags().StringArrayVarP(&headers, "header", "H", nil, "カスタムヘッダー (例: \"Content-Type: application/json\")")
	cmd.Flags().StringVarP(&body, "data", "d", "", "リクエストボディ")
	cmd.Flags().StringVarP(&bodyFile, "data-file", "f", "", "リクエストボディを含むファイルパス")
	cmd.Flags().StringVarP(&requestID, "request-id", "r", "", "X-Request-IDヘッダーの値 (デフォルト: 自動生成)")
	cmd.Flags().StringVarP(&delay, "delay", "", "0s", "リクエスト間の遅延時間 (例: \"100ms\", \"1s\")")
	cmd.Flags().StringVarP(&timeout, "timeout", "t", "30s", "リクエストタイムアウト (例: \"10s\", \"1m\")")
	cmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "JSON形式で出力")
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "バージョン情報を表示")

	return cmd
}