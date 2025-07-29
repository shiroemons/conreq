package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shiroemons/conreq/internal/config"
)

type Response struct {
	RequestID      string
	StatusCode     int
	Headers        http.Header
	Body           string
	Duration       time.Duration
	Timestamp      time.Time
	RequestIndex   int
	Error          error
}

type Client struct {
	httpClient *http.Client
	config     *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		config: cfg,
	}
}

func (c *Client) Do(ctx context.Context, requestIndex int) *Response {
	start := time.Now()
	response := &Response{
		RequestIndex: requestIndex,
		Timestamp:    start,
	}

	req, err := c.createRequest(ctx)
	if err != nil {
		response.Error = err
		response.Duration = time.Since(start)
		return response
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		response.Error = err
		response.Duration = time.Since(start)
		return response
	}
	defer resp.Body.Close()

	response.StatusCode = resp.StatusCode
	response.Headers = resp.Header
	response.RequestID = resp.Header.Get("X-Request-ID")
	response.Duration = time.Since(start)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error = fmt.Errorf("レスポンスボディの読み取りエラー: %w", err)
		return response
	}
	response.Body = string(body)

	return response
}

func (c *Client) createRequest(ctx context.Context) (*http.Request, error) {
	var body io.Reader
	if c.config.Body != "" {
		body = strings.NewReader(c.config.Body)
	}

	req, err := http.NewRequestWithContext(ctx, c.config.Method, c.config.URL, body)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	if c.config.RequestID != "" {
		req.Header.Set("X-Request-ID", c.config.RequestID)
	}

	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) DoWithDelay(ctx context.Context, requestIndex int, delay time.Duration) *Response {
	if delay > 0 {
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return &Response{
				RequestIndex: requestIndex,
				Error:        ctx.Err(),
				Timestamp:    time.Now(),
			}
		}
	}

	return c.Do(ctx, requestIndex)
}