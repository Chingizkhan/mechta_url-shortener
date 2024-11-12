package shortening

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/dto"
	"github.com/Chingizkhan/url-shortener/internal/repository/url_repo"
)

func (s *Service) GetRedirectLink(ctx context.Context, link string) (out dto.RedirectURLOut, err error) {
	if err := s.tx.Exec(ctx, func(txCtx context.Context) error {
		out.Link, err = s.getRedirectLink(txCtx, link)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return out, fmt.Errorf("exec transaction: %w", err)
	}

	return out, nil
}

func (s *Service) getRedirectLink(ctx context.Context, link string) (string, error) {
	shortening, err := s.url.Get(ctx, link)
	if err != nil {
		return "", fmt.Errorf("get: %w", err)
	}

	shortening, err = s.url.Update(ctx, url_repo.UpdateIn{
		Link:   link,
		Visits: shortening.Visits + 1,
	})
	if err != nil {
		return "", fmt.Errorf("update: %w", err)
	}

	return shortening.SourceURL, nil
}
