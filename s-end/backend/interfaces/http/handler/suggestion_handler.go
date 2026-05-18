package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/suggestion"
	"github.com/ocmaster/backend/infrastructure/auth"
	"github.com/ocmaster/backend/interfaces/http/middleware"
)

type SuggestionHandler struct {
	svc *application.SuggestionService
}

func NewSuggestionHandler(svc *application.SuggestionService) *SuggestionHandler {
	return &SuggestionHandler{svc: svc}
}

func (h *SuggestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	var sug suggestion.Suggestion
	if err := json.NewDecoder(r.Body).Decode(&sug); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	sug.MerchantID = claims.MerchantID
	if err := h.svc.Create(r.Context(), &sug); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": sug.ID})
}

func (h *SuggestionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	sug, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, sug)
}

func (h *SuggestionHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	sugs, total, err := h.svc.ListHistory(r.Context(), claims.MerchantID, offset, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "query failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"items": sugs, "total": total})
}

func (h *SuggestionHandler) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	pdfData, err := h.svc.GeneratePDF(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "pdf generation failed"})
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
	w.Write(pdfData)
}
