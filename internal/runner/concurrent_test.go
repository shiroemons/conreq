package runner

import (
	"context"
	"testing"
	"time"

	"github.com/shiroemons/conreq/internal/client"
)

func TestResultMethods(t *testing.T) {
	t.Run("HasErrors", func(t *testing.T) {
		result := &Result{
			Responses: []*client.Response{
				{Error: nil},
				{Error: nil},
			},
		}

		if result.HasErrors() {
			t.Error("HasErrors() = true, want false")
		}

		result.Responses = append(result.Responses, &client.Response{
			Error: context.DeadlineExceeded,
		})

		if !result.HasErrors() {
			t.Error("HasErrors() = false, want true")
		}
	})

	t.Run("ErrorCount", func(t *testing.T) {
		result := &Result{
			Responses: []*client.Response{
				{Error: nil},
				{Error: context.DeadlineExceeded},
				{Error: nil},
				{Error: context.Canceled},
			},
		}

		if got := result.ErrorCount(); got != 2 {
			t.Errorf("ErrorCount() = %d, want 2", got)
		}
	})

	t.Run("SuccessCount", func(t *testing.T) {
		result := &Result{
			Responses: []*client.Response{
				{Error: nil, StatusCode: 200},
				{Error: context.DeadlineExceeded},
				{Error: nil, StatusCode: 201},
				{Error: context.Canceled},
				{Error: nil, StatusCode: 404},
				{Error: nil, StatusCode: 500},
			},
		}

		if got := result.SuccessCount(); got != 2 {
			t.Errorf("SuccessCount() = %d, want 2", got)
		}
	})

	t.Run("AverageDuration", func(t *testing.T) {
		result := &Result{
			Responses: []*client.Response{
				{Error: nil, StatusCode: 200, Duration: 100 * time.Millisecond},
				{Error: nil, StatusCode: 201, Duration: 200 * time.Millisecond},
				{Error: context.DeadlineExceeded, Duration: 500 * time.Millisecond},
				{Error: nil, StatusCode: 404, Duration: 300 * time.Millisecond},
			},
		}

		want := 150 * time.Millisecond
		if got := result.AverageDuration(); got != want {
			t.Errorf("AverageDuration() = %v, want %v", got, want)
		}
	})

	t.Run("AverageDuration with no successful responses", func(t *testing.T) {
		result := &Result{
			Responses: []*client.Response{
				{Error: context.DeadlineExceeded},
				{Error: context.Canceled},
			},
		}

		if got := result.AverageDuration(); got != 0 {
			t.Errorf("AverageDuration() = %v, want 0", got)
		}
	})

	t.Run("AverageDuration with empty responses", func(t *testing.T) {
		result := &Result{
			Responses: []*client.Response{},
		}

		if got := result.AverageDuration(); got != 0 {
			t.Errorf("AverageDuration() = %v, want 0", got)
		}
	})

	t.Run("Status Code Counts", func(t *testing.T) {
		result := &Result{
			Responses: []*client.Response{
				{Error: nil, StatusCode: 200},
				{Error: nil, StatusCode: 201},
				{Error: nil, StatusCode: 301},
				{Error: nil, StatusCode: 302},
				{Error: nil, StatusCode: 400},
				{Error: nil, StatusCode: 404},
				{Error: nil, StatusCode: 500},
				{Error: nil, StatusCode: 503},
				{Error: context.DeadlineExceeded},
			},
		}

		if got := result.Count2xx(); got != 2 {
			t.Errorf("Count2xx() = %d, want 2", got)
		}
		if got := result.Count3xx(); got != 2 {
			t.Errorf("Count3xx() = %d, want 2", got)
		}
		if got := result.Count4xx(); got != 2 {
			t.Errorf("Count4xx() = %d, want 2", got)
		}
		if got := result.Count5xx(); got != 2 {
			t.Errorf("Count5xx() = %d, want 2", got)
		}
	})
}
