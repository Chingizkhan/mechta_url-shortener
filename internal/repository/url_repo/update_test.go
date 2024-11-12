package url_repo

import (
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	ctx, tx, cancel := prependCtx(t)
	defer cancel()
	defer tx.Rollback(ctx)

	link := "JLywKt"
	sourceUrl := "https://example.com"
	initialVisits := 0

	shortening, err := rep.Create(ctx, CreateIn{
		Link:      link,
		SourceURL: sourceUrl,
		ExpireAt: pgtype.Timestamp{
			Time:  time.Now().Add(time.Minute * 30),
			Valid: true,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, shortening)
	assert.Equal(t, initialVisits, shortening.Visits)

	tests := []struct {
		name          string
		updateIn      UpdateIn
		expectedOut   domain.Shortening
		expectedError error
	}{
		{
			name: "successful update",
			updateIn: UpdateIn{
				Link:   link,
				Visits: initialVisits + 1,
			},
			expectedOut: domain.Shortening{
				Link:      link,
				SourceURL: sourceUrl,
				Visits:    initialVisits + 1,
			},
			expectedError: nil,
		},
		{
			name: "non-existent link",
			updateIn: UpdateIn{
				Link:   "non-existent-link",
				Visits: 1,
			},
			expectedOut:   domain.Shortening{},
			expectedError: fmt.Errorf("query_row: scanning one: no rows in result set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Вызываем тестируемую функцию
			out, err := rep.Update(ctx, tt.updateIn)

			// Проверяем результат
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOut.Link, out.Link)
				assert.Equal(t, tt.expectedOut.SourceURL, out.SourceURL)
				assert.Equal(t, tt.updateIn.Visits, out.Visits)
			}
		})
	}
}
