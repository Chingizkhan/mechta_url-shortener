package url_repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) Exists(ctx context.Context, link string) (out bool, err error) {
	sql, args, err := sq.
		Select("true").
		From("url").
		Where(sq.Eq{"link": link}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return out, fmt.Errorf("to sql: %w", err)
	}

	sql = fmt.Sprintf("SELECT EXISTS (%s)", sql)

	if err = r.QueryRow(ctx, &out, sql, args...); err != nil {
		return out, fmt.Errorf("query row: %w", err)
	}

	return out, nil
}
