// Package output provides formatters for displaying request results.
package output

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/shiroemons/conreq/internal/client"
	"github.com/shiroemons/conreq/internal/runner"
)

// Formatter is the interface for result formatters.
type Formatter interface {
	Format(result *runner.Result) error
}

// TextFormatter formats results as plain text.
type TextFormatter struct {
	writer io.Writer
}

// NewTextFormatter creates a new text formatter.
func NewTextFormatter(w io.Writer) *TextFormatter {
	return &TextFormatter{writer: w}
}

// Format formats the result as plain text.
//nolint:errcheck // io.Writer への出力エラーは無視
func (f *TextFormatter) Format(result *runner.Result) error {
	fmt.Fprintln(f.writer, "=== 実行結果 ===")
	fmt.Fprintf(f.writer, "開始時刻: %s\n", result.StartTime.Format(time.RFC3339))
	fmt.Fprintf(f.writer, "終了時刻: %s\n", result.EndTime.Format(time.RFC3339))
	fmt.Fprintf(f.writer, "総実行時間: %v\n", result.EndTime.Sub(result.StartTime))
	fmt.Fprintf(f.writer, "リクエスト数: %d\n", len(result.Responses))
	fmt.Fprintf(f.writer, "成功数: %d\n", result.SuccessCount())
	fmt.Fprintf(f.writer, "エラー数: %d\n", result.ErrorCount())

	if result.SuccessCount() > 0 {
		fmt.Fprintf(f.writer, "平均レスポンス時間: %v\n", result.AverageDuration())
	}

	fmt.Fprintln(f.writer, "\n=== リクエスト詳細 ===")

	sortedResponses := make([]*client.Response, len(result.Responses))
	copy(sortedResponses, result.Responses)
	sort.Slice(sortedResponses, func(i, j int) bool {
		return sortedResponses[i].Timestamp.Before(sortedResponses[j].Timestamp)
	})

	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "No\tRequest ID\tStatus\tDuration\tTime")

	for i, resp := range sortedResponses {
		status := fmt.Sprintf("%d", resp.StatusCode)
		if resp.Error != nil {
			status = "ERROR"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%v\t%s\n",
			i+1,
			truncateString(resp.RequestID, 36),
			status,
			resp.Duration,
			resp.Timestamp.Format("15:04:05.000"),
		)
	}

	if err := w.Flush(); err != nil {
		return err
	}

	if result.HasErrors() {
		fmt.Fprintln(f.writer, "\n=== エラー詳細 ===")
		for i, resp := range sortedResponses {
			if resp.Error != nil {
				fmt.Fprintf(f.writer, "リクエスト %d: %v\n", i+1, resp.Error)
			}
		}
	}

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
