package url_repo

import (
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	ctx, tx, cancel := prependCtx(t)
	defer cancel()
	defer tx.Rollback(ctx)

	tests := []struct {
		name          string
		link          string
		sourceUrl     string
		expectedError error
	}{
		{
			name:          "success",
			link:          "JLywKb",
			sourceUrl:     "https://example.com",
			expectedError: nil,
		},
		{
			name:          "already_exists",
			link:          "JLywKb",
			sourceUrl:     "https://example.com",
			expectedError: domain.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()

			shortening, err := rep.Create(ctx, CreateIn{
				Link:      tt.link,
				SourceURL: tt.sourceUrl,
				ExpireAt: pgtype.Timestamp{
					Time:  now.Add(time.Minute * 30),
					Valid: true,
				},
			})

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, shortening)
				assert.Equal(t, tt.link, shortening.Link)
				assert.Equal(t, tt.sourceUrl, shortening.SourceURL)
				assert.Equal(t, 0, shortening.Visits)
				assert.NotZero(t, now.Add(time.Minute*30), shortening.ExpireAt)
			}
		})
	}
}
