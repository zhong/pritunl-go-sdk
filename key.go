package pritunl

import (
	"fmt"
	"os"
	"strings"
)

// isNotFound checks if an error represents an HTTP 404 response.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "unexpected status 404")
}

// GetKeyLink creates a temporary key link for a user.
// Tries the modern /data/ path first, then falls back to the legacy /key/ path
// for older Pritunl versions (< 1.32.4350 approx).
func (c *Client) GetKeyLink(orgID, userID string) (*KeyLink, error) {
	link := &KeyLink{}
	if err := c.Get("/data/"+orgID+"/"+userID, link); err != nil {
		if !isNotFound(err) {
			return nil, err
		}
		// Fallback for older Pritunl versions
		if err := c.Get("/key/"+orgID+"/"+userID, link); err != nil {
			return nil, err
		}
	}
	return link, nil
}

// DownloadKeyArchive downloads the user key archive in the given format.
// format can be "tar", "zip" or "onc".
// Tries the modern /data/ path first, then falls back to legacy /key/ or /key_onc/ paths.
func (c *Client) DownloadKeyArchive(orgID, userID, format string) ([]byte, error) {
	if format != "tar" && format != "zip" && format != "onc" {
		return nil, fmt.Errorf("unsupported archive format %q", format)
	}

	modernPath := fmt.Sprintf("/data/%s/%s.%s", orgID, userID, format)
	body, _, err := c.Request("GET", modernPath, nil)
	if err == nil {
		return body, nil
	}
	if !isNotFound(err) {
		return nil, err
	}

	// Fallback paths for older Pritunl versions
	var legacyPath string
	if format == "onc" {
		legacyPath = fmt.Sprintf("/key_onc/%s/%s.onc", orgID, userID)
	} else {
		legacyPath = fmt.Sprintf("/key/%s/%s.%s", orgID, userID, format)
	}
	body, _, err = c.Request("GET", legacyPath, nil)
	return body, err
}

// DownloadKeyConfig downloads the OpenVPN configuration for a specific server.
func (c *Client) DownloadKeyConfig(orgID, userID, serverID string) ([]byte, error) {
	modernPath := fmt.Sprintf("/data/%s/%s/%s.key", orgID, userID, serverID)
	body, _, err := c.Request("GET", modernPath, nil)
	if err == nil {
		return body, nil
	}
	if !isNotFound(err) {
		return nil, err
	}

	// Fallback for older Pritunl versions
	legacyPath := fmt.Sprintf("/key/%s/%s/%s.key", orgID, userID, serverID)
	body, _, err = c.Request("GET", legacyPath, nil)
	return body, err
}

// DownloadLinkedKeyArchive downloads a key by key_id.
func (c *Client) DownloadLinkedKeyArchive(keyID, format string) ([]byte, error) {
	if format != "tar" && format != "zip" {
		return nil, fmt.Errorf("unsupported archive format %q", format)
	}
	path := fmt.Sprintf("/key/%s.%s", keyID, format)
	body, _, err := c.Request("GET", path, nil)
	return body, err
}

// SaveUserKeyArchive is a helper that writes the tar/zip archive to a file.
func SaveUserKeyArchive(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0600)
}
