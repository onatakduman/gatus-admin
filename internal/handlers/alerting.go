package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

func (h *Handlers) AlertingGet(w http.ResponseWriter, r *http.Request) {
	doc := &configio.FullAlertingFileDoc{}
	if err := configio.Load(h.path(configio.AlertingFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, r, "alerting.html", map[string]any{"Alerting": doc.Alerting})
}

func (h *Handlers) AlertingPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	doc := &configio.FullAlertingFileDoc{Alerting: configio.FullAlertingSettings{}}

	s := &doc.Alerting

	if enabled(r, "ntfy") {
		priority, _ := strconv.Atoi(formStr(r, "ntfy_priority"))
		s.Ntfy = &configio.NtfyProvider{
			ServerURL:    formStr(r, "ntfy_server_url"),
			Topic:        formStr(r, "ntfy_topic"),
			Token:        formStr(r, "ntfy_token"),
			Priority:     priority,
			Click:        formStr(r, "ntfy_click"),
			DefaultAlert: defaultAlertFrom(r, "ntfy"),
		}
	}
	if enabled(r, "slack") {
		s.Slack = &configio.SlackProvider{
			WebhookURL:   formStr(r, "slack_webhook_url"),
			Channel:      formStr(r, "slack_channel"),
			BotName:      formStr(r, "slack_bot_name"),
			DefaultAlert: defaultAlertFrom(r, "slack"),
		}
	}
	if enabled(r, "discord") {
		s.Discord = &configio.DiscordProvider{
			WebhookURL:     formStr(r, "discord_webhook_url"),
			Title:          formStr(r, "discord_title"),
			MessageContent: formStr(r, "discord_message"),
			DefaultAlert:   defaultAlertFrom(r, "discord"),
		}
	}
	if enabled(r, "telegram") {
		s.Telegram = &configio.TelegramProvider{
			Token:        formStr(r, "telegram_token"),
			ID:           formStr(r, "telegram_id"),
			DefaultAlert: defaultAlertFrom(r, "telegram"),
		}
	}
	if enabled(r, "email") {
		port, _ := strconv.Atoi(formStr(r, "email_port"))
		s.Email = &configio.EmailProvider{
			From:         formStr(r, "email_from"),
			Username:     formStr(r, "email_username"),
			Password:     formStr(r, "email_password"),
			Host:         formStr(r, "email_host"),
			Port:         port,
			To:           formStr(r, "email_to"),
			DefaultAlert: defaultAlertFrom(r, "email"),
		}
	}
	if enabled(r, "pagerduty") {
		s.PagerDuty = &configio.PagerDutyProvider{
			IntegrationKey: formStr(r, "pagerduty_integration_key"),
			Severity:       formStr(r, "pagerduty_severity"),
			DefaultAlert:   defaultAlertFrom(r, "pagerduty"),
		}
	}
	if enabled(r, "opsgenie") {
		s.Opsgenie = &configio.OpsgenieProvider{
			APIKey:       formStr(r, "opsgenie_api_key"),
			Priority:     formStr(r, "opsgenie_priority"),
			Tags:         splitCSV(formStr(r, "opsgenie_tags")),
			DefaultAlert: defaultAlertFrom(r, "opsgenie"),
		}
	}
	if enabled(r, "pushover") {
		priority, _ := strconv.Atoi(formStr(r, "pushover_priority"))
		s.Pushover = &configio.PushoverProvider{
			ApplicationToken: formStr(r, "pushover_app_token"),
			UserKey:          formStr(r, "pushover_user_key"),
			Priority:         priority,
			DefaultAlert:     defaultAlertFrom(r, "pushover"),
		}
	}
	if enabled(r, "gotify") {
		priority, _ := strconv.Atoi(formStr(r, "gotify_priority"))
		s.Gotify = &configio.GotifyProvider{
			ServerURL:    formStr(r, "gotify_server_url"),
			Token:        formStr(r, "gotify_token"),
			Priority:     priority,
			Title:        formStr(r, "gotify_title"),
			DefaultAlert: defaultAlertFrom(r, "gotify"),
		}
	}
	if enabled(r, "github") {
		s.GitHub = &configio.GitHubProvider{
			RepositoryURL: formStr(r, "github_repository_url"),
			Token:         formStr(r, "github_token"),
			DefaultAlert:  defaultAlertFrom(r, "github"),
		}
	}
	if enabled(r, "gitea") {
		s.Gitea = &configio.GiteaProvider{
			RepositoryURL: formStr(r, "gitea_repository_url"),
			Token:         formStr(r, "gitea_token"),
			DefaultAlert:  defaultAlertFrom(r, "gitea"),
		}
	}
	if enabled(r, "gitlab") {
		s.GitLab = &configio.GitLabProvider{
			WebhookURL:       formStr(r, "gitlab_webhook_url"),
			AuthorizationKey: formStr(r, "gitlab_auth_key"),
			Severity:         formStr(r, "gitlab_severity"),
			DefaultAlert:     defaultAlertFrom(r, "gitlab"),
		}
	}
	if enabled(r, "matrix") {
		s.Matrix = &configio.MatrixProvider{
			HomeserverURL: formStr(r, "matrix_homeserver_url"),
			UserID:        formStr(r, "matrix_user_id"),
			Password:      formStr(r, "matrix_password"),
			RoomID:        formStr(r, "matrix_room_id"),
			MsgType:       formStr(r, "matrix_msgtype"),
			DefaultAlert:  defaultAlertFrom(r, "matrix"),
		}
	}
	if enabled(r, "line") {
		s.Line = &configio.LineProvider{
			Token:        formStr(r, "line_token"),
			UserID:       formStr(r, "line_user_id"),
			DefaultAlert: defaultAlertFrom(r, "line"),
		}
	}
	if enabled(r, "twilio") {
		s.Twilio = &configio.TwilioProvider{
			AccountSID:   formStr(r, "twilio_account_sid"),
			AuthToken:    formStr(r, "twilio_auth_token"),
			From:         formStr(r, "twilio_from"),
			To:           formStr(r, "twilio_to"),
			DefaultAlert: defaultAlertFrom(r, "twilio"),
		}
	}
	if enabled(r, "vonage") {
		s.Vonage = &configio.VonageProvider{
			APIKey:       formStr(r, "vonage_api_key"),
			APISecret:    formStr(r, "vonage_api_secret"),
			From:         formStr(r, "vonage_from"),
			To:           formStr(r, "vonage_to"),
			DefaultAlert: defaultAlertFrom(r, "vonage"),
		}
	}
	if enabled(r, "plivo") {
		s.Plivo = &configio.PlivoProvider{
			AuthID:       formStr(r, "plivo_auth_id"),
			AuthToken:    formStr(r, "plivo_auth_token"),
			From:         formStr(r, "plivo_from"),
			To:           formStr(r, "plivo_to"),
			DefaultAlert: defaultAlertFrom(r, "plivo"),
		}
	}
	if enabled(r, "messagebird") {
		s.Messagebird = &configio.MessagebirdProvider{
			AccessKey:    formStr(r, "messagebird_access_key"),
			Originator:   formStr(r, "messagebird_originator"),
			Recipients:   formStr(r, "messagebird_recipients"),
			DefaultAlert: defaultAlertFrom(r, "messagebird"),
		}
	}
	if enabled(r, "sendgrid") {
		s.SendGrid = &configio.SendGridProvider{
			APIKey:       formStr(r, "sendgrid_api_key"),
			From:         formStr(r, "sendgrid_from"),
			To:           formStr(r, "sendgrid_to"),
			DefaultAlert: defaultAlertFrom(r, "sendgrid"),
		}
	}
	if enabled(r, "awsses") {
		s.AwsSES = &configio.AwsSESProvider{
			AccessKeyID:     formStr(r, "awsses_access_key_id"),
			SecretAccessKey: formStr(r, "awsses_secret_access_key"),
			Region:          formStr(r, "awsses_region"),
			From:            formStr(r, "awsses_from"),
			To:              formStr(r, "awsses_to"),
			DefaultAlert:    defaultAlertFrom(r, "awsses"),
		}
	}
	if enabled(r, "signal") {
		s.Signal = &configio.SignalProvider{
			Number:       formStr(r, "signal_number"),
			Recipients:   splitCSV(formStr(r, "signal_recipients")),
			URL:          formStr(r, "signal_url"),
			DefaultAlert: defaultAlertFrom(r, "signal"),
		}
	}
	if enabled(r, "datadog") {
		s.Datadog = &configio.DatadogProvider{
			APIKey:       formStr(r, "datadog_api_key"),
			Site:         formStr(r, "datadog_site"),
			Tags:         splitCSV(formStr(r, "datadog_tags")),
			DefaultAlert: defaultAlertFrom(r, "datadog"),
		}
	}
	if enabled(r, "newrelic") {
		s.NewRelic = &configio.NewRelicProvider{
			APIKey:       formStr(r, "newrelic_api_key"),
			AccountID:    formStr(r, "newrelic_account_id"),
			DefaultAlert: defaultAlertFrom(r, "newrelic"),
		}
	}
	if enabled(r, "splunk") {
		s.Splunk = &configio.SplunkProvider{
			HECURL:       formStr(r, "splunk_hec_url"),
			HECToken:     formStr(r, "splunk_hec_token"),
			Source:       formStr(r, "splunk_source"),
			SourceType:   formStr(r, "splunk_source_type"),
			DefaultAlert: defaultAlertFrom(r, "splunk"),
		}
	}
	if enabled(r, "squadcast") {
		s.Squadcast = &configio.SquadcastProvider{
			WebhookURL:   formStr(r, "squadcast_webhook_url"),
			RoutingKey:   formStr(r, "squadcast_routing_key"),
			DefaultAlert: defaultAlertFrom(r, "squadcast"),
		}
	}
	if enabled(r, "ilert") {
		s.Ilert = &configio.IlertProvider{
			APIKey:       formStr(r, "ilert_api_key"),
			IncidentKey:  formStr(r, "ilert_incident_key"),
			DefaultAlert: defaultAlertFrom(r, "ilert"),
		}
	}
	if enabled(r, "incidentio") {
		s.IncidentIO = &configio.IncidentIOProvider{
			RoutingKey:   formStr(r, "incidentio_routing_key"),
			DefaultAlert: defaultAlertFrom(r, "incidentio"),
		}
	}
	if enabled(r, "signl4") {
		s.Signl4 = &configio.Signl4Provider{
			TeamSecret:   formStr(r, "signl4_team_secret"),
			DefaultAlert: defaultAlertFrom(r, "signl4"),
		}
	}
	if enabled(r, "homeassistant") {
		s.HomeAssistant = &configio.HomeAssistantProvider{
			URL:          formStr(r, "homeassistant_url"),
			Token:        formStr(r, "homeassistant_token"),
			DefaultAlert: defaultAlertFrom(r, "homeassistant"),
		}
	}
	if enabled(r, "ifttt") {
		s.IFTTT = &configio.IFTTTProvider{
			WebhookKey:   formStr(r, "ifttt_webhook_key"),
			EventName:    formStr(r, "ifttt_event_name"),
			DefaultAlert: defaultAlertFrom(r, "ifttt"),
		}
	}
	if enabled(r, "clickup") {
		s.Clickup = &configio.ClickupProvider{
			ListID:       formStr(r, "clickup_list_id"),
			Token:        formStr(r, "clickup_token"),
			DefaultAlert: defaultAlertFrom(r, "clickup"),
		}
	}
	if enabled(r, "zulip") {
		s.Zulip = &configio.ZulipProvider{
			APIURL:       formStr(r, "zulip_api_url"),
			APIKey:       formStr(r, "zulip_api_key"),
			BotName:      formStr(r, "zulip_bot_name"),
			Organization: formStr(r, "zulip_organization"),
			To:           formStr(r, "zulip_to"),
			DefaultAlert: defaultAlertFrom(r, "zulip"),
		}
	}
	// Webhook-only providers
	for _, p := range []struct {
		key string
		dst **configio.WebhookProvider
	}{
		{"mattermost", &s.Mattermost},
		{"googlechat", &s.GoogleChat},
		{"rocketchat", &s.RocketChat},
		{"n8n", &s.N8n},
		{"teams_workflows", &s.Teams},
		{"webex", &s.Webex},
		{"zapier", &s.Zapier},
	} {
		if enabled(r, p.key) {
			*p.dst = &configio.WebhookProvider{
				WebhookURL:   formStr(r, p.key+"_webhook_url"),
				DefaultAlert: defaultAlertFrom(r, p.key),
			}
		}
	}

	if enabled(r, "custom") {
		s.Custom = &configio.CustomConfig{
			URL:          formStr(r, "custom_url"),
			Method:       formStr(r, "custom_method"),
			Body:         r.FormValue("custom_body"),
			DefaultAlert: defaultAlertFrom(r, "custom"),
		}
	}

	if err := configio.Save(h.path(configio.AlertingFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "alerting saved")
	http.Redirect(w, r, "/admin/alerting", http.StatusFound)
}

func enabled(r *http.Request, prefix string) bool {
	return r.FormValue(prefix+"_enabled") == "on"
}

func formStr(r *http.Request, name string) string {
	return strings.TrimSpace(r.FormValue(name))
}

func defaultAlertFrom(r *http.Request, prefix string) *configio.DefaultAlert {
	ft, _ := strconv.Atoi(formStr(r, prefix+"_failure_threshold"))
	st, _ := strconv.Atoi(formStr(r, prefix+"_success_threshold"))
	sor := r.FormValue(prefix+"_send_on_resolved") == "on"
	if ft == 0 && st == 0 && !sor {
		return nil
	}
	on := true
	return &configio.DefaultAlert{
		Enabled:          &on,
		FailureThreshold: ft,
		SuccessThreshold: st,
		SendOnResolved:   &sor,
	}
}
