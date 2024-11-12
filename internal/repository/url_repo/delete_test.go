package url_repo

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	ctx, tx, cancel := prependCtx(t)
	defer cancel()
	defer tx.Rollback(ctx)

	now := time.Now()

	link := "JLywKt"
	sourceUrl := "https://example.com"

	shortening, err := rep.Create(ctx, CreateIn{
		Link:      link,
		SourceURL: sourceUrl,
		ExpireAt: pgtype.Timestamp{
			Time:  now.Add(time.Minute * 30),
			Valid: true,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, shortening)
	assert.Equal(t, link, shortening.Link)
	assert.Equal(t, sourceUrl, shortening.SourceURL)
	assert.Equal(t, 0, shortening.Visits)
	assert.NotZero(t, now.Add(time.Minute*30), shortening.ExpireAt)

	tests := []struct {
		name          string
		link          string
		expectedError error
	}{
		{
			name:          "success",
			link:          link,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = rep.Delete(ctx, link)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
