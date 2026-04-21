package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

func (h *Handlers) AdvancedGet(w http.ResponseWriter, r *http.Request) {
	type fileInfo struct {
		Name    string
		Size    int
		Content string
	}
	var files []fileInfo
	for _, f := range configio.AllFiles {
		content, err := configio.LoadRaw(h.path(f))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		files = append(files, fileInfo{Name: f, Size: len(content), Content: content})
	}
	h.render(w, r, "advanced.html", map[string]any{"Files": files})
}

func (h *Handlers) AdvancedFileGet(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "file")
	if !isAllowedFile(name) {
		http.NotFound(w, r)
		return
	}
	content, err := configio.LoadRaw(h.path(name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, r, "advanced_edit.html", map[string]any{
		"FileName": name,
		"Content":  content,
	})
}

func (h *Handlers) AdvancedFilePost(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "file")
	if !isAllowedFile(name) {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	content := r.FormValue("content")
	if err := configio.SaveRaw(h.path(name), content); err != nil {
		h.flashErr(r, w, err.Error())
		http.Redirect(w, r, "/admin/advanced", http.StatusFound)
		return
	}
	h.flashOk(r, w, name+" saved")
	http.Redirect(w, r, "/admin/advanced", http.StatusFound)
}

func isAllowedFile(name string) bool {
	for _, f := range configio.AllFiles {
		if f == name {
			return true
		}
	}
	return false
}
