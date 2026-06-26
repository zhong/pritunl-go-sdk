package pritunl

import "fmt"

// ListOrganizations returns all organizations.
func (c *Client) ListOrganizations() ([]Organization, error) {
	var orgs []Organization
	if err := c.Get("/organization", &orgs); err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetOrganization returns a single organization by ID.
func (c *Client) GetOrganization(id string) (*Organization, error) {
	org := &Organization{}
	if err := c.Get("/organization/"+id, org); err != nil {
		return nil, err
	}
	return org, nil
}

// CreateOrganization creates a new organization.
func (c *Client) CreateOrganization(name string) (*Organization, error) {
	body := map[string]interface{}{"name": name}
	org := &Organization{}
	if err := c.Post("/organization", body, org); err != nil {
		return nil, err
	}
	return org, nil
}

// UpdateOrganization renames an organization.
func (c *Client) UpdateOrganization(id, name string) (*Organization, error) {
	body := map[string]interface{}{"name": name}
	org := &Organization{}
	if err := c.Put("/organization/"+id, body, org); err != nil {
		return nil, err
	}
	return org, nil
}

// DeleteOrganization deletes an organization.
func (c *Client) DeleteOrganization(id string) error {
	return c.Delete("/organization/"+id, nil)
}

// EnableOrganizationAPI enables API token authentication for an organization.
func (c *Client) EnableOrganizationAPI(id string) (*Organization, error) {
	org, err := c.GetOrganization(id)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"name":     org.Name,
		"auth_api": true,
	}
	updated := &Organization{}
	if err := c.Put("/organization/"+id, body, updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// FindOrganizationByName finds an organization by name.
func (c *Client) FindOrganizationByName(name string) (*Organization, error) {
	orgs, err := c.ListOrganizations()
	if err != nil {
		return nil, err
	}
	for i := range orgs {
		if orgs[i].Name == name {
			return &orgs[i], nil
		}
	}
	return nil, fmt.Errorf("organization %q not found", name)
}
