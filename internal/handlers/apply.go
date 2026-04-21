package handlers

import (
	"net/http"

	"github.com/onatakduman/gatus-admin/internal/gatusctl"
)

func (h *Handlers) Apply(w http.ResponseWriter, r *http.Request) {
	if err := gatusctl.Restart(h.GatusContainer); err != nil {
		h.flashErr(r, w, err.Error())
	} else {
		h.flashOk(r, w, "gatus restarted")
	}
	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/admin"
	}
	http.Redirect(w, r, ref, http.StatusFound)
}
