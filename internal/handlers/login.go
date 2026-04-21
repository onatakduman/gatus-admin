package handlers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) LoginGet(w http.ResponseWriter, r *http.Request) {
	if h.Store.IsAuthed(r) {
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}
	h.render(w, r, "login.html", nil)
}

func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")
	if err := bcrypt.CompareHashAndPassword([]byte(h.PasswordHash), []byte(password)); err != nil {
		h.flashErr(r, w, "invalid password")
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}
	if err := h.Store.SetAuthed(r, w, true); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	_ = h.Store.SetAuthed(r, w, false)
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}
