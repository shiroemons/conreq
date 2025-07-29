package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("conreq version %s (commit: %s, built at: %s)\n", version, commit, date)
				return
			}

			if len(args) == 0 {
				cmd.Help()
				return
			}

			// TODO: 実装
			fmt.Println("実装予定")
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