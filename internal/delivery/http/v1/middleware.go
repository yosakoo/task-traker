package v1

import (
	"net/http"
	"time"
	"fmt"
	"context"
	"strconv"
    "strings"

	"github.com/yosakoo/task-traker/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

func NewMwLogger(log logger.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("Middleware логгера включен")

		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				log.Info(fmt.Sprintf("method: %s, path: %s, remote_addr: %s, user_agent: %s, request_id: %s, status: %d, bytes: %d, duration: %s",
					r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), middleware.GetReqID(r.Context()), ww.Status(), ww.BytesWritten(), time.Since(t1).String()))
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}


func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var accessToken string

        authHeader := r.Header.Get("Authorization")
        if authHeader != "" {
            accessToken = strings.TrimPrefix(authHeader, "Bearer ")
        }
        
        if accessToken == "" {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("not authenticated"))
            return
        }

        userIdStr, err := h.tokenManager.Parse(accessToken)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("not authenticated"))
            return
        }

        userId, err := strconv.Atoi(userIdStr)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte("internal server error"))
            return
        }

        ctx := context.WithValue(r.Context(), "user_id", userId)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
