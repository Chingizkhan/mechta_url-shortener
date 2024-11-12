package v1

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) routes(r *chi.Mux) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/shortener", h.generateShortenLink)
		r.Get("/shortener", h.linkList)
		r.Get("/{link}", h.linkRedirect)
		r.Delete("/{link}", h.linkDelete)
		r.Get("/stats/{link}", h.linkStats)
	})
}
