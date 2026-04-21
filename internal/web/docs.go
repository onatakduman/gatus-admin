package web

import "html/template"

// docsURL maps an active sidebar key to the upstream Gatus docs URL.
func docsURL(active string) string {
	switch active {
	case "endpoints":
		return "https://gatus.io/docs/endpoints"
	case "external-endpoints":
		return "https://gatus.io/docs/external-endpoints"
	case "suites":
		return "https://gatus.io/docs/suites"
	case "alerting":
		return "https://gatus.io/docs/alerting"
	case "ui":
		return "https://gatus.io/docs/appearance"
	case "announcements":
		return "https://gatus.io/docs/announcements"
	case "maintenance":
		return "https://gatus.io/docs/maintenance"
	case "security":
		return "https://gatus.io/docs/security"
	case "tunneling":
		return "https://gatus.io/docs/endpoints"
	case "core":
		return "https://gatus.io/docs/configuration"
	case "advanced":
		return "https://gatus.io/docs"
	}
	return "https://gatus.io/docs"
}

// docsLabel is the short text shown next to the link.
func docsLabel(active string) string {
	if u := docsURL(active); u != "" {
		return u
	}
	return "Gatus docs"
}

// docsContent returns pre-rendered HTML documentation for an active sidebar tab.
// Content is hand-condensed from https://gatus.io/docs/<page> into a quick
// reference plus a typical YAML snippet.
func docsContent(active string) template.HTML {
	if html, ok := docsHTML[active]; ok {
		return html
	}
	return ""
}

