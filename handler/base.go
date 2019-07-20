package handler

import "net/http"

// New returns a mux with registered patterns.
func New() http.Handler {
	// A mux intelligently matches the URL of incoming reqs against registered patterns
	mux := http.NewServeMux()

	// Root
	mux.Handle("/", http.FileServer(http.Dir("tmpl/home.html")))

	// OauthGoogle
	mux.HandleFunc("/auth/google/login", googleLogin)
	mux.HandleFunc("/auth/google/callback", googleCallback)

	return mux
}
