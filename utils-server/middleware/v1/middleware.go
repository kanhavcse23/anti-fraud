package util_middleware_v1

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type contextKey string

const requestIDKey contextKey = "requestID"

type MiddlewareHandler struct {
	logger *logrus.Logger
}

func NewMiddlewareHandler(logger *logrus.Logger) *MiddlewareHandler {
	return &MiddlewareHandler{logger: logger}
}

// wrapper func to handle error for HTTP methods
func (middlewareHandler *MiddlewareHandler) MiddlewareHandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				middlewareHandler.logger.Errorf("Recovered from panic: %v", rec)
				http.Error(w, "An internal server error occurred", http.StatusInternalServerError)
			}

		}()
		reqID := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)

		middlewareHandler.logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"request_id": reqID,
		}).Info("API triggered")

		handler(w, r.WithContext(ctx))
	}
}

// Helper to retrieve the request ID from context in controllers or anywhere else
func GetRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(requestIDKey).(string); ok {
		return val
	}
	return ""
}