var docsHTML = map[string]template.HTML{
	"endpoints": template.HTML(`
<p>An <strong>endpoint</strong> is a single worker that monitors a target at every configured interval. Each endpoint must satisfy its <em>conditions</em> to be healthy.</p>
<table>
<thead><tr><th>Field</th><th>Required</th><th>Description</th></tr></thead>
<tbody>
<tr><td><code>name</code></td><td>✓</td><td>Display name on the dashboard.</td></tr>
<tr><td><code>group</code></td><td>—</td><td>Cluster endpoints visually.</td></tr>
<tr><td><code>url</code></td><td>✓</td><td><code>http(s)://</code>, <code>tcp://</code>, <code>icmp://</code>, <code>dns://</code>, <code>ssh://</code>, <code>starttls://</code>.</td></tr>
<tr><td><code>method</code></td><td>—</td><td>HTTP verb. Default <code>GET</code>.</td></tr>
<tr><td><code>interval</code></td><td>—</td><td>How often (e.g. <code>30s</code>, <code>5m</code>). Default <code>60s</code>.</td></tr>
<tr><td><code>conditions</code></td><td>✓</td><td>List of expressions using placeholders.</td></tr>
<tr><td><code>alerts</code></td><td>—</td><td>Per-endpoint alert rules.</td></tr>
</tbody>
</table>
<h4>Common condition placeholders</h4>
<ul>
<li><code>[STATUS]</code> — HTTP status code</li>
<li><code>[RESPONSE_TIME]</code> — duration in ms</li>
<li><code>[BODY].field</code> — JSONPath into response</li>
<li><code>[CONNECTED]</code> — TCP/ICMP success</li>
<li><code>[CERTIFICATE_EXPIRATION]</code> — TLS cert lifetime</li>
<li><code>[DNS_RCODE]</code> — DNS response code</li>
</ul>
<h4>Example</h4>
<pre><code>endpoints:
  - name: api
    group: Production
    url: https://api.example.com/health
    interval: 1m
    conditions:
      - "[STATUS] == 200"
      - "[RESPONSE_TIME] &lt; 500"
      - "[BODY].status == UP"
    alerts:
      - type: slack
        failure-threshold: 3
        send-on-resolved: true
</code></pre>
<p><a href="https://gatus.io/docs/endpoints" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"external-endpoints": template.HTML(`
<p><strong>External endpoints</strong> are push-based: your scripts/cron jobs POST status updates to Gatus instead of being polled.</p>
<table>
<thead><tr><th>Field</th><th>Required</th><th>Description</th></tr></thead>
<tbody>
<tr><td><code>name</code></td><td>✓</td><td>Display name.</td></tr>
<tr><td><code>token</code></td><td>✓</td><td>Bearer token clients send to authenticate pushes.</td></tr>
<tr><td><code>heartbeat.interval</code></td><td>—</td><td>Auto-fail if no push received within this duration.</td></tr>
</tbody>
</table>
<h4>Pushing a status</h4>
<pre><code>curl -H "Authorization: Bearer YOUR_TOKEN" \
  "https://status.example.com/api/v1/endpoints/Group_Name/external?success=true&duration=120ms"
</code></pre>
<p><a href="https://gatus.io/docs/external-endpoints" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"suites": template.HTML(`
<p><strong>Suites (alpha)</strong> are sequential multi-step health checks that share a context (e.g. login → fetch → logout). Each step is an endpoint that can <code>store</code> values for the next step.</p>
<table>
<thead><tr><th>Field</th><th>Required</th><th>Description</th></tr></thead>
<tbody>
<tr><td><code>name</code></td><td>✓</td><td>Suite identifier.</td></tr>
<tr><td><code>interval</code></td><td>—</td><td>How often the whole suite runs. Default <code>10m</code>.</td></tr>
<tr><td><code>timeout</code></td><td>—</td><td>Total suite timeout. Default <code>5m</code>.</td></tr>
<tr><td><code>context</code></td><td>—</td><td>Initial values, accessible via <code>[CONTEXT].key</code>.</td></tr>
<tr><td><code>endpoints</code></td><td>✓</td><td>Ordered steps. Add via Advanced YAML for now.</td></tr>
</tbody>
</table>
<p><a href="https://gatus.io/docs/suites" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"alerting": template.HTML(`
<p>Alerting providers send notifications when an endpoint fails. Configure providers globally; reference them per-endpoint with <code>alerts: - type: &lt;provider&gt;</code>.</p>
<h4>Common fields per provider</h4>
<ul>
<li><code>default-alert</code> — failure-threshold, success-threshold, send-on-resolved.</li>
<li><code>overrides[]</code> — per-group field overrides.</li>
</ul>
<h4>Example: ntfy + Slack</h4>
<pre><code>alerting:
  ntfy:
    server-url: https://ntfy.sh
    topic: my-alerts
    default-alert:
      failure-threshold: 2
      success-threshold: 2
      send-on-resolved: true
  slack:
    webhook-url: https://hooks.slack.com/services/XXX
</code></pre>
<p>Gatus supports 35+ providers. <a href="https://gatus.io/docs/alerting" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"ui": template.HTML(`
<p>Settings on this page customize the <strong>public</strong> Gatus status page. The admin panel itself has its own branding under <a href="/admin/brand">Admin Brand</a>.</p>
<table>
<thead><tr><th>Field</th><th>Description</th></tr></thead>
<tbody>
<tr><td><code>title</code></td><td>Browser tab title.</td></tr>
<tr><td><code>header</code></td><td>Wordmark at the top of the page.</td></tr>
<tr><td><code>logo</code></td><td>Logo URL (can point to <code>/admin/assets/…</code>).</td></tr>
<tr><td><code>dashboard-heading</code> / <code>-subheading</code></td><td>Hero text.</td></tr>
<tr><td><code>dark-mode</code></td><td>Default theme.</td></tr>
<tr><td><code>buttons[]</code></td><td>Extra header nav buttons.</td></tr>
<tr><td><code>custom-css</code></td><td>Inline CSS injected into the page.</td></tr>
</tbody>
</table>
<p><a href="https://gatus.io/docs/appearance" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"brand": template.HTML(`
<p>This panel's own visual identity. Independent from the public Gatus UI but can <em>inherit</em> the name/tagline from <code>ui.header</code> / <code>ui.description</code> when those are blank here.</p>
<ul>
<li><strong>Logo (full)</strong> — shown on the login screen.</li>
<li><strong>Logo (square)</strong> — shown in the sidebar.</li>
<li><strong>Favicon</strong> — admin browser tab icon.</li>
<li><strong>Colors</strong> — accent, background, surface, border, text.</li>
<li><strong>Fonts</strong> — Google Font name. Empty = system stack.</li>
</ul>
<p>Uploaded assets live under <code>config/brand/assets/</code> and are served from <code>/admin/assets/&lt;name&gt;</code>.</p>
`),

	"announcements": template.HTML(`
<p>Banner messages displayed at the top of the public status page.</p>
<table>
<thead><tr><th>Field</th><th>Description</th></tr></thead>
<tbody>
<tr><td><code>timestamp</code></td><td>RFC3339 UTC. Blank = now.</td></tr>
<tr><td><code>type</code></td><td><code>outage</code> · <code>warning</code> · <code>information</code> · <code>operational</code> · <code>none</code> — controls colour.</td></tr>
<tr><td><code>message</code></td><td>Markdown-formatted body.</td></tr>
<tr><td><code>archived</code></td><td>Hide from main banner; keep in history.</td></tr>
</tbody>
</table>
<p><a href="https://gatus.io/docs/announcements" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"maintenance": template.HTML(`
<p>Scheduled downtime windows. Alerts are suppressed inside an active window.</p>
<table>
<thead><tr><th>Field</th><th>Description</th></tr></thead>
<tbody>
<tr><td><code>start</code></td><td>RFC3339 UTC start time.</td></tr>
<tr><td><code>duration</code></td><td>e.g. <code>1h</code>, <code>30m</code>.</td></tr>
<tr><td><code>recurring</code></td><td><code>none</code> · <code>daily</code> · <code>weekly</code> · <code>monthly</code>.</td></tr>
<tr><td><code>endpoints</code></td><td>Endpoint names or glob patterns this window applies to.</td></tr>
</tbody>
</table>
<p><a href="https://gatus.io/docs/maintenance" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"security": template.HTML(`
<p>Protect the public Gatus page with HTTP Basic Auth or OIDC.</p>
<h4>Basic auth</h4>
<pre><code>security:
  authentication:
    type: basic
    basic:
      username: alice
      password: $2a$12$...   # bcrypt hash
</code></pre>
<h4>OIDC</h4>
<pre><code>security:
  authentication:
    type: oidc
    oidc:
      issuer-url: https://accounts.google.com
      redirect-url: https://status.example.com/authorization-code/callback
      client-id: ...
      client-secret: ...
      scopes: [openid, profile, email]
  authorization:
    admin-groups: ["admin@example.com"]
</code></pre>
<p><a href="https://gatus.io/docs/security" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"tunneling": template.HTML(`
<p>SSH tunnel definitions reusable across endpoints. Each tunnel is keyed by name; reference from an endpoint as <code>client.tunnel: bastion-prod</code>.</p>
<pre><code>tunneling:
  bastion-prod:
    type: SSH
    host: bastion.example.com
    port: 22
    username: gatus
    private-key: |
      -----BEGIN OPENSSH PRIVATE KEY-----
      ...
      -----END OPENSSH PRIVATE KEY-----
</code></pre>
<p><a href="https://gatus.io/docs/endpoints" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"core": template.HTML(`
<p>Server, storage, and global settings. Most setups never touch these.</p>
<table>
<thead><tr><th>Section</th><th>Common keys</th></tr></thead>
<tbody>
<tr><td><code>web</code></td><td><code>address</code>, <code>port</code> — Gatus listen socket.</td></tr>
<tr><td><code>storage</code></td><td><code>type</code> (memory/sqlite/postgres), <code>path</code>, <code>maximum-number-of-results</code>.</td></tr>
<tr><td><code>metrics</code></td><td>Boolean. Exposes <code>/metrics</code> for Prometheus scraping.</td></tr>
<tr><td><code>concurrency</code></td><td>Max parallel health checks. <code>0</code> = unlimited.</td></tr>
</tbody>
</table>
<p><a href="https://gatus.io/docs/configuration" target="_blank" rel="noopener">Full docs →</a></p>
`),

	"advanced": template.HTML(`
<p>Direct YAML editor for every config file. Useful for fields the form-based tabs don't cover yet (e.g. suite endpoints, OAuth2 client config, IAP audience, advanced TLS).</p>
<p>Each file is validated as YAML before saving — invalid input is rejected.</p>
<p><a href="https://gatus.io/docs" target="_blank" rel="noopener">Gatus docs →</a></p>
`),
}
