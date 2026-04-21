package configio

// AlertingProviders — the full Gatus roster of alerting providers.
// Each provider has its own struct; AlertingSettings (types.go) is extended to
// reference them. Every provider also carries `default-alert` and `overrides[]`.

type ProviderOverride struct {
	Group string `yaml:"group"`
	// Fields can differ per provider; use a generic map for the overridden
	// values so we can round-trip all of them without enumerating.
	Fields map[string]any `yaml:",inline"`
}

// Most webhook-style providers share the same simple shape.
// Reused by: slack, discord, mattermost, googlechat, teams-workflows,
// rocketchat, zapier, webex, n8n, squadcast (primarily), signl4-like.
type WebhookProvider struct {
	WebhookURL   string             `yaml:"webhook-url"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

// Slack adds channel + bot-name; kept separate for clarity.
type SlackProvider struct {
	WebhookURL   string             `yaml:"webhook-url"`
	Channel      string             `yaml:"channel,omitempty"`
	BotName      string             `yaml:"bot-name,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type DiscordProvider struct {
	WebhookURL     string             `yaml:"webhook-url"`
	Title          string             `yaml:"title,omitempty"`
	MessageContent string             `yaml:"message-content,omitempty"`
	DefaultAlert   *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides      []ProviderOverride `yaml:"overrides,omitempty"`
}

