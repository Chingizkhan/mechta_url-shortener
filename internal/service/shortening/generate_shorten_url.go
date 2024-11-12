package shortening

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/dto"
	"github.com/Chingizkhan/url-shortener/internal/repository/url_repo"
	"github.com/jackc/pgx/v5/pgtype"
	"math/rand"
	"time"
)

func (s *Service) GenerateShortenUrl(ctx context.Context, in dto.ShortenURLIn) (link string, err error) {
	if err = s.tx.Exec(ctx, func(txCtx context.Context) error {
		if link, err = s.generateShortenUrl(txCtx, in); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return link, fmt.Errorf("exec transaction: %w", err)
	}

	return link, nil
}

func (s *Service) generateShortenUrl(ctx context.Context, in dto.ShortenURLIn) (out string, err error) {
	var link string

	for {
		src := rand.NewSource(time.Now().UnixNano())
		link = s.shortener.ShortenURL(s.linkLen, src)

		exists, err := s.url.Exists(ctx, link)
		if err != nil {
			return out, fmt.Errorf("exists: %w", err)
		}

		if exists {
			s.l.Info("link already exists, regenerate")
			continue
		}
		break
	}

	shortening, err := s.url.Create(ctx, url_repo.CreateIn{
		Link:      link,
		SourceURL: in.URL,
		ExpireAt: pgtype.Timestamp{
			Time:  time.Now().UTC().Add(s.expiration),
			Valid: true,
		},
	})
	if err != nil {
		return out, fmt.Errorf("create url via repo: %w", err)
	}

	return s.host + shortening.Link, nil
}
