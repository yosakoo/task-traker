package v1

import (


	"github.com/go-chi/chi/v5"
	"github.com/yosakoo/task-traker/internal/service"
	"github.com/yosakoo/task-traker/pkg/auth"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
	validate     *validator.Validate
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		validate:     validator.New(),
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(router chi.Router) {
    router.Group(func(v1 chi.Router) {
        h.initUsersRoutes(v1)
		h.initTasksRoutes(v1)
    })
}