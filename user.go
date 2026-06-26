package pritunl

import "fmt"

// CreateUserRequest is the payload for creating a user.
type CreateUserRequest struct {
	Name            string   `json:"name"`
	Email           string   `json:"email,omitempty"`
	Groups          []string `json:"groups,omitempty"`
	Disabled        bool     `json:"disabled,omitempty"`
	BypassSecondary bool     `json:"bypass_secondary,omitempty"`
	ClientToClient  bool     `json:"client_to_client,omitempty"`
}

// ListUsers returns all users in an organization.
func (c *Client) ListUsers(orgID string) ([]User, error) {
	var users []User
	if err := c.Get("/user/"+orgID, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser returns a single user.
func (c *Client) GetUser(orgID, userID string) (*User, error) {
	user := &User{}
	if err := c.Get("/user/"+orgID+"/"+userID, user); err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user in an organization.
// Pritunl returns an array of created users even for a single creation.
func (c *Client) CreateUser(orgID string, req CreateUserRequest) (*User, error) {
	var users []User
	if err := c.Post("/user/"+orgID, req, &users); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user creation returned empty list")
	}
	return &users[0], nil
}

// UpdateUser updates an existing user.
func (c *Client) UpdateUser(orgID, userID string, req CreateUserRequest) (*User, error) {
	user := &User{}
	if err := c.Put("/user/"+orgID+"/"+userID, req, user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(orgID, userID string) error {
	return c.Delete("/user/"+orgID+"/"+userID, nil)
}

// GenerateUserOTPSecret generates a new OTP secret for the user.
// Returns the updated user including the otp_secret field.
func (c *Client) GenerateUserOTPSecret(orgID, userID string) (*User, error) {
	user := &User{}
	if err := c.Put("/user/"+orgID+"/"+userID+"/otp_secret", nil, user); err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByName finds a user by name within an organization.
func (c *Client) FindUserByName(orgID, name string) (*User, error) {
	users, err := c.ListUsers(orgID)
	if err != nil {
		return nil, err
	}
	for i := range users {
		if users[i].Name == name {
			return &users[i], nil
		}
	}
	return nil, fmt.Errorf("user %q not found in organization %q", name, orgID)
}
