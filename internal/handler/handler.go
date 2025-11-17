package handler

import (
	"context"
	"delayedNotifier/internal/entity"
	"encoding/json"
	"net/http"
	"strings"

	"delayedNotifier/internal/service"
)

// NotifyHandler содержит сервис для работы с уведомлениями
type NotifyHandler struct {
	Service *service.NotifierService
}

// NewNotifyHandler создаёт NotifyHandler
func NewNotifyHandler(svc *service.NotifierService) *NotifyHandler {
	return &NotifyHandler{Service: svc}
}

// CreateNotify — POST /notify.
func (h *NotifyHandler) CreateNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req entity.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Sender == "" {
		req.Sender = "mock"
	}

	id, err := h.Service.CreateNotify(r.Context(), req)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": id}, http.StatusCreated)
}

// GetNotify — GET /notify/{id}.
func (h *NotifyHandler) GetNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/notify/")
	if id == "" {
		writeError(w, "id required", http.StatusBadRequest)
		return
	}

	n, err := h.Service.GetNotify(context.Background(), id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, n, http.StatusOK)
}

// DeleteNotify - DELETE /notify/{id}.
func (h *NotifyHandler) DeleteNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/notify/")
	if id == "" {
		writeError(w, "id required", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteNotify(context.Background(), id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	writeJSON(w, map[string]string{"error": msg}, code)
}
