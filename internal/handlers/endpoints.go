package handlers

import (
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v3"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

func (h *Handlers) loadEndpoints() (*configio.EndpointsFileDoc, error) {
	doc := &configio.EndpointsFileDoc{}
	if err := configio.Load(h.path(configio.EndpointsFile), doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (h *Handlers) saveEndpoints(doc *configio.EndpointsFileDoc) error {
	return configio.Save(h.path(configio.EndpointsFile), doc)
}

func (h *Handlers) EndpointsList(w http.ResponseWriter, r *http.Request) {
	doc, err := h.loadEndpoints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	groups := map[string]bool{}
	for _, e := range doc.Endpoints {
		if e.Group != "" {
			groups[e.Group] = true
		}
	}
	groupList := make([]string, 0, len(groups))
	for g := range groups {
		groupList = append(groupList, g)
	}
	sort.Strings(groupList)
	h.render(w, r, "endpoints.html", map[string]any{
		"Endpoints": doc.Endpoints,
		"Groups":    groupList,
	})
}

func parseSharedEndpointFields(r *http.Request) configio.Endpoint {
	e := configio.Endpoint{
		Group:    strings.TrimSpace(r.FormValue("group")),
		Method:   strings.TrimSpace(r.FormValue("method")),
		Body:     r.FormValue("body"),
		Interval: strings.TrimSpace(r.FormValue("interval")),
	}
	if r.FormValue("graphql") == "on" {
		t := true
		e.GraphQL = &t
	}
	for _, c := range strings.Split(r.FormValue("conditions"), "\n") {
		c = strings.TrimSpace(c)
		if c != "" {
			e.Conditions = append(e.Conditions, c)
		}
	}
	headers := map[string]string{}
	for _, line := range strings.Split(r.FormValue("headers"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	if len(headers) > 0 {
		e.Headers = headers
	}
	labels := map[string]string{}
	for _, line := range strings.Split(r.FormValue("extra_labels"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	if len(labels) > 0 {
		e.ExtraLabels = labels
	}
	alertType := strings.TrimSpace(r.FormValue("alert_type"))
	alertDesc := strings.TrimSpace(r.FormValue("alert_description"))
	if alertType != "" {
		e.Alerts = []configio.Alert{{Type: alertType, Description: alertDesc}}
	}

	// Client options
	timeout := strings.TrimSpace(r.FormValue("client_timeout"))
	proxy := strings.TrimSpace(r.FormValue("client_proxy_url"))
	resolver := strings.TrimSpace(r.FormValue("client_dns_resolver"))
	tunnel := strings.TrimSpace(r.FormValue("client_tunnel"))
	insecure := r.FormValue("client_insecure") == "on"
	ignoreRedirect := r.FormValue("client_ignore_redirect") == "on"
	if timeout != "" || proxy != "" || resolver != "" || tunnel != "" || insecure || ignoreRedirect {
		co := &configio.ClientOpts{Timeout: timeout, ProxyURL: proxy, DNSResolver: resolver, Tunnel: tunnel}
		if insecure {
			co.Insecure = &insecure
		}
		if ignoreRedirect {
			co.IgnoreRedirect = &ignoreRedirect
		}
		e.Client = co
	}

	// DNS
	dnsType := strings.TrimSpace(r.FormValue("dns_query_type"))
	dnsName := strings.TrimSpace(r.FormValue("dns_query_name"))
	if dnsType != "" || dnsName != "" {
		e.DNS = &configio.EndpointDNS{QueryType: dnsType, QueryName: dnsName}
	}

	// SSH
	sshUser := strings.TrimSpace(r.FormValue("ssh_username"))
	sshPass := strings.TrimSpace(r.FormValue("ssh_password"))
	if sshUser != "" || sshPass != "" {
		e.SSH = &configio.EndpointSSH{Username: sshUser, Password: sshPass}
	}

	// Endpoint UI hide-*
	hideHost := r.FormValue("ui_hide_hostname") == "on"
	hideURL := r.FormValue("ui_hide_url") == "on"
	hideErr := r.FormValue("ui_hide_errors") == "on"
	if hideHost || hideURL || hideErr {
		e.UI = &configio.EndpointUI{HideHostname: hideHost, HideURL: hideURL, HideErrors: hideErr}
	}

	return e
}

func parseEndpointForm(r *http.Request) (configio.Endpoint, error) {
	if err := r.ParseForm(); err != nil {
		return configio.Endpoint{}, err
	}
	e := parseSharedEndpointFields(r)
	e.Name = strings.TrimSpace(r.FormValue("name"))
	e.URL = strings.TrimSpace(r.FormValue("url"))
	return e, nil
}

func deriveNameFromURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	if u, err := url.Parse(rawURL); err == nil && u.Host != "" {
		return u.Hostname()
	}
	return rawURL
}

func (h *Handlers) EndpointCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shared := parseSharedEndpointFields(r)
	urls := r.Form["url"]
	names := r.Form["name"]

	doc, err := h.loadEndpoints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	existing := map[string]bool{}
	for _, ep := range doc.Endpoints {
		existing[ep.Name] = true
	}

	added, skipped, empty := 0, 0, 0
	var firstErr string
	for i, rawURL := range urls {
		u := strings.TrimSpace(rawURL)
		if u == "" {
			empty++
			continue
		}
		var name string
		if i < len(names) {
			name = strings.TrimSpace(names[i])
		}
		if name == "" {
			name = deriveNameFromURL(u)
		}
		if name == "" {
			if firstErr == "" {
				firstErr = "could not derive name for URL: " + u
			}
			skipped++
			continue
		}
		if existing[name] {
			skipped++
			if firstErr == "" {
				firstErr = "duplicate name skipped: " + name
			}
			continue
		}
		ep := shared
		ep.Name = name
		ep.URL = u
		doc.Endpoints = append(doc.Endpoints, ep)
		existing[name] = true
		added++
	}

	if added == 0 {
		h.flashErr(r, w, "no endpoints added"+ifStr(firstErr != "", " ("+firstErr+")"))
		http.Redirect(w, r, "/admin/endpoints", http.StatusFound)
		return
	}
	if err := h.saveEndpoints(doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg := "added " + itoa(added)
	if skipped > 0 {
		msg += ", skipped " + itoa(skipped)
	}
	if firstErr != "" {
		msg += " (" + firstErr + ")"
	}
	h.flashOk(r, w, msg)
	http.Redirect(w, r, "/admin/endpoints", http.StatusFound)
}

func ifStr(cond bool, s string) string {
	if cond {
		return s
	}
	return ""
}

func (h *Handlers) EndpointEdit(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	doc, err := h.loadEndpoints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var found *configio.Endpoint
	for i := range doc.Endpoints {
		if doc.Endpoints[i].Name == name {
			found = &doc.Endpoints[i]
			break
		}
	}
	if found == nil {
		http.NotFound(w, r)
		return
	}
	headerLines := []string{}
	for k, v := range found.Headers {
		headerLines = append(headerLines, k+": "+v)
	}
	h.render(w, r, "endpoint_edit.html", map[string]any{
		"Endpoint":        found,
		"ConditionsJoined": strings.Join(found.Conditions, "\n"),
		"HeadersJoined":   strings.Join(headerLines, "\n"),
	})
}

func (h *Handlers) EndpointUpdate(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	newE, err := parseEndpointForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	doc, err := h.loadEndpoints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idx := -1
	for i, e := range doc.Endpoints {
		if e.Name == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		http.NotFound(w, r)
		return
	}
	doc.Endpoints[idx] = newE
	if err := h.saveEndpoints(doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "endpoint updated: "+newE.Name)
	http.Redirect(w, r, "/admin/endpoints", http.StatusFound)
}

func (h *Handlers) EndpointDelete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	doc, err := h.loadEndpoints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := doc.Endpoints[:0]
	for _, e := range doc.Endpoints {
		if e.Name != name {
			out = append(out, e)
		}
	}
	doc.Endpoints = out
	if err := h.saveEndpoints(doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "endpoint deleted: "+name)
	http.Redirect(w, r, "/admin/endpoints", http.StatusFound)
}

func (h *Handlers) EndpointsBulkGet(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "endpoints_bulk.html", nil)
}

func (h *Handlers) EndpointsBulkPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	content := r.FormValue("yaml")
	content = strings.TrimSpace(content)
	if content == "" {
		h.flashErr(r, w, "empty input")
		http.Redirect(w, r, "/admin/endpoints/bulk", http.StatusFound)
		return
	}

	var parsed []configio.Endpoint
	if err := yaml.Unmarshal([]byte(content), &parsed); err != nil {
		var wrap struct {
			Endpoints []configio.Endpoint `yaml:"endpoints"`
		}
		if err2 := yaml.Unmarshal([]byte(content), &wrap); err2 != nil {
			h.flashErr(r, w, "parse error: "+err.Error())
			http.Redirect(w, r, "/admin/endpoints/bulk", http.StatusFound)
			return
		}
		parsed = wrap.Endpoints
	}
	if len(parsed) == 0 {
		h.flashErr(r, w, "no endpoints parsed")
		http.Redirect(w, r, "/admin/endpoints/bulk", http.StatusFound)
		return
	}

	doc, err := h.loadEndpoints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	existing := map[string]bool{}
	for _, e := range doc.Endpoints {
		existing[e.Name] = true
	}
	added, skipped := 0, 0
	for _, e := range parsed {
		if e.Name == "" || e.URL == "" {
			skipped++
			continue
		}
		if existing[e.Name] {
			skipped++
			continue
		}
		doc.Endpoints = append(doc.Endpoints, e)
		existing[e.Name] = true
		added++
	}
	if err := h.saveEndpoints(doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "imported "+itoa(added)+", skipped "+itoa(skipped))
	http.Redirect(w, r, "/admin/endpoints", http.StatusFound)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 8)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}
