package dto

import (
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// swagger:parameters linkRedirectRequest
	RedirectURLIn struct {
		// in:path
		// required: true
		ID string `json:"id"`
	}

	RedirectURLOut struct {
		Link string
	}
)

func (in *RedirectURLIn) Parse(r *http.Request) error {
	link := chi.URLParam(r, linkURLParam)
	if link == "" {
		return domain.ErrEmptyLink
	}
	in.ID = link
	return nil
}
