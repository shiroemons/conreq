package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/shiroemons/conreq/internal/client"
	"github.com/shiroemons/conreq/internal/runner"
)

// SpecTextFormatter formats results as plain text according to the specification.
type SpecTextFormatter struct {
	writer io.Writer
	config interface{} // configは後で渡せるように
}

// NewSpecTextFormatter creates a new spec text formatter.
func NewSpecTextFormatter(w io.Writer) *SpecTextFormatter {
	return &SpecTextFormatter{writer: w}
}

// SetConfig sets the configuration for the formatter.
func (f *SpecTextFormatter) SetConfig(cfg interface{}) {
	f.config = cfg
}

// Format formats the result according to the specification.
//
//nolint:errcheck // io.Writer への出力エラーは無視
func (f *SpecTextFormatter) Format(result *runner.Result) error {
	// Request Summary
	fmt.Fprintln(f.writer, "=== Request Summary ===")
	fmt.Fprintf(f.writer, "URL: %s\n", result.Config.URL)
	fmt.Fprintf(f.writer, "Method: %s\n", result.Config.Method)
	fmt.Fprintf(f.writer, "Concurrent: %d\n", result.Config.Count)
	fmt.Fprintf(f.writer, "Total Requests: %d\n", len(result.Responses))

	// Results
	fmt.Fprintln(f.writer, "\n=== Results ===")

	// インデックス順にソート
	sortedResponses := make([]*client.Response, len(result.Responses))
	copy(sortedResponses, result.Responses)
	sort.Slice(sortedResponses, func(i, j int) bool {
		return sortedResponses[i].RequestIndex < sortedResponses[j].RequestIndex
	})

	for _, resp := range sortedResponses {
		index := resp.RequestIndex + 1
		timestamp := resp.Timestamp.Format("2006-01-02 15:04:05.000000")

		if resp.Error != nil {
			// エラーの場合
			fmt.Fprintf(f.writer, "[%d] %s | Status: ERROR | Time: %dms | %s: %s\n",
				index,
				timestamp,
				resp.Duration.Milliseconds(),
				result.Config.RequestIDHeader,
				resp.RequestID,
			)
			fmt.Fprintf(f.writer, "Error: %v\n", resp.Error)
		} else {
			// 成功の場合
			fmt.Fprintf(f.writer, "[%d] %s | Status: %d | Time: %dms | %s: %s\n",
				index,
				timestamp,
				resp.StatusCode,
				resp.Duration.Milliseconds(),
				result.Config.RequestIDHeader,
				resp.RequestID,
			)

			// レスポンスボディ
			if !result.Config.NoBody {
				fmt.Fprintln(f.writer, resp.Body)
			} else {
				fmt.Fprintln(f.writer, "[Body omitted]")
			}
		}

		if index < len(sortedResponses) {
			fmt.Fprintln(f.writer)
		}
	}

	// Summary
	fmt.Fprintln(f.writer, "\n=== Summary ===")

	successCount := result.SuccessCount()
	errorCount := result.ErrorCount()
	total := len(result.Responses)

	successRate := 0.0
	if total > 0 {
		successRate = float64(successCount) / float64(total) * 100
	}

	fmt.Fprintf(f.writer, "Success: %d/%d (%.1f%%)\n", successCount, total, successRate)
	
	// Status code breakdown
	count2xx := result.Count2xx()
	count3xx := result.Count3xx()
	count4xx := result.Count4xx()
	count5xx := result.Count5xx()
	
	fmt.Fprintln(f.writer, "\n=== Status Code Breakdown ===")
	if count2xx > 0 {
		fmt.Fprintf(f.writer, "2xx (Success): %d\n", count2xx)
	}
	if count3xx > 0 {
		fmt.Fprintf(f.writer, "3xx (Redirect): %d\n", count3xx)
	}
	if count4xx > 0 {
		fmt.Fprintf(f.writer, "4xx (Client Error): %d\n", count4xx)
	}
	if count5xx > 0 {
		fmt.Fprintf(f.writer, "5xx (Server Error): %d\n", count5xx)
	}
	if errorCount > 0 {
		fmt.Fprintf(f.writer, "Network/Timeout Errors: %d\n", errorCount)
	}

	if errorCount > 0 {
		errorRate := float64(errorCount) / float64(total) * 100
		fmt.Fprintf(f.writer, "Failed: %d/%d (%.1f%%)\n", errorCount, total, errorRate)
	}

	if successCount > 0 {
		avgDuration := result.AverageDuration()
		fmt.Fprintf(f.writer, "Average Response Time: %dms\n", avgDuration.Milliseconds())
	}

	return nil
}
