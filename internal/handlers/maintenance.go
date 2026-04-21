package handlers

import (
	"net/http"
	"strings"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

func (h *Handlers) MaintenanceGet(w http.ResponseWriter, r *http.Request) {
	doc := &configio.MaintenanceWindowsFileDoc{}
	if err := configio.Load(h.path(configio.MaintenanceFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Migrate legacy single-window schema into windows[].
	if len(doc.Maintenance.Windows) == 0 && (doc.Maintenance.Start != "" || doc.Maintenance.End != "" || doc.Maintenance.Every != "") {
		doc.Maintenance.Windows = []configio.MWindow{{
			Start:     doc.Maintenance.Start,
			Duration:  doc.Maintenance.End,
			Recurring: doc.Maintenance.Every,
		}}
		doc.Maintenance.Start = ""
		doc.Maintenance.End = ""
		doc.Maintenance.Every = ""
	}
	h.render(w, r, "maintenance.html", map[string]any{
		"Windows": doc.Maintenance.Windows,
	})
}

func (h *Handlers) MaintenancePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	names := r.Form["w_name"]
	starts := r.Form["w_start"]
	durations := r.Form["w_duration"]
	recurring := r.Form["w_recurring"]
	endpoints := r.Form["w_endpoints"]

	var windows []configio.MWindow
	for i := 0; i < len(names); i++ {
		start := trimAt(starts, i)
		duration := trimAt(durations, i)
		rec := trimAt(recurring, i)
		name := trimAt(names, i)
		eps := splitCSV(trimAt(endpoints, i))
		if start == "" && duration == "" && rec == "" && name == "" && len(eps) == 0 {
			continue
		}
		windows = append(windows, configio.MWindow{
			Name:      name,
			Start:     start,
			Duration:  duration,
			Recurring: rec,
			Endpoints: eps,
		})
	}
	doc := &configio.MaintenanceWindowsFileDoc{Maintenance: configio.MaintenanceBlock{Windows: windows}}
	if err := configio.Save(h.path(configio.MaintenanceFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "maintenance saved")
	http.Redirect(w, r, "/admin/maintenance", http.StatusFound)
}

func trimAt(vs []string, i int) string {
	if i >= len(vs) {
		return ""
	}
	return strings.TrimSpace(vs[i])
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
