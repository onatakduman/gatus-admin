package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

func (h *Handlers) CoreGet(w http.ResponseWriter, r *http.Request) {
	doc := &configio.CoreFileDoc{}
	if err := configio.Load(h.path(configio.CoreFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, r, "core.html", map[string]any{"Core": doc})
}

func (h *Handlers) CorePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	doc := &configio.CoreFileDoc{}

	webPort, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("web_port")))
	webAddr := strings.TrimSpace(r.FormValue("web_address"))
	if webPort != 0 || webAddr != "" {
		doc.Web = &configio.WebSettings{Address: webAddr, Port: webPort}
	}

	storageType := strings.TrimSpace(r.FormValue("storage_type"))
	storagePath := strings.TrimSpace(r.FormValue("storage_path"))
	if storageType != "" || storagePath != "" {
		doc.Storage = &configio.StorageSettings{Type: storageType, Path: storagePath}
	}

	user := strings.TrimSpace(r.FormValue("security_user"))
	pass := strings.TrimSpace(r.FormValue("security_pass"))
	if user != "" && pass != "" {
		doc.Security = &configio.Security{Basic: &configio.BasicAuth{Username: user, Password: pass}}
	}

	metrics := r.FormValue("metrics") == "on"
	doc.Metrics = &metrics

	conc, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("concurrency")))
	doc.Concurrency = conc

	if err := configio.Save(h.path(configio.CoreFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "core saved")
	http.Redirect(w, r, "/admin/core", http.StatusFound)
}
