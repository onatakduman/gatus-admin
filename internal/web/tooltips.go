package web

// Tooltips maps a semantic field key to a short explanation pulled from
// https://gatus.io/docs/. Used by the `tooltip` template func.
var Tooltips = map[string]string{
	// Endpoint
	"endpoint.name":                    "Unique name for this endpoint. Shown on the status dashboard.",
	"endpoint.enabled":                 "Disable to stop checking without deleting the configuration.",
	"endpoint.group":                   "Group name used to cluster endpoints on the dashboard.",
	"endpoint.url":                     "Target address. http(s)://, tcp://, icmp://, dns://, ssh://, starttls:// are supported.",
	"endpoint.method":                  "HTTP method. Defaults to GET.",
	"endpoint.interval":                "How often to run this health check. Duration: 30s, 1m, 5m, 1h.",
	"endpoint.body":                    "Request body to send with the check.",
	"endpoint.graphql":                 "If true, the body is automatically wrapped as a GraphQL query payload.",
	"endpoint.headers":                 "HTTP headers to send, one per line as 'Key: Value'.",
	"endpoint.conditions":              "Each line must evaluate true for healthy. Placeholders: [STATUS], [RESPONSE_TIME], [BODY].x, [CONNECTED], [CERTIFICATE_EXPIRATION], [DOMAIN_EXPIRATION], [IP], [DNS_RCODE].",
	"endpoint.client.insecure":         "Skip TLS certificate verification. Useful for self-signed certs.",
	"endpoint.client.ignore-redirect":  "Do not follow HTTP 3xx redirects.",
	"endpoint.client.timeout":          "Per-request timeout, e.g. 10s.",
	"endpoint.client.dns-resolver":     "Custom DNS resolver, e.g. tcp://8.8.8.8:53.",
	"endpoint.client.proxy-url":        "Route requests through this HTTPS proxy.",
	"endpoint.client.tunnel":           "Reference an SSH tunnel defined under tunneling.*",
	"endpoint.client.network":          "ICMP protocol: ip, ip4, or ip6.",
	"endpoint.dns.query-type":          "DNS record type: A, AAAA, MX, NS, TXT, SRV, CNAME, PTR.",
	"endpoint.dns.query-name":          "Domain name to query.",
	"endpoint.ssh.username":            "SSH username.",
	"endpoint.ssh.password":            "SSH password (use tunnel.private-key instead where possible).",
	"endpoint.ui.hide-hostname":        "Hide the hostname from the public page.",
	"endpoint.ui.hide-url":             "Hide the full URL from the public page.",
	"endpoint.ui.hide-errors":          "Hide error messages from the public page (show only up/down).",
	"endpoint.ui.badge.response-time":  "Response-time thresholds for the SLA badge. 5 millisecond values.",
	"endpoint.extra-labels":            "Key-value pairs added to the Prometheus metrics for this endpoint.",
	"endpoint.maintenance-windows":     "Per-endpoint scheduled downtime. Alerts are suppressed during windows.",
	"endpoint.alerts":                  "Per-endpoint alert rules. Each references a global alerting provider by type.",
	"endpoint.always-run":              "(Suites only) Run this step even if a prior step failed.",

	// Alert (per-endpoint alert block)
	"alert.type":                       "Provider name (ntfy, slack, pagerduty, etc). Must match a configured alerting provider.",
	"alert.enabled":                    "Activate this specific alert rule.",
	"alert.description":                "Human-readable message included in the notification.",
	"alert.failure-threshold":          "Consecutive failures before triggering the alert.",
	"alert.success-threshold":          "Consecutive successes before resolving the alert.",
	"alert.send-on-resolved":           "Send a notification when the endpoint recovers.",
	"alert.minimum-reminder-interval":  "Suppress repeated notifications within this duration (min 5m).",

	// Alerting — ntfy
	"alerting.ntfy.server-url": "e.g. https://ntfy.sh or a self-hosted ntfy instance.",
	"alerting.ntfy.topic":      "Channel name subscribers watch.",
	"alerting.ntfy.token":      "Optional ntfy access token for protected topics.",
	"alerting.ntfy.priority":   "1-5 (min-max). Higher = more attention.",
	"alerting.ntfy.click":      "URL opened when the user taps the notification.",

	// Alerting — Slack / Discord / Teams / Mattermost / generic webhooks
	"alerting.slack.webhook-url":             "Incoming webhook URL from Slack.",
	"alerting.discord.webhook-url":           "Incoming webhook URL from Discord.",
	"alerting.mattermost.webhook-url":        "Incoming webhook URL from Mattermost.",
	"alerting.googlechat.webhook-url":        "Incoming webhook URL from Google Chat Spaces.",
	"alerting.teams-workflows.webhook-url":   "Microsoft Teams Workflow webhook URL.",
	"alerting.rocketchat.webhook-url":        "Incoming webhook URL from Rocket.Chat.",
	"alerting.zapier.webhook-url":            "Zapier catch-hook URL.",
	"alerting.webex.webhook-url":             "Cisco Webex incoming webhook URL.",

	// Alerting — Telegram
	"alerting.telegram.token":   "Bot token from @BotFather.",
	"alerting.telegram.chat-id": "Chat/channel ID. Use @raw_data_bot to find it.",

	// Alerting — Email
	"alerting.email.from":     "Sender address, e.g. alerts@example.com.",
	"alerting.email.username": "SMTP username. Defaults to the From address.",
	"alerting.email.password": "SMTP password or app-specific password.",
	"alerting.email.host":     "SMTP host, e.g. smtp.gmail.com.",
	"alerting.email.port":     "SMTP port. Common: 587 (STARTTLS), 465 (TLS).",
	"alerting.email.to":       "Recipient addresses (comma-separated).",

	// Alerting — PagerDuty / Opsgenie / Pushover etc.
	"alerting.pagerduty.integration-key": "PagerDuty Events v2 integration key.",
	"alerting.pagerduty.severity":        "critical | error | warning | info.",
	"alerting.opsgenie.api-key":          "Opsgenie API key.",
	"alerting.opsgenie.priority":         "P1 (critical) – P5 (lowest). Default: P3.",
	"alerting.pushover.application-token": "Pushover application API token.",
	"alerting.pushover.user-key":         "Pushover user key.",
	"alerting.pushover.priority":         "-2 (silent) to 2 (emergency). Default 0.",
	"alerting.gotify.server-url":         "Self-hosted Gotify instance URL.",
	"alerting.gotify.token":              "Gotify application token.",
	"alerting.gotify.priority":           "Gotify priority, 1-10.",

	// Alerting — Issue trackers
	"alerting.github.repository-url": "https://github.com/user/repo",
	"alerting.github.token":          "PAT with issue R/W + metadata R on the target repo.",
	"alerting.gitea.repository-url":  "https://gitea.example.com/user/repo",
	"alerting.gitea.token":           "PAT with issue R/W + metadata R on the target repo.",
	"alerting.gitlab.webhook-url":    "Alert webhook URL from GitLab.",
	"alerting.gitlab.authorization-key": "Authorization key from the GitLab webhook config.",

	// Alerting — Chat apps / platforms
	"alerting.matrix.homeserver-url": "e.g. https://matrix.org.",
	"alerting.matrix.user-id":        "Matrix user ID, e.g. @alice:matrix.org.",
	"alerting.matrix.password":       "Matrix account password.",
	"alerting.matrix.room-id":        "Room ID or alias to post into.",
	"alerting.line.token":            "LINE Messaging API channel access token.",
	"alerting.line.user-id":          "LINE target user ID.",
	"alerting.zulip.api-url":         "Zulip API URL.",
	"alerting.zulip.api-key":         "Zulip bot API key.",
	"alerting.zulip.bot-name":        "Zulip bot email/name used to post.",
	"alerting.zulip.to":              "Zulip topic or user to notify.",

	// Alerting — Incident/monitoring platforms
	"alerting.datadog.api-key":       "Datadog API key.",
	"alerting.datadog.site":          "datadoghq.com or datadoghq.eu.",
	"alerting.newrelic.api-key":      "New Relic user key.",
	"alerting.newrelic.account-id":   "New Relic account ID.",
	"alerting.splunk.hec-url":        "Splunk HTTP Event Collector URL.",
	"alerting.splunk.hec-token":      "Splunk HEC token.",
	"alerting.squadcast.webhook-url": "Squadcast incoming webhook URL.",
	"alerting.squadcast.routing-key": "Squadcast routing key.",
	"alerting.ilert.api-key":         "ilert API key.",
	"alerting.ilert.incident-key":    "ilert incident de-duplication key.",
	"alerting.incident-io.routing-key": "incident.io routing key.",
	"alerting.signl4.team-secret":    "SIGNL4 team secret.",
	"alerting.homeassistant.url":     "Home Assistant instance URL, e.g. http://homeassistant:8123.",
	"alerting.homeassistant.token":   "Home Assistant long-lived access token.",

	// Alerting — SMS / email providers
	"alerting.twilio.account-sid": "Twilio account SID.",
	"alerting.twilio.auth-token":  "Twilio auth token.",
	"alerting.twilio.from":        "Sender phone number.",
	"alerting.twilio.to":          "Recipient phone numbers (comma-separated).",
	"alerting.vonage.api-key":     "Vonage API key.",
	"alerting.vonage.api-secret":  "Vonage API secret.",
	"alerting.plivo.auth-id":      "Plivo Auth ID.",
	"alerting.plivo.auth-token":   "Plivo Auth Token.",
	"alerting.messagebird.access-key": "MessageBird access key.",
	"alerting.messagebird.originator": "Sender identifier (phone or alphanumeric).",
	"alerting.messagebird.recipients": "Recipient phone numbers, comma-separated.",
	"alerting.sendgrid.api-key":   "SendGrid API key.",
	"alerting.awsses.region":      "AWS region, e.g. eu-west-1.",
	"alerting.awsses.from":        "Verified sender address.",
	"alerting.awsses.to":          "Recipient addresses (comma-separated).",
	"alerting.signal.number":      "Sender phone number registered with Signal-Cli.",
	"alerting.signal.recipients":  "Recipient phone numbers (+E.164).",
	"alerting.signal.url":         "Signal REST API endpoint.",

	// Alerting — Automation
	"alerting.ifttt.webhook-key": "IFTTT Maker key.",
	"alerting.ifttt.event-name":  "IFTTT event name.",
	"alerting.n8n.webhook-url":   "n8n webhook node URL.",
	"alerting.clickup.list-id":   "ClickUp list ID to create tasks in.",
	"alerting.clickup.token":     "ClickUp API token.",

	// Alerting — Custom
	"alerting.custom.url":     "Any webhook URL. Placeholders in body/headers: [ENDPOINT_NAME], [ENDPOINT_GROUP], [ENDPOINT_URL], [ALERT_DESCRIPTION].",
	"alerting.custom.method":  "HTTP method. Default: POST.",
	"alerting.custom.headers": "HTTP headers to send with the webhook call.",
	"alerting.custom.body":    "Template body with placeholders. JSON or form-encoded.",

	// UI (public page)
	"ui.title":                 "Browser tab title.",
	"ui.description":           "Meta description for SEO.",
	"ui.dashboard-heading":     "Large heading on the dashboard.",
	"ui.dashboard-subheading":  "Subtitle beneath the heading.",
	"ui.header":                "Short wordmark at the top of the public page.",
	"ui.logo":                  "Public-page logo URL (can point to /admin/assets/…).",
	"ui.link":                  "Destination when the logo is clicked.",
	"ui.dark-mode":             "Default theme for first-time visitors.",
	"ui.default-sort-by":       "Initial sort: name / group / health.",
	"ui.default-filter-by":     "Initial filter: none / failing / unstable.",
	"ui.buttons":               "Extra nav buttons in the header.",
	"ui.favicon.default":       "Path to the default .ico favicon.",
	"ui.favicon.size16x16":     "16x16 favicon path.",
	"ui.favicon.size32x32":     "32x32 favicon path.",
	"ui.custom-css":            "Custom CSS injected into the public page AND mirrored into this admin panel so both stay visually in sync.",

	// Admin brand
	"brand.name":               "Admin panel display name. Leave empty to inherit from Gatus UI header.",
	"brand.tagline":            "Small subtitle shown on the login page.",
	"brand.inherit":            "Pull empty Name/Tagline from Gatus ui.yaml (ui.header / ui.description).",
	"brand.logo_full":          "Full wordmark logo shown on the login page. PNG/JPG/SVG/WebP, max 2MB.",
	"brand.logo_square":        "Square icon shown in the sidebar. PNG/JPG/SVG/WebP/ICO, max 2MB.",
	"brand.favicon":            "Admin panel tab icon. PNG/ICO recommended.",
	"brand.colors.accent":      "Primary accent color (buttons, links, active nav).",
	"brand.colors.accent_hover": "Hover state of the accent color.",
	"brand.fonts.heading":      "Google Fonts name for headings and logo. Empty = system sans-serif.",
	"brand.fonts.body":         "Google Fonts name for body text. Empty = system monospace.",

	// Storage
	"storage.type":                          "memory = RAM only, sqlite = single file, postgres = remote DB.",
	"storage.path":                          "File path (sqlite) or postgres://user:pass@host/db URL.",
	"storage.caching":                       "Write-through cache for sqlite/postgres. Faster reads.",
	"storage.maximum-number-of-results":     "How many historical results to keep per endpoint.",
	"storage.maximum-number-of-events":      "How many events (up/down transitions) to keep.",

	// Web
	"web.address":          "Listen address. 0.0.0.0 binds all interfaces.",
	"web.port":             "Listen port. Default 8080.",
	"web.read-buffer-size": "Max size of an HTTP request header block in bytes.",
	"web.tls.certificate-file": "Path to PEM-encoded TLS certificate inside the Gatus container.",
	"web.tls.private-key-file": "Path to PEM-encoded TLS private key inside the Gatus container.",

	// Security
	"security.basic.username":  "Username for HTTP Basic Auth on the public page.",
	"security.basic.password":  "Password for HTTP Basic Auth on the public page.",
	"security.oidc.issuer-url": "Identity provider issuer URL, e.g. https://accounts.google.com.",
	"security.oidc.redirect-url": "OIDC callback URL registered at the provider.",
	"security.oidc.client-id":  "OIDC client ID.",
	"security.oidc.client-secret": "OIDC client secret.",
	"security.oidc.scopes":     "Requested OIDC scopes (comma or space separated).",

	// Global
	"metrics":                    "Expose Prometheus metrics at /metrics.",
	"concurrency":                "Max parallel health checks. 0 = unlimited.",
	"skip-invalid-config-update": "Ignore invalid config reloads instead of crashing.",

	// Announcements
	"announcement.type":      "outage | warning | information | operational | none. Affects banner color.",
	"announcement.timestamp": "RFC3339 UTC. Leave blank to use current time.",
	"announcement.message":   "Markdown-formatted banner text.",
	"announcement.archived":  "Hidden from the primary banner; shown in Past Announcements.",

	// Maintenance
	"maintenance.window.name":      "Friendly label for the window (optional).",
	"maintenance.window.start":     "RFC3339 UTC start time.",
	"maintenance.window.duration":  "Duration, e.g. 1h, 30m.",
	"maintenance.window.recurring": "none | daily | weekly | monthly.",
	"maintenance.window.endpoints": "Endpoint names or glob patterns this window applies to.",

	// Tunneling
	"tunneling.type":        "Only 'SSH' is supported today.",
	"tunneling.host":        "SSH bastion host.",
	"tunneling.port":        "SSH port. Default 22.",
	"tunneling.username":    "SSH username.",
	"tunneling.password":    "SSH password (use private-key when possible).",
	"tunneling.private-key": "PEM-encoded private key.",

	// Suites
	"suite.name":     "Unique suite identifier.",
	"suite.interval": "How often the suite runs. Default 10m.",
	"suite.timeout":  "Total suite timeout. Default 5m.",
	"suite.context":  "Key-value map available to all endpoints as [CONTEXT].key.",

	// External endpoints
	"external-endpoint.name":      "Unique name.",
	"external-endpoint.group":     "Group name on the status page.",
	"external-endpoint.token":     "Bearer token clients POST back to /api/v1/endpoints/{key}/external.",
	"external-endpoint.heartbeat.interval": "Alert if no push received within this duration.",
}

// TooltipFor returns the tooltip text for a field key, or empty string if unknown.
func TooltipFor(key string) string { return Tooltips[key] }
