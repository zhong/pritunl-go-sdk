package pritunl

// GetStatus returns general information about the Pritunl server.
func (c *Client) GetStatus() (*Status, error) {
	status := &Status{}
	if err := c.Get("/status", status); err != nil {
		return nil, err
	}
	return status, nil
}
