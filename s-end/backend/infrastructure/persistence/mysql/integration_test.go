package mysql

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/domain/hardware"
	httpx "github.com/ocmaster/backend/interfaces/http"
	"github.com/ocmaster/backend/interfaces/http/handler"
)

func setupIntegrationServer(t *testing.T) http.Handler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("db open: %v", err)
	}
	if err := AutoMigrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	hwRepo := NewHardwareRepo(db)
	scRepo := NewShareCodeRepo(db)
	hwSvc := application.NewHardwareService(hwRepo, scRepo)

	return httpx.NewRouter(httpx.RouterDeps{
		HardwareHandler:   handler.NewHardwareHandler(hwSvc),
		MerchantHandler:   &handler.MerchantHandler{},
		SuggestionHandler: &handler.SuggestionHandler{},
		ReferenceHandler:  &handler.ReferenceHandler{},
		JWTSecret:         "itest-secret",
	})
}

func TestIntegration_UploadAndRetrieve(t *testing.T) {
	server := setupIntegrationServer(t)

	// Upload
	body, _ := json.Marshal(hardware.HardwareInfo{
		CPU: &hardware.CpuInfo{Model: "i7-13700K", Cores: 16, Threads: 24},
		RAM: &hardware.RamInfo{TotalCapacity: "32 GB", StickCount: 2},
	})
	req := httptest.NewRequest("POST", "/api/v1/hardware/upload", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("upload: expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	code := resp["share_code"]
	if len(code) != 6 {
		t.Fatalf("expected 6-digit code, got %s", code)
	}

	// Retrieve
	req2 := httptest.NewRequest("GET", "/api/v1/hardware/"+code, nil)
	rec2 := httptest.NewRecorder()
	server.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("retrieve: expected 200, got %d", rec2.Code)
	}
	var hw hardware.HardwareUpload
	json.NewDecoder(rec2.Body).Decode(&hw)
	if hw.HardwareData.CPU.Model != "i7-13700K" {
		t.Errorf("expected i7-13700K, got %s", hw.HardwareData.CPU.Model)
	}

	// Delete
	req3 := httptest.NewRequest("DELETE", "/api/v1/hardware/"+code, nil)
	rec3 := httptest.NewRecorder()
	server.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec3.Code)
	}
}

func TestIntegration_HealthCheck(t *testing.T) {
	server := setupIntegrationServer(t)
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("health: expected 200, got %d", rec.Code)
	}
}
