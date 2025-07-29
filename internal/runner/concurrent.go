// Package runner provides functionality for executing concurrent HTTP requests.
package runner

import (
	"context"
	"sync"
	"time"

	"github.com/shiroemons/conreq/internal/client"
	"github.com/shiroemons/conreq/internal/config"
	"github.com/shiroemons/conreq/pkg/requestid"
)

// Progress represents the progress of a single request.
type Progress struct {
	Index      int
	RequestID  string
	Status     string // "pending", "running", "completed", "failed"
	StatusCode int
	Error      error
	StartTime  time.Time
	EndTime    time.Time
}

// Result represents the result of concurrent HTTP requests.
type Result struct {
	Responses []*client.Response
	StartTime time.Time
	EndTime   time.Time
	Config    *config.Config
}

// Runner executes concurrent HTTP requests.
type Runner struct {
	config       *config.Config
	client       *client.Client
	progressChan chan *Progress
}

// NewRunner creates a new Runner.
func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		config:       cfg,
		client:       client.NewClient(cfg),
		progressChan: make(chan *Progress, cfg.Count*3), // buffer for pending, running, completed
	}
}

// ProgressChannel returns the progress channel for streaming updates.
func (r *Runner) ProgressChannel() <-chan *Progress {
	return r.progressChan
}

// Run executes concurrent HTTP requests.
func (r *Runner) Run(ctx context.Context) (*Result, error) {
	result := &Result{
		StartTime: time.Now(),
		Responses: make([]*client.Response, 0, r.config.Count),
		Config:    r.config,
	}

	responseChan := make(chan *client.Response, r.config.Count)
	var wg sync.WaitGroup

	// 同一RequestIDモードの場合、事前に生成
	var sharedRequestID string
	if r.config.SameRequestID && r.config.RequestID == "" {
		sharedRequestID = requestid.Generate()
	}

	for i := 0; i < r.config.Count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			cfg := *r.config
			if r.config.SameRequestID {
				// 同一RequestIDモード
				if r.config.RequestID != "" {
					cfg.RequestID = r.config.RequestID
				} else {
					cfg.RequestID = sharedRequestID
				}
			} else {
				// 個別RequestIDモード
				if cfg.RequestID == "" {
					cfg.RequestID = requestid.Generate()
				}
			}

			// Send pending status
			r.progressChan <- &Progress{
				Index:     index,
				RequestID: cfg.RequestID,
				Status:    "pending",
				StartTime: time.Now(),
			}

			client := client.NewClient(&cfg)

			delay := time.Duration(index) * r.config.Delay
			if delay > 0 {
				time.Sleep(delay)
			}

			// Send running status
			startTime := time.Now()
			r.progressChan <- &Progress{
				Index:     index,
				RequestID: cfg.RequestID,
				Status:    "running",
				StartTime: startTime,
			}

			response := client.Do(ctx, index)

			// Send completed/failed status
			endTime := time.Now()
			status := "completed"
			if response.Error != nil {
				status = "failed"
			}
			r.progressChan <- &Progress{
				Index:      index,
				RequestID:  cfg.RequestID,
				Status:     status,
				StatusCode: response.StatusCode,
				Error:      response.Error,
				StartTime:  startTime,
				EndTime:    endTime,
			}

			responseChan <- response
		}(i)
	}

	go func() {
		wg.Wait()
		close(responseChan)
		close(r.progressChan)
	}()

	for response := range responseChan {
		result.Responses = append(result.Responses, response)
	}

	result.EndTime = time.Now()
	return result, nil
}

// HasErrors returns true if any response has an error.
func (r *Result) HasErrors() bool {
	for _, resp := range r.Responses {
		if resp.Error != nil {
			return true
		}
	}
	return false
}

// ErrorCount returns the number of failed requests.
func (r *Result) ErrorCount() int {
	count := 0
	for _, resp := range r.Responses {
		if resp.Error != nil {
			count++
		}
	}
	return count
}

// SuccessCount returns the number of successful requests.
func (r *Result) SuccessCount() int {
	count := 0
	for _, resp := range r.Responses {
		if resp.Error == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			count++
		}
	}
	return count
}

// AverageDuration returns the average duration of successful requests.
func (r *Result) AverageDuration() time.Duration {
	if len(r.Responses) == 0 {
		return 0
	}

	var total time.Duration
	successCount := 0
	for _, resp := range r.Responses {
		if resp.Error == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			total += resp.Duration
			successCount++
		}
	}

	if successCount == 0 {
		return 0
	}

	return total / time.Duration(successCount)
}

// Count2xx returns the number of 2xx responses.
func (r *Result) Count2xx() int {
	count := 0
	for _, resp := range r.Responses {
		if resp.Error == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			count++
		}
	}
	return count
}

// Count3xx returns the number of 3xx responses.
func (r *Result) Count3xx() int {
	count := 0
	for _, resp := range r.Responses {
		if resp.Error == nil && resp.StatusCode >= 300 && resp.StatusCode < 400 {
			count++
		}
	}
	return count
}

// Count4xx returns the number of 4xx responses.
func (r *Result) Count4xx() int {
	count := 0
	for _, resp := range r.Responses {
		if resp.Error == nil && resp.StatusCode >= 400 && resp.StatusCode < 500 {
			count++
		}
	}
	return count
}

// Count5xx returns the number of 5xx responses.
func (r *Result) Count5xx() int {
	count := 0
	for _, resp := range r.Responses {
		if resp.Error == nil && resp.StatusCode >= 500 && resp.StatusCode < 600 {
			count++
		}
	}
	return count
}
