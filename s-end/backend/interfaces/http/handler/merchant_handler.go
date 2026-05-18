package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/infrastructure/auth"
	"github.com/ocmaster/backend/interfaces/http/middleware"
)

type MerchantHandler struct {
	svc *application.MerchantService
}

func NewMerchantHandler(svc *application.MerchantService) *MerchantHandler {
	return &MerchantHandler{svc: svc}
}

type registerReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (h *MerchantHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	if err := h.svc.Register(r.Context(), req.Phone, req.Password); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "register failed"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "registered, pending review"})
}

type loginReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (h *MerchantHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	token, err := h.svc.Login(r.Context(), req.Phone, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *MerchantHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	m, err := h.svc.GetProfile(r.Context(), claims.MerchantID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (h *MerchantHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	m, err := h.svc.GetProfile(r.Context(), claims.MerchantID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	var req struct {
		Phone        string `json:"phone"`
		RiskTemplate string `json:"risk_template"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	if req.Phone != "" {
		m.Phone = req.Phone
	}
	if req.RiskTemplate != "" {
		m.RiskTemplate = req.RiskTemplate
	}
	if err := h.svc.UpdateProfile(r.Context(), m); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *MerchantHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	if err := h.svc.ChangePassword(r.Context(), claims.MerchantID, req.OldPassword, req.NewPassword); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "password change failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "password changed"})
}
