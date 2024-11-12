package shortening

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/domain"
)

func (s *Service) Delete(ctx context.Context, link string) error {
	if err := s.tx.Exec(ctx, func(txCtx context.Context) error {
		return s.delete(ctx, link)
	}); err != nil {
		return fmt.Errorf("exec transaction: %w", err)
	}

	return nil
}

func (s *Service) delete(ctx context.Context, link string) error {
	exists, err := s.url.Exists(ctx, link)
	if err != nil {
		return fmt.Errorf("exists: %w", err)
	}

	if !exists {
		return domain.ErrNotFound
	}

	if err = s.url.Delete(ctx, link); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
