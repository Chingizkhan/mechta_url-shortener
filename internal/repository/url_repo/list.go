package url_repo

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) List(ctx context.Context) (out []domain.Shortening, err error) {
	sql, args, err := sq.
		Select("*").
		PlaceholderFormat(sq.Dollar).
		From("url").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return out, fmt.Errorf("to sql: %w", err)
	}

	if err = r.Query(ctx, &out, sql, args...); err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return out, nil
}
