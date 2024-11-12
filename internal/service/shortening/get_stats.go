package shortening

import (
	"context"
	"github.com/Chingizkhan/url-shortener/internal/domain"
)

func (s *Service) Get(ctx context.Context, id string) (out domain.Shortening, err error) {
	return s.url.Get(ctx, id)
}
