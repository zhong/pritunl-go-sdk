package pritunl

import "fmt"

// CreateServerRequest is the payload for creating a server.
type CreateServerRequest struct {
	Name          string        `json:"name"`
	Network       string        `json:"network,omitempty"`
	NetworkWG     string        `json:"network_wg,omitempty"`
	Port          int           `json:"port,omitempty"`
	Protocol      string        `json:"protocol,omitempty"`
	WG            bool          `json:"wg,omitempty"`
	OTPAuth       bool          `json:"otp_auth,omitempty"`
	Cipher        string        `json:"cipher,omitempty"`
	Hash          string        `json:"hash,omitempty"`
	LocalNetworks []string      `json:"local_networks,omitempty"`
	DNSServers    []string      `json:"dns_servers,omitempty"`
	DNSSuffix     string        `json:"dns_suffix,omitempty"`
	IPv6          bool          `json:"ipv6,omitempty"`
}

// ListServers returns all servers.
func (c *Client) ListServers() ([]Server, error) {
	var servers []Server
	if err := c.Get("/server", &servers); err != nil {
		return nil, err
	}
	return servers, nil
}

// GetServer returns a single server.
func (c *Client) GetServer(id string) (*Server, error) {
	server := &Server{}
	if err := c.Get("/server/"+id, server); err != nil {
		return nil, err
	}
	return server, nil
}

// CreateServer creates a new server.
func (c *Client) CreateServer(req CreateServerRequest) (*Server, error) {
	server := &Server{}
	if err := c.Post("/server", req, server); err != nil {
		return nil, err
	}
	return server, nil
}

// UpdateServer updates a server.
func (c *Client) UpdateServer(id string, req CreateServerRequest) (*Server, error) {
	server := &Server{}
	if err := c.Put("/server/"+id, req, server); err != nil {
		return nil, err
	}
	return server, nil
}

// DeleteServer deletes a server.
func (c *Client) DeleteServer(id string) error {
	return c.Delete("/server/"+id, nil)
}

// ServerOperation starts, stops or restarts a server.
func (c *Client) ServerOperation(id, operation string) error {
	if operation != "start" && operation != "stop" && operation != "restart" {
		return fmt.Errorf("invalid operation %q", operation)
	}
	return c.Put("/server/"+id+"/operation/"+operation, nil, nil)
}

// AttachOrganization attaches an organization to a server.
func (c *Client) AttachOrganization(serverID, orgID string) error {
	return c.Put("/server/"+serverID+"/organization/"+orgID, nil, nil)
}

// DetachOrganization detaches an organization from a server.
func (c *Client) DetachOrganization(serverID, orgID string) error {
	return c.Delete("/server/"+serverID+"/organization/"+orgID, nil)
}

// AddRoute adds a route to a server.
func (c *Client) AddRoute(serverID string, route ServerRoute) error {
	return c.Post("/server/"+serverID+"/route", route, nil)
}

// DeleteRoute deletes a route from a server.
func (c *Client) DeleteRoute(serverID, network string) error {
	return c.Delete("/server/"+serverID+"/route/"+network, nil)
}
