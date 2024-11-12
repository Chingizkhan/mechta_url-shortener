package dto

import (
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// swagger:parameters deleteURLRequest
	DeleteURLIn struct {
		// in:path
		// required: true
		ID string `json:"id"`
	}

	// swagger:response deleteURLResponse
	deleteURLResponse struct {
		Body DeleteURLOut `json:"body"`
	}

	DeleteURLOut struct {
		// example: ok
		Status string `json:"status"`
	}
)

func (in *DeleteURLIn) Parse(r *http.Request) error {
	link := chi.URLParam(r, linkURLParam)
	if link == "" {
		return domain.ErrEmptyLink
	}
	in.ID = link
	return nil
}
