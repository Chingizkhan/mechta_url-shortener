package url_repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Get(ctx context.Context, link string) (out domain.Shortening, err error) {
	sql, args, err := sq.
		Select("*").
		From("url").
		Where(sq.Eq{"link": link}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return out, fmt.Errorf("to sql: %w", err)
	}

	if err = r.QueryRow(ctx, &out, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, domain.ErrNotFound
		}
		return out, fmt.Errorf("query row: %w", err)
	}

	return out, nil
}