type EmailProvider struct {
	From         string             `yaml:"from"`
	Username     string             `yaml:"username,omitempty"`
	Password     string             `yaml:"password,omitempty"`
	Host         string             `yaml:"host"`
	Port         int                `yaml:"port"`
	To           string             `yaml:"to"`
	Client       *ClientOpts        `yaml:"client,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type TelegramProvider struct {
	Token        string             `yaml:"token"`
	ID           string             `yaml:"id"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type NtfyProvider struct {
	ServerURL    string             `yaml:"server-url,omitempty"`
	URL          string             `yaml:"url,omitempty"` // legacy field
	Topic        string             `yaml:"topic"`
	Token        string             `yaml:"token,omitempty"`
	Priority     int                `yaml:"priority,omitempty"`
	Click        string             `yaml:"click,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type PagerDutyProvider struct {
	IntegrationKey string             `yaml:"integration-key"`
	AutoResolve    *bool              `yaml:"auto-resolve,omitempty"`
	Severity       string             `yaml:"severity,omitempty"`
	DefaultAlert   *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides      []ProviderOverride `yaml:"overrides,omitempty"`
}

type OpsgenieProvider struct {
	APIKey       string             `yaml:"api-key"`
	Priority     string             `yaml:"priority,omitempty"`
	Tags         []string           `yaml:"tags,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type PushoverProvider struct {
	ApplicationToken string             `yaml:"application-token"`
	UserKey          string             `yaml:"user-key"`
	Priority         int                `yaml:"priority,omitempty"`
	DefaultAlert     *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides        []ProviderOverride `yaml:"overrides,omitempty"`
}

type GotifyProvider struct {
	ServerURL    string             `yaml:"server-url"`
	Token        string             `yaml:"token"`
	Priority     int                `yaml:"priority,omitempty"`
	Title        string             `yaml:"title,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type GitHubProvider struct {
	RepositoryURL string             `yaml:"repository-url"`
	Token         string             `yaml:"token"`
	DefaultAlert  *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides     []ProviderOverride `yaml:"overrides,omitempty"`
}

type GiteaProvider struct {
	RepositoryURL string             `yaml:"repository-url"`
	Token         string             `yaml:"token"`
	DefaultAlert  *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides     []ProviderOverride `yaml:"overrides,omitempty"`
}

type GitLabProvider struct {
	WebhookURL       string             `yaml:"webhook-url"`
	AuthorizationKey string             `yaml:"authorization-key"`
	Severity         string             `yaml:"severity,omitempty"`
	MonitoringTool   string             `yaml:"monitoring-tool,omitempty"`
	EnvironmentName  string             `yaml:"environment-name,omitempty"`
	Service          string             `yaml:"service,omitempty"`
	DefaultAlert     *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides        []ProviderOverride `yaml:"overrides,omitempty"`
}

type MatrixProvider struct {
	HomeserverURL string             `yaml:"homeserver-url"`
	UserID        string             `yaml:"user-id"`
	Password      string             `yaml:"password"`
	RoomID        string             `yaml:"room-id"`
	MsgType       string             `yaml:"msgtype,omitempty"`
	DefaultAlert  *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides     []ProviderOverride `yaml:"overrides,omitempty"`
}

type LineProvider struct {
	Token        string             `yaml:"token"`
	UserID       string             `yaml:"user-id"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type TwilioProvider struct {
	AccountSID   string             `yaml:"account-sid"`
	AuthToken    string             `yaml:"auth-token"`
	From         string             `yaml:"from"`
	To           string             `yaml:"to"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type VonageProvider struct {
	APIKey       string             `yaml:"api-key"`
	APISecret    string             `yaml:"api-secret"`
	From         string             `yaml:"from"`
	To           string             `yaml:"to"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type PlivoProvider struct {
	AuthID       string             `yaml:"auth-id"`
	AuthToken    string             `yaml:"auth-token"`
	From         string             `yaml:"from"`
	To           string             `yaml:"to"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type MessagebirdProvider struct {
	AccessKey    string             `yaml:"access-key"`
	Originator   string             `yaml:"originator"`
	Recipients   string             `yaml:"recipients"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type SendGridProvider struct {
	APIKey       string             `yaml:"api-key"`
	From         string             `yaml:"from"`
	To           string             `yaml:"to"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type AwsSESProvider struct {
	AccessKeyID     string             `yaml:"access-key-id,omitempty"`
	SecretAccessKey string             `yaml:"secret-access-key,omitempty"`
	Region          string             `yaml:"region"`
	From            string             `yaml:"from"`
	To              string             `yaml:"to"`
	DefaultAlert    *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides       []ProviderOverride `yaml:"overrides,omitempty"`
}

type SignalProvider struct {
	Number       string             `yaml:"number"`
	Recipients   []string           `yaml:"recipients"`
	URL          string             `yaml:"url,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type DatadogProvider struct {
	APIKey       string             `yaml:"api-key"`
	Site         string             `yaml:"site,omitempty"`
	Tags         []string           `yaml:"tags,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type NewRelicProvider struct {
	APIKey       string             `yaml:"api-key"`
	AccountID    string             `yaml:"account-id"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type SplunkProvider struct {
	HECURL       string             `yaml:"hec-url"`
	HECToken     string             `yaml:"hec-token"`
	Source       string             `yaml:"source,omitempty"`
	SourceType   string             `yaml:"source-type,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type SquadcastProvider struct {
	WebhookURL   string             `yaml:"webhook-url"`
	RoutingKey   string             `yaml:"routing-key"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type IlertProvider struct {
	APIKey       string             `yaml:"api-key"`
	IncidentKey  string             `yaml:"incident-key"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type IncidentIOProvider struct {
	RoutingKey   string             `yaml:"routing-key"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type Signl4Provider struct {
	TeamSecret   string             `yaml:"team-secret"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type HomeAssistantProvider struct {
	URL          string             `yaml:"url"`
	Token        string             `yaml:"token"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type IFTTTProvider struct {
	WebhookKey   string             `yaml:"webhook-key"`
	EventName    string             `yaml:"event-name"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type ClickupProvider struct {
	ListID       string             `yaml:"list-id"`
	Token        string             `yaml:"token"`
	APIURL       string             `yaml:"api-url,omitempty"`
	Assignees    []string           `yaml:"assignees,omitempty"`
	Status       string             `yaml:"status,omitempty"`
	Priority     string             `yaml:"priority,omitempty"`
	NotifyAll    *bool              `yaml:"notify-all,omitempty"`
	Name         string             `yaml:"name,omitempty"`
	Content      string             `yaml:"content,omitempty"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type ZulipProvider struct {
	APIURL       string             `yaml:"api-url"`
	APIKey       string             `yaml:"api-key"`
	BotName      string             `yaml:"bot-name"`
	Organization string             `yaml:"organization,omitempty"`
	To           string             `yaml:"to"`
	DefaultAlert *DefaultAlert      `yaml:"default-alert,omitempty"`
	Overrides    []ProviderOverride `yaml:"overrides,omitempty"`
}

type JenkinsProvider = WebhookProvider // filler to avoid unused import warnings

// FullAlertingSettings is the union of every provider Gatus supports.
// Only non-nil pointers are written, so an "unconfigured" provider vanishes
// from disk entirely — matching Gatus behaviour where missing = disabled.
type FullAlertingSettings struct {
	// Existing simple fields preserved for backward compat:
	Ntfy     *NtfyProvider    `yaml:"ntfy,omitempty"`
	Slack    *SlackProvider   `yaml:"slack,omitempty"`
	Discord  *DiscordProvider `yaml:"discord,omitempty"`
	Email    *EmailProvider   `yaml:"email,omitempty"`
	Telegram *TelegramProvider `yaml:"telegram,omitempty"`
	Custom   *CustomConfig    `yaml:"custom,omitempty"`

	// Expanded roster:
	AwsSES        *AwsSESProvider        `yaml:"awsses,omitempty"`
	Clickup       *ClickupProvider       `yaml:"clickup,omitempty"`
	Datadog       *DatadogProvider       `yaml:"datadog,omitempty"`
	Gitea         *GiteaProvider         `yaml:"gitea,omitempty"`
	GitHub        *GitHubProvider        `yaml:"github,omitempty"`
	GitLab        *GitLabProvider        `yaml:"gitlab,omitempty"`
	GoogleChat    *WebhookProvider       `yaml:"googlechat,omitempty"`
	Gotify        *GotifyProvider        `yaml:"gotify,omitempty"`
	HomeAssistant *HomeAssistantProvider `yaml:"homeassistant,omitempty"`
	IFTTT         *IFTTTProvider         `yaml:"ifttt,omitempty"`
	Ilert         *IlertProvider         `yaml:"ilert,omitempty"`
	IncidentIO    *IncidentIOProvider    `yaml:"incident-io,omitempty"`
	Line          *LineProvider          `yaml:"line,omitempty"`
	Matrix        *MatrixProvider        `yaml:"matrix,omitempty"`
	Mattermost    *WebhookProvider       `yaml:"mattermost,omitempty"`
	Messagebird   *MessagebirdProvider   `yaml:"messagebird,omitempty"`
	N8n           *WebhookProvider       `yaml:"n8n,omitempty"`
	NewRelic      *NewRelicProvider      `yaml:"newrelic,omitempty"`
	Opsgenie      *OpsgenieProvider      `yaml:"opsgenie,omitempty"`
	PagerDuty     *PagerDutyProvider     `yaml:"pagerduty,omitempty"`
	Plivo         *PlivoProvider         `yaml:"plivo,omitempty"`
	Pushover      *PushoverProvider      `yaml:"pushover,omitempty"`
	RocketChat    *WebhookProvider       `yaml:"rocketchat,omitempty"`
	SendGrid      *SendGridProvider      `yaml:"sendgrid,omitempty"`
	Signal        *SignalProvider        `yaml:"signal,omitempty"`
	Signl4        *Signl4Provider        `yaml:"signl4,omitempty"`
	Splunk        *SplunkProvider        `yaml:"splunk,omitempty"`
	Squadcast     *SquadcastProvider     `yaml:"squadcast,omitempty"`
	Teams         *WebhookProvider       `yaml:"teams-workflows,omitempty"`
	Twilio        *TwilioProvider        `yaml:"twilio,omitempty"`
	Vonage        *VonageProvider        `yaml:"vonage,omitempty"`
	Webex         *WebhookProvider       `yaml:"webex,omitempty"`
	Zapier        *WebhookProvider       `yaml:"zapier,omitempty"`
	Zulip         *ZulipProvider         `yaml:"zulip,omitempty"`
}

type FullAlertingFileDoc struct {
	Alerting FullAlertingSettings `yaml:"alerting"`
}
