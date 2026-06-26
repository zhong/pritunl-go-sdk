package pritunl

import "net/url"

// Status represents the response from GET /status.
type Status struct {
	HostCount     int      `json:"host_count"`
	ServersOnline int      `json:"servers_online"`
	HostsOnline   int      `json:"hosts_online"`
	ServerCount   int      `json:"server_count"`
	ServerVersion string   `json:"server_version"`
	PublicIP      string   `json:"public_ip"`
	UserCount     int      `json:"user_count"`
	UsersOnline   int      `json:"users_online"`
	OrgCount      int      `json:"org_count"`
	LocalNetworks []string `json:"local_networks"`
	CurrentHost   string   `json:"current_host"`
	Notification  string   `json:"notification"`
}

// Organization represents a Pritunl organization.
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UserCount int       `json:"user_count,omitempty"`
	AuthAPI   bool      `json:"auth_api,omitempty"`
	AuthToken string    `json:"auth_token,omitempty"`
	AuthSecret string   `json:"auth_secret,omitempty"`
}

// UserServer represents server attachment info embedded in a User response.
type UserServer struct {
	ID             string `json:"id"`
	ServerID       string `json:"server_id"`
	Name           string `json:"name"`
	Status         bool   `json:"status"`
	DeviceName     string `json:"device_name,omitempty"`
	Platform       string `json:"platform,omitempty"`
	RealAddress    string `json:"real_address,omitempty"`
	VirtAddress    string `json:"virt_address,omitempty"`
	VirtAddress6   string `json:"virt_address6,omitempty"`
	ConnectedSince string `json:"connected_since,omitempty"`
}

// User represents a Pritunl user.
type User struct {
	ID              string       `json:"id"`
	Organization    string       `json:"organization"`
	OrganizationName string      `json:"organization_name,omitempty"`
	Name            string       `json:"name"`
	Email           string       `json:"email,omitempty"`
	Groups          []string     `json:"groups,omitempty"`
	Type            string       `json:"type,omitempty"`
	AuthType        string       `json:"auth_type,omitempty"`
	Disabled        bool         `json:"disabled,omitempty"`
	OTPSecret       string       `json:"otp_secret,omitempty"`
	YubicoID        string       `json:"yubico_id,omitempty"`
	BypassSecondary bool         `json:"bypass_secondary,omitempty"`
	ClientToClient  bool         `json:"client_to_client,omitempty"`
	DNSMapping      string       `json:"dns_mapping,omitempty"`
	NetworkLinks    []string     `json:"network_links,omitempty"`
	Servers         []UserServer `json:"servers,omitempty"`
	Devices         []Device     `json:"devices,omitempty"`
}

// Device represents a registered device.
type Device struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Platform string `json:"platform,omitempty"`
}

// ServerRoute represents a route attached to a server.
type ServerRoute struct {
	Network string `json:"network"`
	NAT     bool   `json:"nat,omitempty"`
	NetGateway bool `json:"net_gateway,omitempty"`
}

// Server represents a Pritunl VPN server.
type Server struct {
	ID              string        `json:"id,omitempty"`
	Name            string        `json:"name"`
	Status          string        `json:"status,omitempty"`
	Network         string        `json:"network"`
	NetworkWG       string        `json:"network_wg,omitempty"`
	NetworkMode     string        `json:"network_mode,omitempty"`
	NetworkStart    string        `json:"network_start,omitempty"`
	NetworkEnd      string        `json:"network_end,omitempty"`
	Port            int           `json:"port,omitempty"`
	PortWG          int           `json:"port_wg,omitempty"`
	Protocol        string        `json:"protocol,omitempty"`
	WG              bool          `json:"wg,omitempty"`
	OTPAuth         bool          `json:"otp_auth,omitempty"`
	DHParamBits     int           `json:"dh_param_bits,omitempty"`
	Groups          []string      `json:"groups,omitempty"`
	BindAddress     string        `json:"bind_address,omitempty"`
	LocalNetworks   []string      `json:"local_networks,omitempty"`
	DNSservers      []string      `json:"dns_servers,omitempty"`
	DNSsuffix       string        `json:"dns_suffix,omitempty"`
	IPv6            bool          `json:"ipv6,omitempty"`
	Routes          []ServerRoute `json:"routes,omitempty"`
	Organizations   []string      `json:"organizations,omitempty"`
	Hosts           []string      `json:"hosts,omitempty"`
}

// KeyLink represents the temporary key download URI response from Pritunl.
// The backend returns paths; use FullURI() to get the client-importable URI.
type KeyLink struct {
	ID        string    `json:"id"`
	KeyURL    string    `json:"key_url"`
	KeyZipURL string    `json:"key_zip_url"`
	KeyOncURL string    `json:"key_onc_url"`
	ViewURL   string    `json:"view_url"`
	URIURL    string    `json:"uri_url"`
}

// FullURI returns the pritunl:// URI that can be imported into Pritunl clients.
// The host argument may be a plain hostname or a full URL (e.g. https://host:port).
func (k *KeyLink) FullURI(host string) string {
	if k.URIURL == "" {
		return ""
	}

	// Extract just the host part if a full URL was passed.
	if u, err := url.Parse(host); err == nil && u.Host != "" {
		host = u.Host
	}

	return "pritunl://" + host + k.URIURL
}

// PagedOrganizations is used when listing organizations with pagination.
type PagedOrganizations struct {
	Page          int            `json:"page"`
	PageTotal     int            `json:"page_total"`
	Organizations []Organization `json:"organizations"`
}

// PagedUsers is used when listing users with pagination.
type PagedUsers struct {
	Page        int    `json:"page"`
	PageTotal   int    `json:"page_total"`
	ServerCount int    `json:"server_count"`
	Users       []User `json:"users"`
}

// PagedServers is used when listing servers with pagination.
type PagedServers struct {
	Page      int      `json:"page"`
	PageTotal int      `json:"page_total"`
	Servers   []Server `json:"servers"`
}
