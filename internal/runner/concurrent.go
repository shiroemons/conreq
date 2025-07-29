package runner

import (
	"context"
	"sync"
	"time"

	"github.com/shiroemons/conreq/internal/client"
	"github.com/shiroemons/conreq/internal/config"
	"github.com/shiroemons/conreq/pkg/requestid"
)

type Result struct {
	Responses []*client.Response
	StartTime time.Time
	EndTime   time.Time
}

type Runner struct {
	config *config.Config
	client *client.Client
}

func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		config: cfg,
		client: client.NewClient(cfg),
	}
}

func (r *Runner) Run(ctx context.Context) (*Result, error) {
	result := &Result{
		StartTime: time.Now(),
		Responses: make([]*client.Response, 0, r.config.Count),
	}

	responseChan := make(chan *client.Response, r.config.Count)
	var wg sync.WaitGroup

	for i := 0; i < r.config.Count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			cfg := *r.config
			if cfg.RequestID == "" {
				cfg.RequestID = requestid.Generate()
			}

			client := client.NewClient(&cfg)
			
			delay := time.Duration(index) * r.config.Delay
			response := client.DoWithDelay(ctx, index, delay)
			
			responseChan <- response
		}(i)
	}

	go func() {
		wg.Wait()
		close(responseChan)
	}()

	for response := range responseChan {
		result.Responses = append(result.Responses, response)
	}

	result.EndTime = time.Now()
	return result, nil
}

func (r *Result) HasErrors() bool {
	for _, resp := range r.Responses {
		if resp.Error != nil {
			return true
		}
	}
	return false
}

func (r *Result) ErrorCount() int {
	count := 0
	for _, resp := range r.Responses {
		if resp.Error != nil {
			count++
		}
	}
	return count
}

func (r *Result) SuccessCount() int {
	count := 0
	for _, resp := range r.Responses {
		if resp.Error == nil {
			count++
		}
	}
	return count
}

func (r *Result) AverageDuration() time.Duration {
	if len(r.Responses) == 0 {
		return 0
	}

	var total time.Duration
	successCount := 0
	for _, resp := range r.Responses {
		if resp.Error == nil {
			total += resp.Duration
			successCount++
		}
	}

	if successCount == 0 {
		return 0
	}

	return total / time.Duration(successCount)
}