package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

const (
	defaultBaseURL = "https://api.prod.whoop.com/developer"
	userAgent      = "whoop-cli/1.0.0"
)

// HTTPDoer is the minimal interface for making HTTP requests.
// Allows tests to inject a mock http client.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is an authenticated Whoop API client with rate-limit retry awareness.
type Client struct {
	baseURL string
	http    HTTPDoer
}

// NewClient creates a Client whose HTTP requests are authenticated via the
// provided oauth2.TokenSource (which handles automatic token refresh).
func NewClient(ctx context.Context, ts oauth2.TokenSource) *Client {
	return &Client{
		baseURL: defaultBaseURL,
		http:    oauth2.NewClient(ctx, ts),
	}
}

// NewClientWithBaseURL creates a Client with a custom base URL and HTTP doer.
// Used in tests to point at a mock server.
func NewClientWithBaseURL(baseURL string, httpClient HTTPDoer) *Client {
	return &Client{
		baseURL: baseURL,
		http:    httpClient,
	}
}

// get performs a GET request, retrying once on HTTP 429.
func (c *Client) get(ctx context.Context, path string, query url.Values, dest any) error {
	rawURL := c.baseURL + path
	if len(query) > 0 {
		rawURL += "?" + query.Encode()
	}
	return c.doWithRetry(ctx, http.MethodGet, rawURL, nil, dest)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, body io.Reader) error {
	return c.doWithRetry(ctx, http.MethodPost, c.baseURL+path, body, nil)
}

// doWithRetry executes the request; on 429 it waits for the rate-limit reset
// and retries once.
func (c *Client) doWithRetry(ctx context.Context, method, rawURL string, body io.Reader, dest any) error {
	resp, err := c.execute(ctx, method, rawURL, body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		waitForRateLimit(resp)
		resp, err = c.execute(ctx, method, rawURL, body)
		if err != nil {
			return err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Message: string(b)}
	}
	if dest == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

func (c *Client) execute(ctx context.Context, method, rawURL string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	return resp, nil
}

// waitForRateLimit reads X-RateLimit-Reset from resp and sleeps until that time.
func waitForRateLimit(resp *http.Response) {
	resp.Body.Close()
	resetStr := resp.Header.Get("X-RateLimit-Reset")
	if resetStr != "" {
		if resetEpoch, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
			waitDur := time.Until(time.Unix(resetEpoch, 0))
			if waitDur > 0 && waitDur < 2*time.Minute {
				time.Sleep(waitDur)
				return
			}
		}
	}
	// Fallback: short fixed wait.
	time.Sleep(500 * time.Millisecond)
}

// buildQuery converts ListParams into url.Values.
func buildQuery(p ListParams) url.Values {
	q := url.Values{}
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.NextToken != "" {
		q.Set("nextToken", p.NextToken)
	}
	if p.Start != nil {
		q.Set("start", p.Start.UTC().Format(time.RFC3339))
	}
	if p.End != nil {
		q.Set("end", p.End.UTC().Format(time.RFC3339))
	}
	return q
}
