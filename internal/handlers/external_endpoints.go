package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

const ExternalEndpointsFile = "external-endpoints.yaml"

func (h *Handlers) ExternalList(w http.ResponseWriter, r *http.Request) {
	doc := &configio.ExternalEndpointsFileDoc{}
	if err := configio.Load(h.path(ExternalEndpointsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, r, "external_endpoints.html", map[string]any{
		"Endpoints": doc.ExternalEndpoints,
	})
}

func (h *Handlers) ExternalCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	token := strings.TrimSpace(r.FormValue("token"))
	if name == "" || token == "" {
		h.flashErr(r, w, "name and token are required")
		http.Redirect(w, r, "/admin/external-endpoints", http.StatusFound)
		return
	}
	doc := &configio.ExternalEndpointsFileDoc{}
	if err := configio.Load(h.path(ExternalEndpointsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, e := range doc.ExternalEndpoints {
		if e.Name == name {
			h.flashErr(r, w, "duplicate name: "+name)
			http.Redirect(w, r, "/admin/external-endpoints", http.StatusFound)
			return
		}
	}
	ee := configio.ExternalEndpoint{
		Name:  name,
		Group: strings.TrimSpace(r.FormValue("group")),
		Token: token,
	}
	if hb := strings.TrimSpace(r.FormValue("heartbeat_interval")); hb != "" {
		ee.Heartbeat = &configio.Heartbeat{Interval: hb}
	}
	doc.ExternalEndpoints = append(doc.ExternalEndpoints, ee)
	if err := configio.Save(h.path(ExternalEndpointsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "external endpoint added: "+name)
	http.Redirect(w, r, "/admin/external-endpoints", http.StatusFound)
}

func (h *Handlers) ExternalDelete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	doc := &configio.ExternalEndpointsFileDoc{}
	if err := configio.Load(h.path(ExternalEndpointsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := doc.ExternalEndpoints[:0]
	for _, e := range doc.ExternalEndpoints {
		if e.Name != name {
			out = append(out, e)
		}
	}
	doc.ExternalEndpoints = out
	if err := configio.Save(h.path(ExternalEndpointsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "external endpoint deleted: "+name)
	http.Redirect(w, r, "/admin/external-endpoints", http.StatusFound)
}
