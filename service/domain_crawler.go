package service

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
)

// Custom client with redirect count check
var client = &http.Client{
	Timeout: 15 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 { // Limit redirects to 3
			return fmt.Errorf("stopped after 3 redirects")
		}
		return nil // Allow redirect
	},
}

func CrawlDomain(domain, pageType string) ([]byte, string, error) {
	url := fmt.Sprintf("http://%s/%s", domain, pageType)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, url, fmt.Errorf("failed to create request: %w", err)
	}

	// Adding custom headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		// main.TotalErrors++
		return nil, url, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// var resp *http.Response
	// retryCount := 2
	// for i := 0; i < retryCount; i++ {
	// 	resp, err = client.Do(req)
	// 	if err != nil {
	// 		if strings.Contains(err.Error(), "i/o timeout") && i < retryCount-1 {
	// 			time.Sleep(2 * time.Second) // wait before retrying
	// 			continue
	// 		}
	// 		return nil, url, fmt.Errorf("failed to fetch URL: %w", err)
	// 	}
	// 	defer resp.Body.Close()
	// }

	if resp.StatusCode != http.StatusOK {
		// totalErrors++
		return nil, url, fmt.Errorf("non-200 status code received: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, url, fmt.Errorf("failed to read response body: %w", err)
	}

	if !isValidAdsTxt(body) {
		return nil, url, fmt.Errorf("response does not look like an ads.txt file")
	}
	return body, url, nil
}

// isValidAdsTxt checks if the content looks like an ads.txt file.
func isValidAdsTxt(body []byte) bool {
	content := string(body)
	return !strings.Contains(content, "<html") && !strings.Contains(content, "<body")
}
