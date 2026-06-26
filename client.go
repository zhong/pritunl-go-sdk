// Package pritunl provides a Go client for the Pritunl VPN server REST API.
package pritunl

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// randReader is used to generate nonces. It can be overridden in tests.
var randReader io.Reader = rand.Reader

// Client is the Pritunl API client.
type Client struct {
	BaseURL    string
	APIToken   string
	APISecret  string
	HTTPClient *http.Client
}

// NewClient creates a new Pritunl API client.
// If insecure is true, TLS certificate verification is skipped.
func NewClient(baseURL, apiToken, apiSecret string, insecure bool) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	return &Client{
		BaseURL:   baseURL,
		APIToken:  apiToken,
		APISecret: apiSecret,
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

// signRequest computes the Pritunl API signature and returns the auth headers.
func (c *Client) signRequest(method, path string) map[string]string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := generateNonce()

	// Signature string: TOKEN & TIMESTAMP & NONCE & METHOD & PATH
	// Path must not include query parameters.
	signPath := path
	if idx := stringsIndex(signPath, '?'); idx != -1 {
		signPath = signPath[:idx]
	}
	authString := c.APIToken + "&" + timestamp + "&" + nonce + "&" + method + "&" + signPath

	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(authString))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return map[string]string{
		"Auth-Token":     c.APIToken,
		"Auth-Timestamp": timestamp,
		"Auth-Nonce":     nonce,
		"Auth-Signature": signature,
		"Accept":         "application/json",
	}
}

// generateNonce returns a 32-char hex random nonce.
func generateNonce() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(randReader, b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}

// stringsIndex is a tiny helper to avoid importing strings just for Index.
func stringsIndex(s string, substr byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == substr {
			return i
		}
	}
	return -1
}

// Request performs an authenticated HTTP request and returns the response body.
func (c *Client) Request(method, path string, body interface{}) ([]byte, *http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	u, err := url.JoinPath(c.BaseURL, path)
	if err != nil {
		return nil, nil, fmt.Errorf("build url: %w", err)
	}

	req, err := http.NewRequest(method, u, bodyReader)
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}

	for k, v := range c.signRequest(method, path) {
		req.Header.Set(k, v)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, resp, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, resp, nil
}

// Get performs an authenticated GET request and unmarshals the JSON response into v.
func (c *Client) Get(path string, v interface{}) error {
	body, _, err := c.Request("GET", path, nil)
	if err != nil {
		return err
	}
	if v != nil {
		return json.Unmarshal(body, v)
	}
	return nil
}

// Post performs an authenticated POST request.
func (c *Client) Post(path string, body, v interface{}) error {
	return c.doJSON("POST", path, body, v)
}

// Put performs an authenticated PUT request.
func (c *Client) Put(path string, body, v interface{}) error {
	return c.doJSON("PUT", path, body, v)
}

// Delete performs an authenticated DELETE request.
func (c *Client) Delete(path string, v interface{}) error {
	return c.doJSON("DELETE", path, nil, v)
}

func (c *Client) doJSON(method, path string, body, v interface{}) error {
	respBody, _, err := c.Request(method, path, body)
	if err != nil {
		return err
	}
	if v != nil {
		return json.Unmarshal(respBody, v)
	}
	return nil
}
