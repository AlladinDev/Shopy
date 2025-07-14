// Package test is for testing various utility functions
package test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

func isRetriableStatus(code int) bool {
	return code == http.StatusInternalServerError ||
		code == http.StatusBadGateway ||
		code == http.StatusGatewayTimeout ||
		code == http.StatusServiceUnavailable
}

func GoHit(method string, url string, bodyData []byte, timeout time.Duration, contentType string, maxRetries int, delay time.Duration) (*http.Response, error) {
	retryCount := 0
	client := http.Client{}

	for retryCount <= maxRetries {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyData))
		if err != nil {
			cancel()
			return nil, err
		}

		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}

		resp, err := client.Do(req)
		cancel()

		if err == nil {
			// Success
			return resp, nil
		}

		//only retry when statuscode is 500 or 502 means server error not on 400 401 like those errors
		if !isRetriableStatus(resp.StatusCode) {
			return nil, err
		}

		// Retry on error
		retryCount++
		if retryCount <= maxRetries {
			time.Sleep(delay)
			delay *= 2
		}
	}

	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

func main() {

}
