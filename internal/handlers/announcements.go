package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

func (h *Handlers) AnnouncementsList(w http.ResponseWriter, r *http.Request) {
	doc := &configio.AnnouncementsFileDoc{}
	if err := configio.Load(h.path(configio.AnnouncementsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, r, "announcements.html", map[string]any{
		"Announcements": doc.Announcements,
	})
}

func (h *Handlers) AnnouncementCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg := strings.TrimSpace(r.FormValue("message"))
	if msg == "" {
		h.flashErr(r, w, "message required")
		http.Redirect(w, r, "/admin/announcements", http.StatusFound)
		return
	}
	ts := strings.TrimSpace(r.FormValue("timestamp"))
	if ts == "" {
		ts = time.Now().UTC().Format(time.RFC3339)
	}
	doc := &configio.AnnouncementsFileDoc{}
	if err := configio.Load(h.path(configio.AnnouncementsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	doc.Announcements = append(doc.Announcements, configio.Announcement{
		Timestamp: ts,
		Type:      strings.TrimSpace(r.FormValue("type")),
		Message:   msg,
		Archived:  r.FormValue("archived") == "on",
	})
	if err := configio.Save(h.path(configio.AnnouncementsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "announcement added")
	http.Redirect(w, r, "/admin/announcements", http.StatusFound)
}

func (h *Handlers) AnnouncementDelete(w http.ResponseWriter, r *http.Request) {
	idx, err := strconv.Atoi(chi.URLParam(r, "idx"))
	if err != nil {
		http.Error(w, "bad idx", http.StatusBadRequest)
		return
	}
	doc := &configio.AnnouncementsFileDoc{}
	if err := configio.Load(h.path(configio.AnnouncementsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if idx < 0 || idx >= len(doc.Announcements) {
		http.NotFound(w, r)
		return
	}
	doc.Announcements = append(doc.Announcements[:idx], doc.Announcements[idx+1:]...)
	if err := configio.Save(h.path(configio.AnnouncementsFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "announcement deleted")
	http.Redirect(w, r, "/admin/announcements", http.StatusFound)
}
