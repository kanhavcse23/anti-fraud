package util_middleware_v1

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

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

		middlewareHandler.logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("API triggered")

		handler(w, r)
	}
}
