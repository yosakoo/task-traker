package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/yosakoo/task-traker/internal/delivery/http/v1"
	"github.com/yosakoo/task-traker/internal/service"
	"github.com/yosakoo/task-traker/pkg/auth"
	"github.com/yosakoo/task-traker/pkg/logger"
	"net/http"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(l logger.Interface) *chi.Mux {
	router := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})
	router.Use(c.Handler)
	router.Use(v1.NewMwLogger(l))
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	h.initAPI(router)
	return router
}

func (h *Handler) initAPI(router chi.Router) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)
	router.Route("/api", func(api chi.Router) {
		handlerV1.Init(api)
	})
}
