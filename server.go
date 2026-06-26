package pritunl

import (
	"fmt"
	"net/url"
)

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
// Pritunl social edition API behavior for route deletion.
func (c *Client) DeleteRoute(serverID, network string) error {
	// Try standard DELETE first
	encodedNetwork := url.QueryEscape(network)
	deleteErr := c.Delete("/server/"+serverID+"/route/"+encodedNetwork, nil)

	// If successful, return
	if deleteErr == nil {
		return nil
	}

	// If DELETE fails with 404, try alternative: maybe Pritunl expects PUT with empty body
	fmt.Printf("[SDK] DELETE failed with: %v\n", deleteErr)
	fmt.Printf("[SDK] Trying PUT method as alternative...\n")

	putErr := c.Put("/server/"+serverID+"/route/"+encodedNetwork, nil, nil)
	if putErr == nil {
		return nil
	}
	fmt.Printf("[SDK] PUT also failed with: %v\n", putErr)

	// If both DELETE and PUT fail, return the original DELETE error
	return deleteErr
}

// ReplaceRoutes replaces all routes on a server with a new set.
// This deletes all existing routes and adds the new ones.
func (c *Client) ReplaceRoutes(serverID string, newRoutes []ServerRoute) error {
	// First, try to update the server with the new routes list
	// This is the most direct approach - send PUT with the new routes
	type ServerRoutesUpdate struct {
		Routes []ServerRoute `json:"routes"`
	}

	updateReq := ServerRoutesUpdate{Routes: newRoutes}
	err := c.Put("/server/"+serverID, updateReq, nil)

	if err == nil {
		return nil
	}

	// If direct update fails, try alternative: clear first, then add
	// Get current routes
	currentRoutes, err := c.GetServerRoutes(serverID)
	if err != nil {
		return fmt.Errorf("get current routes: %w", err)
	}

	// Delete old routes one by one
	for _, route := range currentRoutes {
		encodedNetwork := url.QueryEscape(route.Network)
		_ = c.Delete("/server/"+serverID+"/route/"+encodedNetwork, nil)
	}

	// Add new routes
	for _, route := range newRoutes {
		if err := c.AddRoute(serverID, route); err != nil {
			return fmt.Errorf("add route %s: %w", route.Network, err)
		}
	}

	return nil
}

// ListRoutes returns all routes for a server.
// This gets the routes from the server detail endpoint.
func (c *Client) ListRoutes(serverID string) ([]ServerRoute, error) {
	server, err := c.GetServer(serverID)
	if err != nil {
		return nil, err
	}
	return server.Routes, nil
}

// GetServerRoutes fetches routes directly from the /server/{id}/route endpoint.
// This is the preferred method if GetServer doesn't include routes.
func (c *Client) GetServerRoutes(serverID string) ([]ServerRoute, error) {
	var routes []ServerRoute
	if err := c.Get("/server/"+serverID+"/route", &routes); err != nil {
		return nil, err
	}
	return routes, nil
}
