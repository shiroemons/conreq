package config

import (
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	
	if cfg.Method != "GET" {
		t.Errorf("Expected default method to be GET, got %s", cfg.Method)
	}
	
	if cfg.Count != 1 {
		t.Errorf("Expected default count to be 1, got %d", cfg.Count)
	}
	
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30s, got %v", cfg.Timeout)
	}
	
	if cfg.Headers == nil {
		t.Error("Expected headers map to be initialized")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				URL:     "https://example.com",
				Method:  "GET",
				Count:   3,
				Timeout: 30 * time.Second,
				Delay:   0,
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			config: &Config{
				Method:  "GET",
				Count:   1,
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			config: &Config{
				URL:     "https://example.com",
				Method:  "INVALID",
				Count:   1,
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "count too low",
			config: &Config{
				URL:     "https://example.com",
				Method:  "GET",
				Count:   0,
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "count too high",
			config: &Config{
				URL:     "https://example.com",
				Method:  "GET",
				Count:   6,
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			config: &Config{
				URL:     "https://example.com",
				Method:  "GET",
				Count:   1,
				Timeout: -1 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "negative delay",
			config: &Config{
				URL:     "https://example.com",
				Method:  "GET",
				Count:   1,
				Timeout: 30 * time.Second,
				Delay:   -1 * time.Second,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "valid headers",
			headers: []string{
				"Content-Type: application/json",
				"Authorization: Bearer token123",
			},
			want: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer token123",
			},
			wantErr: false,
		},
		{
			name: "header with extra spaces",
			headers: []string{
				"Content-Type:   application/json   ",
			},
			want: map[string]string{
				"Content-Type": "application/json",
			},
			wantErr: false,
		},
		{
			name:    "invalid header format",
			headers: []string{"InvalidHeader"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty headers",
			headers: []string{},
			want:    map[string]string{},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewConfig()
			err := cfg.ParseHeaders(tt.headers)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if !tt.wantErr {
				if len(cfg.Headers) != len(tt.want) {
					t.Errorf("Expected %d headers, got %d", len(tt.want), len(cfg.Headers))
				}
				
				for key, value := range tt.want {
					if cfg.Headers[key] != value {
						t.Errorf("Expected header %s = %s, got %s", key, value, cfg.Headers[key])
					}
				}
			}
		})
	}
}

func TestIsValidHTTPMethod(t *testing.T) {
	tests := []struct {
		method string
		want   bool
	}{
		{"GET", true},
		{"get", true},
		{"POST", true},
		{"PUT", true},
		{"DELETE", true},
		{"PATCH", true},
		{"HEAD", true},
		{"OPTIONS", true},
		{"INVALID", false},
		{"", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			if got := isValidHTTPMethod(tt.method); got != tt.want {
				t.Errorf("isValidHTTPMethod(%s) = %v, want %v", tt.method, got, tt.want)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "valid duration - seconds",
			input:   "10s",
			want:    10 * time.Second,
			wantErr: false,
		},
		{
			name:    "valid duration - milliseconds",
			input:   "500ms",
			want:    500 * time.Millisecond,
			wantErr: false,
		},
		{
			name:    "valid duration - minutes",
			input:   "2m",
			want:    2 * time.Minute,
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			want:    0,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}