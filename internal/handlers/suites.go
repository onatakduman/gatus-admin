package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

const SuitesFile = "suites.yaml"

func (h *Handlers) SuitesList(w http.ResponseWriter, r *http.Request) {
	doc := &configio.SuitesFileDoc{}
	if err := configio.Load(h.path(SuitesFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, r, "suites.html", map[string]any{"Suites": doc.Suites})
}

func (h *Handlers) SuiteCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		h.flashErr(r, w, "name is required")
		http.Redirect(w, r, "/admin/suites", http.StatusFound)
		return
	}
	doc := &configio.SuitesFileDoc{}
	if err := configio.Load(h.path(SuitesFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, s := range doc.Suites {
		if s.Name == name {
			h.flashErr(r, w, "duplicate name: "+name)
			http.Redirect(w, r, "/admin/suites", http.StatusFound)
			return
		}
	}
	s := configio.Suite{
		Name:     name,
		Group:    strings.TrimSpace(r.FormValue("group")),
		Interval: strings.TrimSpace(r.FormValue("interval")),
		Timeout:  strings.TrimSpace(r.FormValue("timeout")),
	}
	// Context as multi-line key=value input
	if ctxRaw := strings.TrimSpace(r.FormValue("context")); ctxRaw != "" {
		ctx := map[string]any{}
		for _, line := range strings.Split(ctxRaw, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			ctx[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		if len(ctx) > 0 {
			s.Context = ctx
		}
	}
	doc.Suites = append(doc.Suites, s)
	if err := configio.Save(h.path(SuitesFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "suite added: "+name+" (add endpoints via Advanced YAML)")
	http.Redirect(w, r, "/admin/suites", http.StatusFound)
}

func (h *Handlers) SuiteDelete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	doc := &configio.SuitesFileDoc{}
	if err := configio.Load(h.path(SuitesFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := doc.Suites[:0]
	for _, s := range doc.Suites {
		if s.Name != name {
			out = append(out, s)
		}
	}
	doc.Suites = out
	if err := configio.Save(h.path(SuitesFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "suite deleted: "+name)
	http.Redirect(w, r, "/admin/suites", http.StatusFound)
}
