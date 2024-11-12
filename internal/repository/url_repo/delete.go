package url_repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) Delete(ctx context.Context, link string) error {
	sql, args, err := sq.
		Delete("url").
		Where(sq.Eq{"link": link}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("to sql: %w", err)
	}

	if err = r.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
