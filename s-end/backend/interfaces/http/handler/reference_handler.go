package handler

import (
	"net/http"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/reference"
)

type ReferenceHandler struct {
	svc *application.ReferenceService
}

func NewReferenceHandler(svc *application.ReferenceService) *ReferenceHandler {
	return &ReferenceHandler{svc: svc}
}

func (h *ReferenceHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	data, err := h.svc.GetByCategory(r.Context(), reference.Category(category))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "query failed"})
		return
	}
	if data == nil {
		data = []*reference.ReferenceData{}
	}
	writeJSON(w, http.StatusOK, data)
}
