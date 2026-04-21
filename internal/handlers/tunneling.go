package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/onatakduman/gatus-admin/internal/configio"
)

const TunnelingFile = "tunneling.yaml"

func (h *Handlers) TunnelingGet(w http.ResponseWriter, r *http.Request) {
	doc := &configio.TunnelingFileDoc{}
	if err := configio.Load(h.path(TunnelingFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Flatten into a slice for ordered rendering.
	type row struct {
		Key    string
		Tunnel configio.Tunnel
	}
	var rows []row
	for k, v := range doc.Tunneling {
		rows = append(rows, row{Key: k, Tunnel: v})
	}
	h.render(w, r, "tunneling.html", map[string]any{"Rows": rows})
}

func (h *Handlers) TunnelingPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	keys := r.Form["t_key"]
	hosts := r.Form["t_host"]
	ports := r.Form["t_port"]
	users := r.Form["t_user"]
	passwords := r.Form["t_password"]
	pkeys := r.Form["t_pkey"]

	out := map[string]configio.Tunnel{}
	for i := 0; i < len(keys); i++ {
		k := trimAt(keys, i)
		if k == "" {
			continue
		}
		p, _ := strconv.Atoi(trimAt(ports, i))
		out[k] = configio.Tunnel{
			Type:       "SSH",
			Host:       trimAt(hosts, i),
			Port:       p,
			Username:   trimAt(users, i),
			Password:   trimAt(passwords, i),
			PrivateKey: strings.TrimSpace(pkeyAt(pkeys, i)),
		}
	}
	doc := &configio.TunnelingFileDoc{Tunneling: out}
	if err := configio.Save(h.path(TunnelingFile), doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "tunneling saved")
	http.Redirect(w, r, "/admin/tunneling", http.StatusFound)
}

func pkeyAt(vs []string, i int) string {
	if i >= len(vs) {
		return ""
	}
	return vs[i]
}
