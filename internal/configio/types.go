package configio

type Endpoint struct {
	Enabled            *bool             `yaml:"enabled,omitempty"`
	Name               string            `yaml:"name"`
	Group              string            `yaml:"group,omitempty"`
	URL                string            `yaml:"url"`
	Method             string            `yaml:"method,omitempty"`
	Body               string            `yaml:"body,omitempty"`
	GraphQL            *bool             `yaml:"graphql,omitempty"`
	Headers            map[string]string `yaml:"headers,omitempty"`
	Interval           string            `yaml:"interval,omitempty"`
	Conditions         []string          `yaml:"conditions,omitempty"`
	Alerts             []Alert           `yaml:"alerts,omitempty"`
	Client             *ClientOpts       `yaml:"client,omitempty"`
	UI                 *EndpointUI       `yaml:"ui,omitempty"`
	DNS                *EndpointDNS      `yaml:"dns,omitempty"`
	SSH                *EndpointSSH      `yaml:"ssh,omitempty"`
	ExtraLabels        map[string]string `yaml:"extra-labels,omitempty"`
	MaintenanceWindows []MWindow         `yaml:"maintenance-windows,omitempty"`
}

type Alert struct {
	Type                    string `yaml:"type"`
	Description             string `yaml:"description,omitempty"`
	Enabled                 *bool  `yaml:"enabled,omitempty"`
	FailureThreshold        int    `yaml:"failure-threshold,omitempty"`
	SuccessThreshold        int    `yaml:"success-threshold,omitempty"`
	SendOnResolved          *bool  `yaml:"send-on-resolved,omitempty"`
	MinimumReminderInterval string `yaml:"minimum-reminder-interval,omitempty"`
}

type ClientOpts struct {
	Insecure             *bool          `yaml:"insecure,omitempty"`
	IgnoreRedirect       *bool          `yaml:"ignore-redirect,omitempty"`
	Timeout              string         `yaml:"timeout,omitempty"`
	DNSResolver          string         `yaml:"dns-resolver,omitempty"`
	ProxyURL             string         `yaml:"proxy-url,omitempty"`
	Network              string         `yaml:"network,omitempty"`
	Tunnel               string         `yaml:"tunnel,omitempty"`
	OAuth2               *OAuth2Config  `yaml:"oauth2,omitempty"`
	IdentityAwareProxy   *IAPConfig     `yaml:"identity-aware-proxy,omitempty"`
	TLS                  *TLSConfig     `yaml:"tls,omitempty"`
}

type OAuth2Config struct {
	TokenURL     string   `yaml:"token-url"`
	ClientID     string   `yaml:"client-id"`
	ClientSecret string   `yaml:"client-secret"`
	Scopes       []string `yaml:"scopes,omitempty"`
}

type IAPConfig struct {
	Audience string `yaml:"audience"`
}

type TLSConfig struct {
	CertificateFile string `yaml:"certificate-file,omitempty"`
	PrivateKeyFile  string `yaml:"private-key-file,omitempty"`
	Renegotiation   string `yaml:"renegotiation,omitempty"`
}

type EndpointUI struct {
	HideConditions              bool           `yaml:"hide-conditions,omitempty"`
	HideHostname                bool           `yaml:"hide-hostname,omitempty"`
	HidePort                    bool           `yaml:"hide-port,omitempty"`
	HideURL                     bool           `yaml:"hide-url,omitempty"`
	HideErrors                  bool           `yaml:"hide-errors,omitempty"`
	DontResolveFailedConditions bool           `yaml:"dont-resolve-failed-conditions,omitempty"`
	ResolveSuccessfulConditions bool           `yaml:"resolve-successful-conditions,omitempty"`
	Badge                       *EndpointBadge `yaml:"badge,omitempty"`
}

type EndpointsFileDoc struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

type UIFileDoc struct {
	UI UISettings `yaml:"ui"`
}

type UISettings struct {
	Title               string    `yaml:"title,omitempty"`
	Description         string    `yaml:"description,omitempty"`
	Header              string    `yaml:"header,omitempty"`
	Logo                string    `yaml:"logo,omitempty"`
	Link                string    `yaml:"link,omitempty"`
	DashboardHeading    string    `yaml:"dashboard-heading,omitempty"`
	DashboardSubheading string    `yaml:"dashboard-subheading,omitempty"`
	DarkMode            *bool     `yaml:"dark-mode,omitempty"`
	DefaultSortBy       string    `yaml:"default-sort-by,omitempty"`
	DefaultFilterBy     string    `yaml:"default-filter-by,omitempty"`
	Buttons             []UIButton `yaml:"buttons,omitempty"`
	CustomCSS           string    `yaml:"custom-css,omitempty"`
	Favicon             *Favicon  `yaml:"favicon,omitempty"`
}

