package shortening

import (
	"context"
	"github.com/Chingizkhan/url-shortener/internal/domain"
)

func (s *Service) List(ctx context.Context) (out []domain.Shortening, err error) {
	return s.url.List(ctx)
}
