package v1

import (
	"database/sql"
	"net/http"

	"hakathon-mvp/internal/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func NewRouter(
	citizenReportUC *usecases.CitizenReportUseCase,
	db *sql.DB,
	redisClient *redis.Client,
	rateLimit int,
) http.Handler {
	r := chi.NewRouter()

	// Глобальные middleware
	r.Use(middleware.RequestID)
	r.Use(RecoverMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(ContentTypeJSONMiddleware)
	r.Use(RateLimitMiddleware(rateLimit))

	// Health checks
	healthHandler := NewHealthHandler(db, redisClient)
	healthRouter := chi.NewRouter()
	healthHandler.RegisterRoutes(healthRouter)
	r.Mount("/health", healthRouter)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		cReportHandler := NewCitizenReportHandler(citizenReportUC)

		r.Route("/citizen-reports", func(r chi.Router) {
			r.Post("/", cReportHandler.CreateCitizenReport)
			r.Get("/", cReportHandler.ListCitizenReports)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", cReportHandler.GetCitizenReport)
				r.Put("/", cReportHandler.UpdateCitizenReport)
				r.Delete("/", cReportHandler.DeleteCitizenReport)
			})
		})
	})

	return r
}
