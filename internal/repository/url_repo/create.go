package url_repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateIn struct {
	Link      string
	SourceURL string
	ExpireAt  pgtype.Timestamp
}

func (r *Repository) Create(ctx context.Context, in CreateIn) (out domain.Shortening, err error) {
	sql, args, err := sq.
		Insert("url").
		PlaceholderFormat(sq.Dollar).
		Columns("link", "source_url", "expire_at").
		Values(in.Link, in.SourceURL, in.ExpireAt).
		Suffix("returning *").
		ToSql()
	if err != nil {
		return out, fmt.Errorf("to sql: %w", err)
	}

	if err = r.QueryRow(ctx, &out, sql, args...); err != nil {
		if err = r.checkCreateOrderConstraints(err); err != nil {
			return out, err
		}
		return out, fmt.Errorf("query row: %w", err)
	}

	return out, nil
}

const (
	constraintUrlPKey = "url_pkey"
)

func (r *Repository) checkCreateOrderConstraints(err error) error {
	var e *pgconn.PgError
	if errors.As(err, &e) {
		switch e.Code {

		case pgerrcode.UniqueViolation:
			switch e.ConstraintName {
			case constraintUrlPKey:
				return domain.ErrAlreadyExists
			default:
				return fmt.Errorf("database error: %w", errors.New(e.ConstraintName))
			}
		}
	}
	return err
}
