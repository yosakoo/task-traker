package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/yosakoo/task-traker/internal/domain"
	"github.com/yosakoo/task-traker/internal/service"
	
)

func (h *Handler) initUsersRoutes(router chi.Router) {
	router.Route("/users", func(r chi.Router) {
		r.Post("/sign-up", h.userSignUp)
		r.Post("/sign-in", h.userSignIn)
		r.Post("/auth/refresh", h.userRefresh)

		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware)
			r.Get("/", h.getCurrentUser)
		})
	})
}

type userSignUpInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type userSignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type authUser struct {
	ID int      `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type refreshInput struct {
	Token string `json:"token"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) userSignUp(w http.ResponseWriter, r *http.Request) {
	var input userSignUpInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}
	if err := h.validate.Struct(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	ctx := r.Context()
	res, err := h.services.Users.SignUp(ctx, service.UserSignUpInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("this email is already taken"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		w.Write([]byte("could not sign up user"))
		return
	}
	
	response := tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
	
}
func (h *Handler) userSignIn(w http.ResponseWriter, r *http.Request) {
    var input userSignInInput
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("invalid request body"))
        return
    }
	if err := h.validate.Struct(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

    ctx := r.Context()
    res, err := h.services.Users.SignIn(ctx, service.UserSignInInput{
        Email:    input.Email,
        Password: input.Password,
    })
    if err != nil {
		fmt.Println(err)
		if errors.Is(err, domain.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not found"))
			return
		}
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("could not sign in user"))
        return
    }
    response := tokenResponse{
        AccessToken:  res.AccessToken,
        RefreshToken: res.RefreshToken,
    }
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("could not marshal response"))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(jsonResponse)
}

func (h *Handler) userRefresh(w http.ResponseWriter, r *http.Request) {
	var input refreshInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}

	ctx := r.Context()
	fmt.Println(input.Token)

	res, err := h.services.Users.RefreshTokens(ctx, input.Token)
	if err != nil {
		if errors.Is(err, domain.ErrTokenExpired) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("token has expired"))
			return
		}

		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not refresh"))
		return
	}

	response := tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)

}

func (h *Handler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	
    userId := r.Context().Value("user_id").(int)
    user, err := h.services.Users.GetUserByID(r.Context(), userId)
    if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not found"))
			return
		}
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("could not get user"))
        return
    }
    jsonResponse, err := json.Marshal(user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("could not marshal response"))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}
