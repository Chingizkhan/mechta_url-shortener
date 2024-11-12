package url_repo

import (
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFindExpired(t *testing.T) {
	ctx, tx, cancel := prependCtx(t)
	defer cancel()
	defer tx.Rollback(ctx)

	now := time.Now().UTC()

	expLink1, err := rep.Create(ctx, CreateIn{
		Link:      "expiredLink1",
		SourceURL: "https://example.com/expired1",
		ExpireAt: pgtype.Timestamp{
			Time:  now.Add(-time.Hour),
			Valid: true,
		},
	})
	assert.NoError(t, err)

	expLink2, err := rep.Create(ctx, CreateIn{
		Link:      "expiredLink2",
		SourceURL: "https://example.com/expired2",
		ExpireAt: pgtype.Timestamp{
			Time:  now.Add(-time.Minute * 30),
			Valid: true,
		},
	})
	assert.NoError(t, err)

	_, err = rep.Create(ctx, CreateIn{
		Link:      "validLink",
		SourceURL: "https://example.com/valid",
		ExpireAt: pgtype.Timestamp{
			Time:  now.Add(time.Hour),
			Valid: true,
		},
	})
	assert.NoError(t, err)

	tests := []struct {
		name          string
		in            FindIn
		expectedOut   []domain.Shortening
		expectedError error
	}{
		{
			name: "find expired links",
			in:   FindIn{Limit: 10},
			expectedOut: []domain.Shortening{
				expLink1,
				expLink2,
			},
			expectedError: nil,
		},
		{
			name:          "no expired links",
			in:            FindIn{Limit: 10},
			expectedOut:   nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "no expired links" {
				tx.Rollback(ctx)
				ctx, tx, cancel = prependCtx(t)
				defer cancel()
				defer tx.Rollback(ctx)
			}

			out, err := rep.FindExpired(ctx, tt.in)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOut, out)
			}
		})
	}
}
