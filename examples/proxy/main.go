// Example: Proxy Configuration
//
// This example demonstrates all proxy configuration options.
//
// Run with:
//
//	REVE_API_KEY=your-key go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	reve "github.com/shamspias/reve-go"
)

func main() {
	apiKey := os.Getenv("REVE_API_KEY")
	if apiKey == "" {
		log.Fatal("REVE_API_KEY required")
	}

	ctx := context.Background()

	// =========================================================================
	// Example 1: HTTP Proxy
	// =========================================================================
	fmt.Println("=== HTTP Proxy Configuration ===")
	fmt.Println(`
// Simple HTTP proxy
client := reve.NewClient(apiKey,
    reve.WithHTTPProxy("http://proxy.example.com:8080"),
)

// HTTP proxy with authentication
client := reve.NewClient(apiKey,
    reve.WithHTTPProxy("http://user:password@proxy.example.com:8080"),
)
`)

	// =========================================================================
	// Example 2: HTTPS Proxy
	// =========================================================================
	fmt.Println("=== HTTPS Proxy Configuration ===")
	fmt.Println(`
// HTTPS proxy (same as HTTP proxy for most cases)
client := reve.NewClient(apiKey,
    reve.WithHTTPSProxy("https://proxy.example.com:8443"),
)
`)

	// =========================================================================
	// Example 3: SOCKS5 Proxy
	// =========================================================================
	fmt.Println("=== SOCKS5 Proxy Configuration ===")
	fmt.Println(`
// SOCKS5 proxy without authentication
client := reve.NewClient(apiKey,
    reve.WithSOCKS5Proxy("127.0.0.1:1080", "", ""),
)

// SOCKS5 proxy with authentication
client := reve.NewClient(apiKey,
    reve.WithSOCKS5Proxy("proxy.example.com:1080", "username", "password"),
)
`)

	// =========================================================================
	// Example 4: Environment Proxy
	// =========================================================================
	fmt.Println("=== Environment Proxy Configuration ===")
	fmt.Println(`
// Uses HTTP_PROXY, HTTPS_PROXY, NO_PROXY environment variables
// Set in shell:
//   export HTTP_PROXY=http://proxy:8080
//   export HTTPS_PROXY=https://proxy:8443
//   export NO_PROXY=localhost,127.0.0.1

client := reve.NewClient(apiKey,
    reve.WithProxyFromEnvironment(),
)
`)

	// =========================================================================
	// Example 5: Custom Transport
	// =========================================================================
	fmt.Println("=== Custom Transport Configuration ===")
	fmt.Println(`
// For advanced proxy configurations, use a custom transport
transport := &http.Transport{
    Proxy: http.ProxyURL(proxyURL),
    MaxIdleConns: 10,
    IdleConnTimeout: 30 * time.Second,
}

client := reve.NewClient(apiKey,
    reve.WithTransport(transport),
)
`)

	// =========================================================================
	// Example 6: Combining with other options
	// =========================================================================
	fmt.Println("=== Combined Configuration ===")
	fmt.Println(`
// Combine proxy with other options
client := reve.NewClient(apiKey,
    reve.WithHTTPProxy("http://proxy:8080"),
    reve.WithTimeout(60*time.Second),
    reve.WithRetry(5, time.Second, 30*time.Second),
    reve.WithDebug(true),
    reve.WithLogger(func(format string, args ...any) {
        log.Printf("[REVE] "+format, args...)
    }),
)
`)

	// =========================================================================
	// Live test (if proxy env is set)
	// =========================================================================
	httpProxy := os.Getenv("HTTP_PROXY")
	if httpProxy != "" {
		fmt.Printf("\n=== Live Test with HTTP_PROXY=%s ===\n", httpProxy)

		client := reve.NewClient(apiKey,
			reve.WithProxyFromEnvironment(),
			reve.WithTimeout(30*time.Second),
			reve.WithDebug(true),
		)

		result, err := client.Images.Create(ctx, &reve.CreateParams{
			Prompt: "A simple test image",
		})
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Success! Credits used: %d\n", result.CreditsUsed)
		}
	} else {
		fmt.Println("\n=== Live Test (No Proxy) ===")

		// Test without proxy
		client := reve.NewClient(apiKey,
			reve.WithTimeout(30*time.Second),
		)

		result, err := client.Images.Create(ctx, &reve.CreateParams{
			Prompt: "A simple test image",
		})
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Success! Credits used: %d\n", result.CreditsUsed)
			result.SaveTo("proxy_test.png")
		}
	}

	// =========================================================================
	// Example 7: Practical usage scenarios
	// =========================================================================
	fmt.Println("\n=== Practical Scenarios ===")
	fmt.Println(`
// Scenario 1: Corporate network with proxy
client := reve.NewClient(apiKey,
    reve.WithHTTPProxy("http://corporate-proxy.internal:3128"),
)

// Scenario 2: VPN/Tor through SOCKS5
client := reve.NewClient(apiKey,
    reve.WithSOCKS5Proxy("127.0.0.1:9050", "", ""), // Tor default port
)

// Scenario 3: Development with local proxy (e.g., mitmproxy, Charles)
client := reve.NewClient(apiKey,
    reve.WithHTTPProxy("http://localhost:8888"),
)

// Scenario 4: Cloud environment with auto-detection
client := reve.NewClient(apiKey,
    reve.WithProxyFromEnvironment(), // Respects cloud proxy settings
)
`)

	// Suppress unused warnings
	_ = ctx
	_ = http.Transport{}
	_ = time.Second
}
