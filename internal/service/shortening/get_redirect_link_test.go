package shortening

import (
	"context"
	"errors"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/Chingizkhan/url-shortener/internal/repository/url_repo"
	"github.com/Chingizkhan/url-shortener/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRedirectLink(t *testing.T) {
	type mockBehavior func(urlRepo *mocks.UrlRepository, link string)

	tests := []struct {
		name          string
		link          string
		mockBehavior  mockBehavior
		expectedURL   string
		expectedError error
	}{
		{
			name: "successful get and update",
			link: "some-link",
			mockBehavior: func(urlRepo *mocks.UrlRepository, link string) {
				urlRepo.On("Get", mock.Anything, link).Once().Return(domain.Shortening{
					SourceURL: "http://example.com",
					Visits:    10,
				}, nil)
				urlRepo.On("Update", mock.Anything, url_repo.UpdateIn{
					Link:   link,
					Visits: 11,
				}).Once().Return(domain.Shortening{
					SourceURL: "http://example.com",
					Visits:    11,
				}, nil)
			},
			expectedURL:   "http://example.com",
			expectedError: nil,
		},
		{
			name: "get method error",
			link: "some-link",
			mockBehavior: func(urlRepo *mocks.UrlRepository, link string) {
				urlRepo.
					On("Get", mock.Anything, link).
					Once().
					Return(domain.Shortening{}, errors.New("get error"))
			},
			expectedURL:   "",
			expectedError: errors.New("get: get error"),
		},
		{
			name: "update method error",
			link: "some-link",
			mockBehavior: func(urlRepo *mocks.UrlRepository, link string) {
				urlRepo.On("Get", mock.Anything, link).Once().Return(domain.Shortening{
					SourceURL: "http://example.com",
					Visits:    10,
				}, nil)
				urlRepo.On("Update", mock.Anything, url_repo.UpdateIn{
					Link:   link,
					Visits: 11,
				}).Once().Return(domain.Shortening{}, errors.New("update error"))
			},
			expectedURL:   "",
			expectedError: errors.New("update: update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlRepo := new(mocks.UrlRepository)

			tt.mockBehavior(urlRepo, tt.link)

			svc := &Service{
				url: urlRepo,
			}

			resultURL, err := svc.getRedirectLink(context.Background(), tt.link)

			assert.Equal(t, tt.expectedURL, resultURL)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			urlRepo.AssertExpectations(t)
		})
	}
}
