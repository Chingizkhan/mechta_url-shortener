package url_repo

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type UpdateIn struct {
	Link   string
	Visits int
}

func (r *Repository) Update(ctx context.Context, in UpdateIn) (out domain.Shortening, err error) {
	sql, args, err := sq.
		Update("url").
		PlaceholderFormat(sq.Dollar).
		SetMap(map[string]interface{}{
			"visits":     in.Visits,
			"updated_at": time.Now(),
		}).
		Where(sq.Eq{"link": in.Link}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return out, fmt.Errorf("to_sql: %w", err)
	}

	if err = r.QueryRow(ctx, &out, sql, args...); err != nil {
		return out, fmt.Errorf("query_row: %w", err)
	}

	return out, nil
}
