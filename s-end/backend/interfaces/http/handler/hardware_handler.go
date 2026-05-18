package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/hardware"
)

type HardwareHandler struct {
	svc *application.HardwareService
}

func NewHardwareHandler(svc *application.HardwareService) *HardwareHandler {
	return &HardwareHandler{svc: svc}
}

func (h *HardwareHandler) Upload(w http.ResponseWriter, r *http.Request) {
	var info hardware.HardwareInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	code, err := h.svc.Upload(r.Context(), &info)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "upload failed"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"share_code": code})
}

func (h *HardwareHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	upload, err := h.svc.GetByCode(r.Context(), code)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, upload)
}

func (h *HardwareHandler) DeleteByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if err := h.svc.DeleteByCode(r.Context(), code); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "delete failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
