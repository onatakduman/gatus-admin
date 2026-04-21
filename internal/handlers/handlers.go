package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/onatakduman/gatus-admin/internal/auth"
	"github.com/onatakduman/gatus-admin/internal/brand"
	"github.com/onatakduman/gatus-admin/internal/configio"
	"github.com/onatakduman/gatus-admin/internal/web"
)

type Handlers struct {
	ConfigDir      string
	GatusContainer string
	Tpl            *web.Templates
	Store          *auth.Store
	PasswordHash   string
	Dev            bool
	StatusURL      string
	BrandStore     brand.Store
}

func (h *Handlers) path(name string) string {
	return filepath.Join(h.ConfigDir, name)
}

func (h *Handlers) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.Store.IsAuthed(r) {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) render(w http.ResponseWriter, r *http.Request, name string, data map[string]any) {
	if data == nil {
		data = map[string]any{}
	}
	data["FlashError"] = h.Store.PopFlashes(r, w, auth.FlashErrorKey)
	data["FlashOk"] = h.Store.PopFlashes(r, w, auth.FlashOkKey)
	data["StatusURL"] = h.StatusURL
	if _, exists := data["Brand"]; !exists {
		b, _ := h.BrandStore.Load()
		data["Brand"] = h.computeBrand(b)
	}
	// Admin-only custom CSS — always applied.
	if effective, ok := data["Brand"].(brand.Effective); ok {
		data["BrandCustomCSS"] = template.CSS(effective.CustomCSS)
		// Mirror Gatus public-page custom-css ONLY when the operator
		// explicitly enabled the toggle in the Brand panel.
		if effective.MirrorPublicCSS {
			uiDoc := &configio.UIFileDoc{}
			_ = configio.Load(h.path(configio.UIFile), uiDoc)
			data["GatusCustomCSS"] = template.CSS(uiDoc.UI.CustomCSS)
		} else {
			data["GatusCustomCSS"] = template.CSS("")
		}
	}
	if err := h.Tpl.Render(w, r, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) flashErr(r *http.Request, w http.ResponseWriter, msg string) {
	_ = h.Store.AddFlash(r, w, auth.FlashErrorKey, msg)
}

func (h *Handlers) flashOk(r *http.Request, w http.ResponseWriter, msg string) {
	_ = h.Store.AddFlash(r, w, auth.FlashOkKey, msg)
}
