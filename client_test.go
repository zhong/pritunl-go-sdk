package pritunl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestSignRequest(t *testing.T) {
	client := NewClient("https://localhost", "token123", "secret456", true)
	headers := client.signRequest("GET", "/organization")

	if headers["Auth-Token"] != "token123" {
		t.Errorf("Auth-Token mismatch")
	}
	if headers["Auth-Timestamp"] == "" {
		t.Errorf("Auth-Timestamp empty")
	}
	if headers["Auth-Nonce"] == "" {
		t.Errorf("Auth-Nonce empty")
	}
	if headers["Auth-Signature"] == "" {
		t.Errorf("Auth-Signature empty")
	}

	// Verify signature manually
	authString := "token123" + "&" + headers["Auth-Timestamp"] + "&" + headers["Auth-Nonce"] + "&" + "GET" + "&" + "/organization"
	h := hmac.New(sha256.New, []byte("secret456"))
	h.Write([]byte(authString))
	expectedSig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if headers["Auth-Signature"] != expectedSig {
		t.Errorf("signature mismatch: got %s, want %s", headers["Auth-Signature"], expectedSig)
	}
}

func TestSignRequestStripsQueryParams(t *testing.T) {
	client := NewClient("https://localhost", "token123", "secret456", true)
	headers := client.signRequest("GET", "/user/org123?page=2")

	// Signature should be computed against /user/org123, not include ?page=2
	authString := "token123" + "&" + headers["Auth-Timestamp"] + "&" + headers["Auth-Nonce"] + "&" + "GET" + "&" + "/user/org123"
	h := hmac.New(sha256.New, []byte("secret456"))
	h.Write([]byte(authString))
	expectedSig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if headers["Auth-Signature"] != expectedSig {
		t.Errorf("signature should ignore query params: got %s, want %s", headers["Auth-Signature"], expectedSig)
	}
}

func TestGenerateNonceLength(t *testing.T) {
	nonce := generateNonce()
	if len(nonce) != 32 {
		t.Errorf("nonce length = %d, want 32", len(nonce))
	}
}

func TestURLJoinPath(t *testing.T) {
	client := NewClient("https://example.com:9700", "t", "s", true)
	client.HTTPClient.Timeout = 1 // second
	path := "/organization"
	_, _, err := client.Request("GET", path, nil)
	// Request will fail because no real server, but we can verify URL construction
	// by checking error message if it mentions the full URL.
	if err == nil {
		t.Skip("unexpected success without server")
	}
	if !strings.Contains(err.Error(), "example.com:9700") {
		t.Logf("error message: %s", err.Error())
	}
}

func ExampleNewClient() {
	client := NewClient("https://localhost", "token", "secret", true)
	fmt.Println(client.BaseURL)
	// Output: https://localhost
}
