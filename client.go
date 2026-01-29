package reve

import (
	"net/http"
	"time"

	"github.com/shamspias/reve-go/image"
	"github.com/shamspias/reve-go/internal/transport"
)

// Default configuration values.
const (
	DefaultBaseURL      = "https://api.reve.com"
	DefaultTimeout      = 120 * time.Second
	DefaultMaxRetries   = 3
	DefaultRetryMinWait = 1 * time.Second
	DefaultRetryMaxWait = 30 * time.Second
	DefaultUserAgent    = "reve-go/" + Version
)

// Client is the Reve API client.
type Client struct {
	// Images provides image generation operations.
	Images *image.Service

	config *Config
}

// Config holds client configuration.
type Config struct {
	APIKey       string
	BaseURL      string
	Timeout      time.Duration
	MaxRetries   int
	RetryMinWait time.Duration
	RetryMaxWait time.Duration
	UserAgent    string
	Debug        bool
	Logger       func(format string, args ...any)
	Transport    http.RoundTripper
}

// NewClient creates a new Reve API client.
//
// Example:
//
//	// Basic
//	client := reve.NewClient(apiKey)
//
//	// With options
//	client := reve.NewClient(apiKey,
//		reve.WithTimeout(60*time.Second),
//		reve.WithDebug(true),
//	)
//
//	// With HTTP proxy
//	client := reve.NewClient(apiKey,
//		reve.WithHTTPProxy("http://proxy:8080"),
//	)
//
//	// With SOCKS5 proxy
//	client := reve.NewClient(apiKey,
//		reve.WithSOCKS5Proxy("127.0.0.1:1080", "user", "pass"),
//	)
func NewClient(apiKey string, opts ...Option) *Client {
	config := &Config{
		APIKey:       apiKey,
		BaseURL:      DefaultBaseURL,
		Timeout:      DefaultTimeout,
		MaxRetries:   DefaultMaxRetries,
		RetryMinWait: DefaultRetryMinWait,
		RetryMaxWait: DefaultRetryMaxWait,
		UserAgent:    DefaultUserAgent,
	}

	for _, opt := range opts {
		opt(config)
	}

	t := transport.New(&transport.Config{
		BaseURL:      config.BaseURL,
		APIKey:       config.APIKey,
		UserAgent:    config.UserAgent,
		Timeout:      config.Timeout,
		MaxRetries:   config.MaxRetries,
		RetryMinWait: config.RetryMinWait,
		RetryMaxWait: config.RetryMaxWait,
		Debug:        config.Debug,
		Logger:       config.Logger,
		Transport:    config.Transport,
	})

	return &Client{
		Images: image.NewService(t),
		config: config,
	}
}

// Config returns a copy of the client configuration.
func (c *Client) Config() Config {
	return *c.config
}
