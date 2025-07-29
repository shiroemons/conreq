package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shiroemons/conreq/internal/runner"
)

func TestProgressFormatter(t *testing.T) {
	tests := []struct {
		name       string
		totalCount int
		progress   *runner.Progress
		wantOutput []string // parts that should be in the output
	}{
		{
			name:       "pending status",
			totalCount: 3,
			progress: &runner.Progress{
				Index:     0,
				RequestID: "test-id-123",
				Status:    "pending",
				StartTime: time.Now(),
			},
			wantOutput: []string{"Request 0", "⏳", "PENDING", "test-id-123"},
		},
		{
			name:       "running status",
			totalCount: 5,
			progress: &runner.Progress{
				Index:     2,
				RequestID: "test-id-456",
				Status:    "running",
				StartTime: time.Now(),
			},
			wantOutput: []string{"Request 2", "🔄", "RUNNING", "test-id-456"},
		},
		{
			name:       "completed status",
			totalCount: 2,
			progress: &runner.Progress{
				Index:      1,
				RequestID:  "test-id-789",
				Status:     "completed",
				StatusCode: 200,
				StartTime:  time.Now(),
				EndTime:    time.Now().Add(1500 * time.Millisecond),
			},
			wantOutput: []string{"Request 1", "✅", "DONE", "200", "test-id-789"},
		},
		{
			name:       "failed status",
			totalCount: 1,
			progress: &runner.Progress{
				Index:     0,
				RequestID: "test-id-error",
				Status:    "failed",
				Error:     &testError{msg: "connection timeout"},
				StartTime: time.Now(),
			},
			wantOutput: []string{"Request 0", "❌", "FAILED", "connection timeout", "test-id-error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			formatter := NewProgressFormatter(&buf, tt.totalCount)
			formatter.FormatProgress(tt.progress)

			output := buf.String()
			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("output does not contain %q\nGot: %s", want, output)
				}
			}
		})
	}
}

func TestProgressFormatterStart(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewProgressFormatter(&buf, 5)
	formatter.Start()

	output := buf.String()
	if !strings.Contains(output, "Starting 5 concurrent requests at") {
		t.Errorf("Start() output missing expected message\nGot: %s", output)
	}
	if !strings.Contains(output, "+Time") && !strings.Contains(output, "Time") && !strings.Contains(output, "Request") && !strings.Contains(output, "Status") {
		t.Errorf("Start() output missing header\nGot: %s", output)
	}
}

func TestProgressFormatterFinish(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewProgressFormatter(&buf, 3)
	formatter.Finish()

	output := buf.String()
	if !strings.Contains(output, "All requests completed in") || !strings.Contains(output, " at ") {
		t.Errorf("Finish() output missing completion message\nGot: %s", output)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{100 * time.Microsecond, "100µs"},
		{500 * time.Microsecond, "500µs"},
		{1 * time.Millisecond, "1ms"},
		{500 * time.Millisecond, "500ms"},
		{1 * time.Second, "1.00s"},
		{2500 * time.Millisecond, "2.50s"},
	}

	for _, tt := range tests {
		got := formatDuration(tt.duration)
		if got != tt.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
		}
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}