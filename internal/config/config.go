// Package config provides configuration management for the conreq application.
package config

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Config holds all configuration parameters for concurrent requests.
type Config struct {
	URL             string
	Method          string
	Count           int
	Headers         map[string]string
	Body            string
	RequestID       string
	SameRequestID   bool
	RequestIDHeader string
	Delay           time.Duration
	Timeout         time.Duration
	OutputJSON      bool
	NoBody          bool
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		Method:          "GET",
		Count:           1,
		Headers:         make(map[string]string),
		Timeout:         30 * time.Second,
		Delay:           0,
		RequestIDHeader: "X-Request-ID",
		SameRequestID:   false,
		NoBody:          false,
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("URLが指定されていません")
	}

	if !isValidHTTPMethod(c.Method) {
		return fmt.Errorf("無効なHTTPメソッド: %s", c.Method)
	}

	if c.Count < 1 || c.Count > 5 {
		return fmt.Errorf("同時リクエスト数は1-5の範囲で指定してください: %d", c.Count)
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("タイムアウトは正の値を指定してください: %s", c.Timeout)
	}

	if c.Delay < 0 {
		return fmt.Errorf("遅延時間は0以上の値を指定してください: %s", c.Delay)
	}

	return nil
}

// ParseHeaders parses header strings and adds them to the config.
func (c *Config) ParseHeaders(headers []string) error {
	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("無効なヘッダー形式: %s", header)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		c.Headers[key] = value
	}
	return nil
}

func isValidHTTPMethod(method string) bool {
	validMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, valid := range validMethods {
		if strings.EqualFold(method, valid) {
			return true
		}
	}
	return false
}

// ParseDuration parses a duration string.
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	return time.ParseDuration(s)
}
