package url_repo

import (
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	ctx, tx, cancel := prependCtx(t)
	defer cancel()
	defer tx.Rollback(ctx)

	link := "JLywKt"
	sourceUrl := "https://example.com"
	expireAt := time.Now().UTC().Add(time.Minute * 30)

	shortening, err := rep.Create(ctx, CreateIn{
		Link:      link,
		SourceURL: sourceUrl,
		ExpireAt: pgtype.Timestamp{
			Time:  expireAt,
			Valid: true,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, shortening)

	tests := []struct {
		name          string
		link          string
		expectedOut   domain.Shortening
		expectedError error
	}{
		{
			name: "successful retrieval",
			link: link,
			expectedOut: domain.Shortening{
				Link:      link,
				SourceURL: sourceUrl,
				ExpireAt:  expireAt.UTC(),
			},
			expectedError: nil,
		},
		{
			name:          "link not found",
			link:          "non-existent-link",
			expectedOut:   domain.Shortening{},
			expectedError: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := rep.Get(ctx, tt.link)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOut.Link, out.Link)
				assert.Equal(t, tt.expectedOut.SourceURL, out.SourceURL)
				assert.Equal(t, tt.expectedOut.ExpireAt, out.ExpireAt)
			}
		})
	}
}
