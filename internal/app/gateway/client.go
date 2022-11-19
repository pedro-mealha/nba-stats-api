package gateway

import (
	"net/http"
	"time"
)

// NewClientWithTimeout create an http.Client with timeout.
func NewClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
