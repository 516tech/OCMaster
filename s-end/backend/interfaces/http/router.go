package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/ocmaster/backend/application"
	"github.com/ocmaster/backend/interfaces/http/handler"
	"github.com/ocmaster/backend/interfaces/http/middleware"
)

type RouterDeps struct {
	HardwareHandler   *handler.HardwareHandler
	MerchantHandler   *handler.MerchantHandler
	SuggestionHandler *handler.SuggestionHandler
	ReferenceHandler  *handler.ReferenceHandler
	CleanupService    *application.CleanupService
	JWTSecret         string
}

func NewRouter(deps RouterDeps) chi.Router {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.LoggerMiddleware)
	r.Use(chimw.Recoverer)

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// 公开路由
		r.Post("/hardware/upload", deps.HardwareHandler.Upload)
		r.Get("/hardware/{code}", deps.HardwareHandler.GetByCode)
		r.Delete("/hardware/{code}", deps.HardwareHandler.DeleteByCode)
		r.Post("/merchant/register", deps.MerchantHandler.Register)
		r.Post("/merchant/login", deps.MerchantHandler.Login)

		// 认证路由
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(deps.JWTSecret))
			r.Get("/merchant/profile", deps.MerchantHandler.GetProfile)
			r.Put("/merchant/profile", deps.MerchantHandler.UpdateProfile)
			r.Put("/merchant/password", deps.MerchantHandler.ChangePassword)
			r.Post("/suggestions", deps.SuggestionHandler.Create)
			r.Get("/suggestions/history", deps.SuggestionHandler.ListHistory)
			r.Get("/suggestions/{id}", deps.SuggestionHandler.GetByID)
			r.Get("/suggestions/{id}/pdf", deps.SuggestionHandler.DownloadPDF)
			r.Get("/reference/cpu", deps.ReferenceHandler.GetByCategory)
			r.Get("/reference/ram", deps.ReferenceHandler.GetByCategory)
			r.Get("/reference/cooler", deps.ReferenceHandler.GetByCategory)
		})
	})

	return r
}
