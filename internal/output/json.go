package output

import (
	"encoding/json"
	"io"
	"time"

	"github.com/shiroemons/conreq/internal/runner"
)

type JSONFormatter struct {
	writer io.Writer
}

func NewJSONFormatter(w io.Writer) *JSONFormatter {
	return &JSONFormatter{writer: w}
}

type JSONResponse struct {
	RequestID    string            `json:"request_id"`
	StatusCode   int               `json:"status_code,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body,omitempty"`
	Duration     string            `json:"duration"`
	Timestamp    string            `json:"timestamp"`
	RequestIndex int               `json:"request_index"`
	Error        string            `json:"error,omitempty"`
}

type JSONResult struct {
	StartTime    string          `json:"start_time"`
	EndTime      string          `json:"end_time"`
	TotalTime    string          `json:"total_time"`
	RequestCount int             `json:"request_count"`
	SuccessCount int             `json:"success_count"`
	ErrorCount   int             `json:"error_count"`
	AvgDuration  string          `json:"average_duration,omitempty"`
	Responses    []JSONResponse  `json:"responses"`
}

func (f *JSONFormatter) Format(result *runner.Result) error {
	jsonResult := JSONResult{
		StartTime:    result.StartTime.Format(time.RFC3339),
		EndTime:      result.EndTime.Format(time.RFC3339),
		TotalTime:    result.EndTime.Sub(result.StartTime).String(),
		RequestCount: len(result.Responses),
		SuccessCount: result.SuccessCount(),
		ErrorCount:   result.ErrorCount(),
		Responses:    make([]JSONResponse, 0, len(result.Responses)),
	}

	if result.SuccessCount() > 0 {
		jsonResult.AvgDuration = result.AverageDuration().String()
	}

	for _, resp := range result.Responses {
		jsonResp := JSONResponse{
			RequestID:    resp.RequestID,
			Duration:     resp.Duration.String(),
			Timestamp:    resp.Timestamp.Format(time.RFC3339Nano),
			RequestIndex: resp.RequestIndex,
		}

		if resp.Error != nil {
			jsonResp.Error = resp.Error.Error()
		} else {
			jsonResp.StatusCode = resp.StatusCode
			jsonResp.Body = resp.Body
			
			jsonResp.Headers = make(map[string]string)
			for key, values := range resp.Headers {
				if len(values) > 0 {
					jsonResp.Headers[key] = values[0]
				}
			}
		}

		jsonResult.Responses = append(jsonResult.Responses, jsonResp)
	}

	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jsonResult)
}