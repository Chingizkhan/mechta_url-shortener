package url_repo

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	sq "github.com/Masterminds/squirrel"
)

type FindIn struct {
	Limit uint64
}

func (r *Repository) FindExpired(ctx context.Context, in FindIn) (out []domain.Shortening, err error) {
	sql, args, err := sq.
		Select("*").
		From("url").
		PlaceholderFormat(sq.Dollar).
		Where("expire_at <= current_timestamp").
		OrderBy("expire_at").
		Limit(in.Limit).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("to sql: %w", err)
	}

	if err = r.Query(ctx, &out, sql, args...); err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return out, nil
}
