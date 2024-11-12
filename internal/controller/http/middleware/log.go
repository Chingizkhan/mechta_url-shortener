package middleware

import (
	"fmt"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"net/http"
)

func Logging(l logger.ILogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info(fmt.Sprintf("[%s]: %s", r.Method, r.RequestURI))

			next.ServeHTTP(w, r)
		})
	}
}
