package main

import "net/http"

func route(h *handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/b/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.getBallot(w, r)
		case http.MethodPost:
			h.postBallot(w, r)
		default:
			serve405(w, r)
		}
	})

	mux.HandleFunc("/e/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.getElection(w, r)
		case http.MethodPost:
			h.postElection(w, r)
		default:
			serve405(w, r)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.getIndex(w, r)
		case http.MethodPost:
			h.postIndex(w, r)
		default:
			serve405(w, r)
		}
	})

	return mux
}

func serve404(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func serve405(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
