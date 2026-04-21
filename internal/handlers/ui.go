package handlers

import (
	"net/http"
	"strings"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

func (h *Handlers) UIGet(w http.ResponseWriter, r *http.Request) {
	doc := &configio.UIFileDoc{}
	if err := configio.Load(h.path(configio.UIFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, r, "ui.html", map[string]any{"UI": doc.UI})
}

func (h *Handlers) UIPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	doc := &configio.UIFileDoc{}
	if err := configio.Load(h.path(configio.UIFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dark := r.FormValue("dark_mode") == "on"
	doc.UI = configio.UISettings{
		Title:               strings.TrimSpace(r.FormValue("title")),
		Description:         strings.TrimSpace(r.FormValue("description")),
		Header:              strings.TrimSpace(r.FormValue("header")),
		Logo:                strings.TrimSpace(r.FormValue("logo")),
		Link:                strings.TrimSpace(r.FormValue("link")),
		DashboardHeading:    strings.TrimSpace(r.FormValue("dashboard_heading")),
		DashboardSubheading: strings.TrimSpace(r.FormValue("dashboard_subheading")),
		DarkMode:            &dark,
		DefaultSortBy:       strings.TrimSpace(r.FormValue("default_sort_by")),
		DefaultFilterBy:     strings.TrimSpace(r.FormValue("default_filter_by")),
		CustomCSS:           r.FormValue("custom_css"),
	}
	faviconDefault := strings.TrimSpace(r.FormValue("favicon_default"))
	favicon16 := strings.TrimSpace(r.FormValue("favicon_16"))
	favicon32 := strings.TrimSpace(r.FormValue("favicon_32"))
	if faviconDefault != "" || favicon16 != "" || favicon32 != "" {
		doc.UI.Favicon = &configio.Favicon{
			Default:   faviconDefault,
			Size16x16: favicon16,
			Size32x32: favicon32,
		}
	}
	buttonNames := r.Form["button_name"]
	buttonLinks := r.Form["button_link"]
	var buttons []configio.UIButton
	for i := range buttonNames {
		name := strings.TrimSpace(buttonNames[i])
		var link string
		if i < len(buttonLinks) {
			link = strings.TrimSpace(buttonLinks[i])
		}
		if name == "" && link == "" {
			continue
		}
		buttons = append(buttons, configio.UIButton{Name: name, Link: link})
	}
	doc.UI.Buttons = buttons

	if err := configio.Save(h.path(configio.UIFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "UI saved")
	http.Redirect(w, r, "/admin/ui", http.StatusFound)
}
