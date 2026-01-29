package reve

import (
	"net/http"
	"time"

	"github.com/shamspias/reve-go/internal/transport"
)

// Option is a functional option for Client configuration.
type Option func(*Config)

// WithBaseURL sets a custom API base URL.
//
// Example:
//
//	client := reve.NewClient(apiKey, reve.WithBaseURL("https://custom.api.com"))
func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithTimeout sets the request timeout.
//
// Example:
//
//	client := reve.NewClient(apiKey, reve.WithTimeout(60*time.Second))
func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Timeout = d
	}
}

// WithRetry configures retry behavior.
//
// Example:
//
//	client := reve.NewClient(apiKey, reve.WithRetry(5, time.Second, 30*time.Second))
func WithRetry(maxRetries int, minWait, maxWait time.Duration) Option {
	return func(c *Config) {
		c.MaxRetries = maxRetries
		c.RetryMinWait = minWait
		c.RetryMaxWait = maxWait
	}
}

// WithNoRetry disables automatic retries.
//
// Example:
//
//	client := reve.NewClient(apiKey, reve.WithNoRetry())
func WithNoRetry() Option {
	return func(c *Config) {
		c.MaxRetries = 0
	}
}

// WithUserAgent sets a custom User-Agent header.
//
// Example:
//
//	client := reve.NewClient(apiKey, reve.WithUserAgent("MyApp/1.0"))
func WithUserAgent(ua string) Option {
	return func(c *Config) {
		c.UserAgent = ua
	}
}

// WithDebug enables debug logging.
//
// Example:
//
//	client := reve.NewClient(apiKey, reve.WithDebug(true))
func WithDebug(enabled bool) Option {
	return func(c *Config) {
		c.Debug = enabled
	}
}

// WithLogger sets a custom logger.
//
// Example:
//
//	client := reve.NewClient(apiKey, reve.WithLogger(func(format string, args ...any) {
//		log.Printf("[REVE] "+format, args...)
//	}))
func WithLogger(logger func(format string, args ...any)) Option {
	return func(c *Config) {
		c.Logger = logger
		c.Debug = true
	}
}

// WithTransport sets a custom HTTP transport.
//
// Example:
//
//	transport := &http.Transport{MaxIdleConns: 10}
//	client := reve.NewClient(apiKey, reve.WithTransport(transport))
func WithTransport(t http.RoundTripper) Option {
	return func(c *Config) {
		c.Transport = t
	}
}

// WithHTTPProxy configures an HTTP/HTTPS proxy.
//
// Example:
//
//	// Simple proxy
//	client := reve.NewClient(apiKey, reve.WithHTTPProxy("http://proxy:8080"))
//
//	// With authentication
//	client := reve.NewClient(apiKey, reve.WithHTTPProxy("http://user:pass@proxy:8080"))
func WithHTTPProxy(proxyURL string) Option {
	return func(c *Config) {
		t, err := transport.CreateHTTPProxyTransport(proxyURL)
		if err == nil {
			c.Transport = t
		}
	}
}

// WithHTTPSProxy is an alias for WithHTTPProxy.
//
// Example:
//
//	client := reve.NewClient(apiKey, reve.WithHTTPSProxy("https://proxy:8443"))
func WithHTTPSProxy(proxyURL string) Option {
	return WithHTTPProxy(proxyURL)
}

// WithSOCKS5Proxy configures a SOCKS5 proxy.
//
// Example:
//
//	// Without authentication
//	client := reve.NewClient(apiKey, reve.WithSOCKS5Proxy("127.0.0.1:1080", "", ""))
//
//	// With authentication
//	client := reve.NewClient(apiKey, reve.WithSOCKS5Proxy("proxy:1080", "user", "pass"))
func WithSOCKS5Proxy(addr, username, password string) Option {
	return func(c *Config) {
		t, err := transport.CreateSOCKS5ProxyTransport(addr, username, password)
		if err == nil {
			c.Transport = t
		}
	}
}

// WithProxyFromEnvironment uses proxy from environment variables.
// Reads HTTP_PROXY, HTTPS_PROXY, NO_PROXY.
//
// Example:
//
//	// Set env: HTTP_PROXY=http://proxy:8080
//	client := reve.NewClient(apiKey, reve.WithProxyFromEnvironment())
func WithProxyFromEnvironment() Option {
	return func(c *Config) {
		c.Transport = transport.CreateEnvProxyTransport()
	}
}
