package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const csrfCookie = "csrf_token"

func getOrSetCSRF(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie(csrfCookie); err == nil && c.Value != "" {
		return c.Value
	}
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)
	http.SetCookie(w, &http.Cookie{
		Name: csrfCookie, Value: token, Path: "/",
		HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	return token
}

func validateCSRF(r *http.Request) bool {
	c, err := r.Cookie(csrfCookie)
	if err != nil {
		return false
	}
	return r.FormValue("csrf_token") != "" && r.FormValue("csrf_token") == c.Value
}
