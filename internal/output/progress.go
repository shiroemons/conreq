// Package output provides formatting and display functionality for request results.
package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/shiroemons/conreq/internal/runner"
)

// ProgressFormatter formats progress updates for streaming output.
type ProgressFormatter struct {
	writer       io.Writer
	startTime    time.Time
	totalCount   int
	requestWidth int
}

// NewProgressFormatter creates a new progress formatter.
func NewProgressFormatter(w io.Writer, totalCount int) *ProgressFormatter {
	// Calculate the width needed for request index display
	requestWidth := len(fmt.Sprintf("%d", totalCount-1))
	if requestWidth < 7 { // minimum width for "Request"
		requestWidth = 7
	}

	return &ProgressFormatter{
		writer:       w,
		startTime:    time.Now(),
		totalCount:   totalCount,
		requestWidth: requestWidth,
	}
}

// Start prints the initial header.
func (f *ProgressFormatter) Start() {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, _ = fmt.Fprintf(f.writer, "🚀 Starting %d concurrent requests at %s\n\n", f.totalCount, now)
	f.printHeader()
}

// FormatProgress formats a single progress update.
func (f *ProgressFormatter) FormatProgress(p *runner.Progress) {
	elapsed := time.Since(f.startTime)
	timeStr := time.Now().Format("15:04:05.000")
	requestStr := fmt.Sprintf("Request %*d", f.requestWidth-7, p.Index)

	var statusIcon, statusText, httpCode, duration string
	switch p.Status {
	case "pending":
		statusIcon = "⏳"
		statusText = "PENDING"
		httpCode = "-"
		duration = "-"
	case "running":
		statusIcon = "🔄"
		statusText = "RUNNING"
		httpCode = "-"
		duration = "-"
	case "completed":
		statusIcon = "✅"
		if p.StatusCode >= 400 {
			statusIcon = "⚠️"
		}
		if p.StatusCode >= 500 {
			statusIcon = "❌"
		}
		statusText = "DONE"
		httpCode = fmt.Sprintf("%d", p.StatusCode)
		duration = formatDuration(p.EndTime.Sub(p.StartTime))
	case "failed":
		statusIcon = "❌"
		statusText = "FAILED"
		httpCode = "-"
		duration = "-"
		if p.Error != nil {
			statusText = fmt.Sprintf("FAILED: %s", p.Error.Error())
		}
	}

	_, _ = fmt.Fprintf(f.writer, "[%8s] %s | %-*s  %s %-7s  %4s  %9s  %s\n",
		formatDuration(elapsed),
		timeStr,
		f.requestWidth, requestStr,
		statusIcon,
		statusText,
		httpCode,
		duration,
		p.RequestID,
	)
}

// Finish prints the completion message.
func (f *ProgressFormatter) Finish() {
	elapsed := time.Since(f.startTime)
	now := time.Now().Format("2006-01-02 15:04:05")
	_, _ = fmt.Fprintf(f.writer, "\n🎉 All requests completed in %s at %s\n", formatDuration(elapsed), now)
	_, _ = fmt.Fprintln(f.writer, strings.Repeat("=", 110))
}

func (f *ProgressFormatter) printHeader() {
	_, _ = fmt.Fprintf(f.writer, "%-10s %-12s | %-*s  %s %-7s  %4s  %9s  %s\n",
		"+Time", "Time", f.requestWidth, "Request", " ", "Status", "Code", "Duration", "Request-ID")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat("─", 110))
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}