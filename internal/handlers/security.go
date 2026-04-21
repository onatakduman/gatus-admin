package handlers

import (
	"net/http"
	"strings"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

const SecurityFile = "security.yaml"

func (h *Handlers) SecurityGet(w http.ResponseWriter, r *http.Request) {
	doc := &configio.SecurityFileDoc{}
	if err := configio.Load(h.path(SecurityFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, r, "security.html", map[string]any{"Security": doc.Security})
}

func (h *Handlers) SecurityPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	authType := strings.TrimSpace(r.FormValue("auth_type"))
	doc := &configio.SecurityFileDoc{}

	if authType == "none" || authType == "" {
		doc.Security = nil
	} else {
		sec := &configio.FullSecurity{Authentication: &configio.Authentication{Type: authType}}
		switch authType {
		case "basic":
			user := strings.TrimSpace(r.FormValue("basic_username"))
			pass := strings.TrimSpace(r.FormValue("basic_password"))
			if user != "" || pass != "" {
				sec.Authentication.Basic = &configio.BasicAuth{Username: user, Password: pass}
			}
		case "oidc":
			scopes := splitCSV(r.FormValue("oidc_scopes"))
			sec.Authentication.OIDC = &configio.OIDC{
				IssuerURL:    strings.TrimSpace(r.FormValue("oidc_issuer_url")),
				RedirectURL:  strings.TrimSpace(r.FormValue("oidc_redirect_url")),
				ClientID:     strings.TrimSpace(r.FormValue("oidc_client_id")),
				ClientSecret: strings.TrimSpace(r.FormValue("oidc_client_secret")),
				Scopes:       scopes,
			}
		}
		defRole := strings.TrimSpace(r.FormValue("default_role"))
		admin := splitCSV(r.FormValue("admin_groups"))
		viewer := splitCSV(r.FormValue("viewer_groups"))
		if defRole != "" || len(admin) > 0 || len(viewer) > 0 {
			sec.Authorization = &configio.Authorization{
				DefaultRole:  defRole,
				AdminGroups:  admin,
				ViewerGroups: viewer,
			}
		}
		doc.Security = sec
	}

	if err := configio.Save(h.path(SecurityFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "security saved")
	http.Redirect(w, r, "/admin/security", http.StatusFound)
}
