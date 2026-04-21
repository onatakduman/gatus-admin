package configio

// ── Maintenance (windows[] schema) ────────────────────────────────────────

type MaintenanceWindowsFileDoc struct {
	Maintenance MaintenanceBlock `yaml:"maintenance"`
}

type MaintenanceBlock struct {
	// Legacy single-window fields kept for migration tolerance.
	Start string `yaml:"start,omitempty"`
	End   string `yaml:"end,omitempty"`
	Every string `yaml:"every,omitempty"`
	// Preferred: list of windows.
	Windows []MWindow `yaml:"windows,omitempty"`
}

type MWindow struct {
	Name      string   `yaml:"name,omitempty"`
	Start     string   `yaml:"start,omitempty"`
	Duration  string   `yaml:"duration,omitempty"`
	Recurring string   `yaml:"recurring,omitempty"`
	Endpoints []string `yaml:"endpoints,omitempty"`
}

// ── Security (basic + OIDC) ───────────────────────────────────────────────

type SecurityFileDoc struct {
	Security *FullSecurity `yaml:"security,omitempty"`
}

type FullSecurity struct {
	Authentication *Authentication `yaml:"authentication,omitempty"`
	Authorization  *Authorization  `yaml:"authorization,omitempty"`
	// Legacy short-hand: security.basic — kept for back-compat.
	Basic *BasicAuth `yaml:"basic,omitempty"`
}

type Authentication struct {
	Type  string     `yaml:"type,omitempty"` // basic | oidc | none
	Basic *BasicAuth `yaml:"basic,omitempty"`
	OIDC  *OIDC      `yaml:"oidc,omitempty"`
}

type OIDC struct {
	IssuerURL    string   `yaml:"issuer-url,omitempty"`
	RedirectURL  string   `yaml:"redirect-url,omitempty"`
	ClientID     string   `yaml:"client-id,omitempty"`
	ClientSecret string   `yaml:"client-secret,omitempty"`
	Scopes       []string `yaml:"scopes,omitempty"`
}

type Authorization struct {
	DefaultRole  string   `yaml:"default-role,omitempty"`
	AdminGroups  []string `yaml:"admin-groups,omitempty"`
	ViewerGroups []string `yaml:"viewer-groups,omitempty"`
}

// ── Tunneling ─────────────────────────────────────────────────────────────

type TunnelingFileDoc struct {
	Tunneling map[string]Tunnel `yaml:"tunneling,omitempty"`
}

type Tunnel struct {
	Type       string `yaml:"type,omitempty"` // only SSH today
	Host       string `yaml:"host"`
	Port       int    `yaml:"port,omitempty"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"private-key,omitempty"`
}

// ── External endpoints ────────────────────────────────────────────────────

type ExternalEndpointsFileDoc struct {
	ExternalEndpoints []ExternalEndpoint `yaml:"external-endpoints"`
}

type ExternalEndpoint struct {
	Enabled   *bool      `yaml:"enabled,omitempty"`
	Name      string     `yaml:"name"`
	Group     string     `yaml:"group,omitempty"`
	Token     string     `yaml:"token"`
	Alerts    []Alert    `yaml:"alerts,omitempty"`
	Heartbeat *Heartbeat `yaml:"heartbeat,omitempty"`
}

type Heartbeat struct {
	Interval string `yaml:"interval"`
}

// ── Suites ────────────────────────────────────────────────────────────────

type SuitesFileDoc struct {
	Suites []Suite `yaml:"suites"`
}

type Suite struct {
	Enabled   *bool          `yaml:"enabled,omitempty"`
	Name      string         `yaml:"name"`
	Group     string         `yaml:"group,omitempty"`
	Interval  string         `yaml:"interval,omitempty"`
	Timeout   string         `yaml:"timeout,omitempty"`
	Context   map[string]any `yaml:"context,omitempty"`
	Endpoints []Endpoint     `yaml:"endpoints"`
}

// ── Endpoint advanced (DNS, SSH, OAuth2, IAP, TLS, badge) ─────────────────

// Extended client options to be spliced into Endpoint.Client via the same
// ClientOpts struct in types.go. We add fields here by shadowing the field
// names in wire protocol; they only serialise when non-zero.

type EndpointDNS struct {
	QueryType string `yaml:"query-type,omitempty"`
	QueryName string `yaml:"query-name,omitempty"`
}

type EndpointSSH struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// Badge thresholds (5 millisecond values) — nested inside endpoint.ui.badge.
type EndpointBadge struct {
	ResponseTime []int `yaml:"response-time,omitempty"`
}