type UIButton struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}

type Favicon struct {
	Default   string `yaml:"default,omitempty"`
	Size16x16 string `yaml:"size16x16,omitempty"`
	Size32x32 string `yaml:"size32x32,omitempty"`
}

type AlertingFileDoc struct {
	Alerting AlertingSettings `yaml:"alerting"`
}

type AlertingSettings struct {
	Ntfy     *NtfyConfig     `yaml:"ntfy,omitempty"`
	Slack    *WebhookConfig  `yaml:"slack,omitempty"`
	Discord  *WebhookConfig  `yaml:"discord,omitempty"`
	Telegram *TelegramConfig `yaml:"telegram,omitempty"`
	Email    *EmailConfig    `yaml:"email,omitempty"`
	Custom   *CustomConfig   `yaml:"custom,omitempty"`
}

type DefaultAlert struct {
	Enabled                 *bool  `yaml:"enabled,omitempty"`
	FailureThreshold        int    `yaml:"failure-threshold,omitempty"`
	SuccessThreshold        int    `yaml:"success-threshold,omitempty"`
	SendOnResolved          *bool  `yaml:"send-on-resolved,omitempty"`
	MinimumReminderInterval string `yaml:"minimum-reminder-interval,omitempty"`
}

type NtfyConfig struct {
	URL          string        `yaml:"url"`
	Topic        string        `yaml:"topic"`
	Token        string        `yaml:"token,omitempty"`
	Priority     int           `yaml:"priority,omitempty"`
	Click        string        `yaml:"click,omitempty"`
	DefaultAlert *DefaultAlert `yaml:"default-alert,omitempty"`
}

type WebhookConfig struct {
	WebhookURL   string        `yaml:"webhook-url"`
	DefaultAlert *DefaultAlert `yaml:"default-alert,omitempty"`
}

type TelegramConfig struct {
	Token        string        `yaml:"token"`
	ID           string        `yaml:"id"`
	DefaultAlert *DefaultAlert `yaml:"default-alert,omitempty"`
}

type EmailConfig struct {
	From         string        `yaml:"from"`
	To           string        `yaml:"to"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Username     string        `yaml:"username,omitempty"`
	Password     string        `yaml:"password,omitempty"`
	DefaultAlert *DefaultAlert `yaml:"default-alert,omitempty"`
}

type CustomConfig struct {
	URL          string            `yaml:"url"`
	Method       string            `yaml:"method,omitempty"`
	Headers      map[string]string `yaml:"headers,omitempty"`
	Body         string            `yaml:"body,omitempty"`
	DefaultAlert *DefaultAlert     `yaml:"default-alert,omitempty"`
}

type AnnouncementsFileDoc struct {
	Announcements []Announcement `yaml:"announcements"`
}

type Announcement struct {
	Timestamp string `yaml:"timestamp,omitempty"`
	Type      string `yaml:"type,omitempty"`
	Message   string `yaml:"message"`
	Archived  bool   `yaml:"archived,omitempty"`
}

type MaintenanceFileDoc struct {
	Maintenance *MaintenanceWindow `yaml:"maintenance,omitempty"`
}

type MaintenanceWindow struct {
	Start string `yaml:"start,omitempty"`
	End   string `yaml:"end,omitempty"`
	Every string `yaml:"every,omitempty"`
}

type CoreFileDoc struct {
	Web         *WebSettings     `yaml:"web,omitempty"`
	Storage     *StorageSettings `yaml:"storage,omitempty"`
	Security    *Security        `yaml:"security,omitempty"`
	Metrics     *bool            `yaml:"metrics,omitempty"`
	Concurrency int              `yaml:"concurrency,omitempty"`
}

type WebSettings struct {
	Address        string `yaml:"address,omitempty"`
	Port           int    `yaml:"port,omitempty"`
	ReadBufferSize int    `yaml:"read-buffer-size,omitempty"`
}

type StorageSettings struct {
	Type                    string `yaml:"type,omitempty"`
	Path                    string `yaml:"path,omitempty"`
	Caching                 *bool  `yaml:"caching,omitempty"`
	MaximumNumberOfResults  int    `yaml:"maximum-number-of-results,omitempty"`
	MaximumNumberOfEvents   int    `yaml:"maximum-number-of-events,omitempty"`
}

type Security struct {
	Basic *BasicAuth `yaml:"basic,omitempty"`
}

type BasicAuth struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}
