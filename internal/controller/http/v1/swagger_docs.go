package v1

import (
	urlShortener "github.com/Chingizkhan/url-shortener"
	"github.com/go-chi/chi/v5"
	swagmiddleware "github.com/go-openapi/runtime/middleware"
	"net/http"
)

func (h *Handler) swaggerDocs(r *chi.Mux) {
	r.Handle("/", http.RedirectHandler("/docs", http.StatusMovedPermanently))
	r.Handle("/swagger.yaml", http.FileServer(http.FS(urlShortener.SwaggerSpecs)))
	opts := swagmiddleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "docs", Title: "REST API для сервиса кошелька."}
	sh := swagmiddleware.Redoc(opts, nil)
	r.Handle("/docs", sh)
}
