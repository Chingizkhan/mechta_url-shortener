package url_repo

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExists(t *testing.T) {
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

	tests := []struct {
		name          string
		link          string
		expectedOut   bool
		expectedError error
	}{
		{
			name:          "link exists",
			link:          link,
			expectedOut:   true,
			expectedError: nil,
		},
		{
			name:          "link does not exist",
			link:          "non-existent-link",
			expectedOut:   false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := rep.Exists(ctx, tt.link)

			assert.Equal(t, tt.expectedOut, out)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
