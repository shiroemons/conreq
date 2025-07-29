package output

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/shiroemons/conreq/internal/client"
	"github.com/shiroemons/conreq/internal/config"
	"github.com/shiroemons/conreq/internal/runner"
)

// SpecJSONFormatter formats results as JSON according to the specification.
type SpecJSONFormatter struct {
	writer io.Writer
	config *config.Config
}

// NewSpecJSONFormatter creates a new spec JSON formatter.
func NewSpecJSONFormatter(w io.Writer, cfg *config.Config) *SpecJSONFormatter {
	return &SpecJSONFormatter{writer: w, config: cfg}
}

// SpecJSONMetadata represents metadata in the JSON output.
type SpecJSONMetadata struct {
	URL             string `json:"url"`
	Method          string `json:"method"`
	Concurrent      int    `json:"concurrent"`
	TotalRequests   int    `json:"total_requests"`
	StartedAt       string `json:"started_at"`
	CompletedAt     string `json:"completed_at"`
	TotalDurationMs int64  `json:"total_duration_ms"`
}

// SpecJSONRequest represents a request in the JSON output.
type SpecJSONRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

// SpecJSONResponse represents a response in the JSON output.
type SpecJSONResponse struct {
	StatusCode int               `json:"status_code"`
	StatusText string            `json:"status_text"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// SpecJSONResult represents a single result in the JSON output.
type SpecJSONResult struct {
	Index       int               `json:"index"`
	RequestID   string            `json:"request_id"`
	StartedAt   string            `json:"started_at"`
	CompletedAt string            `json:"completed_at"`
	DurationMs  int64             `json:"duration_ms"`
	Request     SpecJSONRequest   `json:"request"`
	Response    *SpecJSONResponse `json:"response"`
	Error       interface{}       `json:"error"`
}

// SpecJSONSummary represents the summary in the JSON output.
type SpecJSONSummary struct {
	Total               int            `json:"total"`
	Successful          int            `json:"successful"`
	Failed              int            `json:"failed"`
	SuccessRate         float64        `json:"success_rate"`
	AverageDurationMs   int64          `json:"average_duration_ms"`
	MinDurationMs       int64          `json:"min_duration_ms"`
	MaxDurationMs       int64          `json:"max_duration_ms"`
	StatusCodes         map[string]int `json:"status_codes"`
	StatusCodeBreakdown struct {
		Count2xx      int `json:"2xx"`
		Count3xx      int `json:"3xx"`
		Count4xx      int `json:"4xx"`
		Count5xx      int `json:"5xx"`
		NetworkErrors int `json:"network_errors"`
	} `json:"status_code_breakdown"`
}

// SpecJSONOutput represents the complete JSON output structure.
type SpecJSONOutput struct {
	Metadata SpecJSONMetadata `json:"metadata"`
	Results  []SpecJSONResult `json:"results"`
	Summary  SpecJSONSummary  `json:"summary"`
}

// Format formats the result as JSON according to the specification.
//
//nolint:funlen // 仕様に従った複雑な出力のため
func (f *SpecJSONFormatter) Format(result *runner.Result) error {
	output := SpecJSONOutput{
		Metadata: SpecJSONMetadata{
			URL:             f.config.URL,
			Method:          f.config.Method,
			Concurrent:      f.config.Count,
			TotalRequests:   len(result.Responses),
			StartedAt:       result.StartTime.Format(time.RFC3339Nano),
			CompletedAt:     result.EndTime.Format(time.RFC3339Nano),
			TotalDurationMs: result.EndTime.Sub(result.StartTime).Milliseconds(),
		},
		Results: make([]SpecJSONResult, 0, len(result.Responses)),
	}

	// ソートしてインデックス順に処理
	sortedResponses := make([]*client.Response, len(result.Responses))
	copy(sortedResponses, result.Responses)
	sort.Slice(sortedResponses, func(i, j int) bool {
		return sortedResponses[i].RequestIndex < sortedResponses[j].RequestIndex
	})

	statusCodes := make(map[string]int)
	var totalDuration, minDuration, maxDuration int64
	successCount := 0
	firstSuccess := true

	for _, resp := range sortedResponses {
		headers := make(map[string]string)
		for key, value := range f.config.Headers {
			headers[key] = value
		}
		if resp.RequestID != "" {
			headers[f.config.RequestIDHeader] = resp.RequestID
		}
		if f.config.Body != "" && headers["Content-Type"] == "" {
			headers["Content-Type"] = "application/json"
		}

		var body interface{}
		if f.config.Body != "" {
			body = f.config.Body
		} else {
			body = nil
		}

		result := SpecJSONResult{
			Index:       resp.RequestIndex + 1,
			RequestID:   resp.RequestID,
			StartedAt:   resp.Timestamp.Format(time.RFC3339Nano),
			CompletedAt: resp.Timestamp.Add(resp.Duration).Format(time.RFC3339Nano),
			DurationMs:  resp.Duration.Milliseconds(),
			Request: SpecJSONRequest{
				Method:  f.config.Method,
				URL:     f.config.URL,
				Headers: headers,
				Body:    body,
			},
			Response: nil,
			Error:    nil,
		}

		if resp.Error != nil {
			// エラーの場合
			if resp.Error.Error() == "context deadline exceeded" {
				result.Error = "request timeout: context deadline exceeded"
			} else {
				result.Error = resp.Error.Error()
			}
		} else {
			// 成功の場合
			statusText := http.StatusText(resp.StatusCode)
			if statusText == "" {
				statusText = "Unknown"
			}

			respHeaders := make(map[string]string)
			for key, values := range resp.Headers {
				if len(values) > 0 {
					respHeaders[key] = values[0]
				}
			}

			result.Response = &SpecJSONResponse{
				StatusCode: resp.StatusCode,
				StatusText: statusText,
				Headers:    respHeaders,
				Body:       resp.Body,
			}

			statusCode := resp.StatusCode
			statusCodeStr := fmt.Sprintf("%d", statusCode)
			statusCodes[statusCodeStr]++
			successCount++

			durationMs := resp.Duration.Milliseconds()
			totalDuration += durationMs

			if firstSuccess {
				minDuration = durationMs
				maxDuration = durationMs
				firstSuccess = false
			} else {
				if durationMs < minDuration {
					minDuration = durationMs
				}
				if durationMs > maxDuration {
					maxDuration = durationMs
				}
			}
		}

		output.Results = append(output.Results, result)
	}

	// サマリーを計算
	total := len(result.Responses)
	failed := total - successCount
	successRate := 0.0
	avgDuration := int64(0)

	if total > 0 {
		successRate = float64(successCount) / float64(total) * 100
	}
	if successCount > 0 {
		avgDuration = totalDuration / int64(successCount)
	}

	output.Summary = SpecJSONSummary{
		Total:             total,
		Successful:        successCount,
		Failed:            failed,
		SuccessRate:       successRate,
		AverageDurationMs: avgDuration,
		MinDurationMs:     minDuration,
		MaxDurationMs:     maxDuration,
		StatusCodes:       statusCodes,
	}

	// Status code breakdown
	output.Summary.StatusCodeBreakdown.Count2xx = result.Count2xx()
	output.Summary.StatusCodeBreakdown.Count3xx = result.Count3xx()
	output.Summary.StatusCodeBreakdown.Count4xx = result.Count4xx()
	output.Summary.StatusCodeBreakdown.Count5xx = result.Count5xx()
	output.Summary.StatusCodeBreakdown.NetworkErrors = result.ErrorCount()

	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
