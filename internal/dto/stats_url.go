package dto

import (
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type (
	// swagger:parameters statsURLRequest
	StatsURLIn struct {
		// in:path
		// required: true
		ID string `json:"id"`
	}

	// swagger:response statsURLResponse
	statsURLResponse struct {
		Body domain.Shortening `json:"body"`
	}

	StatsURLResponse struct {
		SourceURL       string    `json:"source_url"`
		Visits          int       `json:"visits"`
		CreatedAt       time.Time `json:"created_at"`
		LastTimeVisited time.Time `json:"last_time_visited"`
	}
)

func (in *StatsURLIn) Parse(r *http.Request) error {
	link := chi.URLParam(r, linkURLParam)
	if link == "" {
		return domain.ErrEmptyLink
	}
	in.ID = link
	return nil
}
