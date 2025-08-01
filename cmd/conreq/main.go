// Package main provides the command-line interface for conreq.
package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/shiroemons/conreq/internal/config"
	"github.com/shiroemons/conreq/internal/output"
	"github.com/shiroemons/conreq/internal/runner"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "dev" && info.Main.Version != "" {
			version = info.Main.Version
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if commit == "none" && setting.Value != "" {
					commit = setting.Value
					if len(commit) > 7 {
						commit = commit[:7]
					}
				}
			case "vcs.time":
				if date == "unknown" && setting.Value != "" {
					date = setting.Value
				}
			}
		}
	}
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		method          string
		concurrent      int
		headers         []string
		data            string
		requestID       string
		sameRequestID   bool
		requestIDHeader string
		delay           string
		timeout         string
		noBody          bool
		outputJSON      bool
		outputFile      string
		showVersion     bool
		streamOutput    bool
	)

	cmd := &cobra.Command{
		Use:   "conreq [URL]",
		Short: "同一エンドポイントへの並行HTTPリクエストツール",
		Long: `conreqは、同一のAPIエンドポイントに対して複数の並行HTTPリクエストを送信し、
APIの挙動を検証するためのツールです。`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				versionStr := fmt.Sprintf("conreq version %s", version)
				if commit != "none" {
					versionStr += fmt.Sprintf(" (commit: %s", commit)
					if date != "unknown" {
						versionStr += fmt.Sprintf(", built at: %s", date)
					}
					versionStr += ")"
				}
				fmt.Println(versionStr)
				return nil
			}

			if len(args) == 0 {
				return cmd.Help()
			}

			// 設定を作成
			cfg := config.NewConfig()
			cfg.URL = args[0]
			cfg.Method = strings.ToUpper(method)
			cfg.Count = concurrent
			cfg.RequestID = requestID
			cfg.SameRequestID = sameRequestID
			cfg.RequestIDHeader = requestIDHeader
			cfg.OutputJSON = outputJSON
			cfg.NoBody = noBody

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
			if data != "" {
				if strings.HasPrefix(data, "@") {
					// ファイルから読み込み
					filename := data[1:]
					content, err := os.ReadFile(filename) //nolint:gosec // CLI argument
					if err != nil {
						return fmt.Errorf("ファイル読み込みエラー: %w", err)
					}
					cfg.Body = string(content)
				} else {
					cfg.Body = data
				}
			}

			// 設定を検証
			if err := cfg.Validate(); err != nil {
				return err
			}

			// リクエストを実行
			r := runner.NewRunner(cfg)

			// ストリーミング出力の設定（--streamフラグが有効で、JSON出力でない場合のみ）
			if streamOutput && !cfg.OutputJSON && outputFile == "" {
				progressFormatter := output.NewProgressFormatter(os.Stderr, cfg.Count)
				progressFormatter.Start()

				// プログレスチャネルを別goroutineで監視
				progressDone := make(chan struct{})
				go func() {
					for progress := range r.ProgressChannel() {
						progressFormatter.FormatProgress(progress)
					}
					progressFormatter.Finish()
					close(progressDone)
				}()

				// リクエストを実行
				result, err := r.Run(context.Background())
				if err != nil {
					return err
				}

				// プログレス出力の完了を待つ
				<-progressDone

				// 結果を標準出力に出力
				_, _ = fmt.Fprintln(os.Stdout, "\nFinal Results:")
				_, _ = fmt.Fprintln(os.Stdout, "")
				formatter := output.NewSpecTextFormatter(os.Stdout)
				return formatter.Format(result)
			}

			// JSON出力またはファイル出力の場合は従来通り
			result, err := r.Run(context.Background())
			if err != nil {
				return err
			}

			// 出力先を決定
			var outputWriter *os.File
			if outputFile != "" {
				file, err := os.Create(outputFile) //nolint:gosec // CLI argument
				if err != nil {
					return fmt.Errorf("出力ファイル作成エラー: %w", err)
				}
				defer func() { _ = file.Close() }()
				outputWriter = file
			} else {
				outputWriter = os.Stdout
			}

			// 結果を出力
			var formatter output.Formatter
			if cfg.OutputJSON {
				formatter = output.NewSpecJSONFormatter(outputWriter, cfg)
			} else {
				formatter = output.NewSpecTextFormatter(outputWriter)
			}

			return formatter.Format(result)
		},
	}

	cmd.Flags().StringVarP(&method, "method", "X", "GET", "HTTPメソッド (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)")
	cmd.Flags().IntVarP(&concurrent, "concurrent", "c", 1, "同時リクエスト数 (1-5)")
	cmd.Flags().StringArrayVarP(&headers, "header", "H", nil, "カスタムヘッダー (例: \"Content-Type: application/json\")")
	cmd.Flags().StringVarP(&data, "data", "d", "", "リクエストボディ (@でファイル指定可)")
	cmd.Flags().StringVar(&requestID, "request-id", "", "カスタムRequest ID値を指定")
	cmd.Flags().BoolVar(&sameRequestID, "same-request-id", false, "全リクエストで同一のRequest IDを使用")
	cmd.Flags().StringVar(&requestIDHeader, "request-id-header", "X-Request-ID", "Request IDヘッダー名")
	cmd.Flags().StringVar(&delay, "delay", "0s", "リクエスト間の遅延時間 (例: \"100ms\", \"1s\")")
	cmd.Flags().StringVar(&timeout, "timeout", "30s", "タイムアウト時間 (例: \"10s\", \"30s\")")
	cmd.Flags().BoolVar(&noBody, "no-body", false, "レスポンスボディを非表示（JSON出力時は無視）")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "JSON形式で出力")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "結果をファイルに出力")
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "バージョン情報を表示")
	cmd.Flags().BoolVar(&streamOutput, "stream", false, "リアルタイムで進行状況を表示")

	return cmd
}
