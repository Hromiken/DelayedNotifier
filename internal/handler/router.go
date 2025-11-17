package handler

import (
	"net/http"
)

// NewRouter - инициализация роутера.
func NewRouter(h *NotifyHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/notify", h.CreateNotify)
	mux.HandleFunc("/notify/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetNotify(w, r)
		case http.MethodDelete:
			h.DeleteNotify(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
