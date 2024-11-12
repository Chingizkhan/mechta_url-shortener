package shortening

import (
	"context"
	"errors"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/Chingizkhan/url-shortener/internal/dto"
	"github.com/Chingizkhan/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGenerateShortenUrl(t *testing.T) {
	type mockBehavior func(urlRepo *mocks.UrlRepository, shortener *mocks.Shortener, l *mocks.ILogger, linkGenerated string)

	tests := []struct {
		name          string
		input         dto.ShortenURLIn
		generatedLink string
		mockBehavior  mockBehavior
		expectedURL   string
		expectedError error
	}{
		{
			name: "successful link creation",
			input: dto.ShortenURLIn{
				URL: "http://example.com",
			},
			mockBehavior: func(urlRepo *mocks.UrlRepository, shortener *mocks.Shortener, l *mocks.ILogger, linkGenerated string) {
				urlRepo.On("Exists", mock.Anything, linkGenerated).Return(false, nil)
				urlRepo.On("Create", mock.Anything, mock.AnythingOfType("CreateIn")).Return(domain.Shortening{
					Link: linkGenerated,
				}, nil)
				shortener.
					On("ShortenURL", 6, mock.Anything).
					Return(linkGenerated)
			},
			expectedURL:   "http://localhost:8888/v1/JLywKv",
			expectedError: nil,
		},
		{
			name: "exists method error",
			input: dto.ShortenURLIn{
				URL: "http://example.com",
			},
			mockBehavior: func(urlRepo *mocks.UrlRepository, shortener *mocks.Shortener, l *mocks.ILogger, linkGenerated string) {
				shortener.
					On("ShortenURL", 6, mock.Anything).
					Return(linkGenerated)
				urlRepo.On("Exists", mock.Anything, linkGenerated).Return(true, errors.New("exists error"))
			},
			expectedURL:   "",
			expectedError: errors.New("exists: exists error"),
		},
		{
			name: "create method error",
			input: dto.ShortenURLIn{
				URL: "http://example.com",
			},
			mockBehavior: func(urlRepo *mocks.UrlRepository, shortener *mocks.Shortener, l *mocks.ILogger, linkGenerated string) {
				shortener.
					On("ShortenURL", 6, mock.Anything).
					Return(linkGenerated)
				urlRepo.On("Exists", mock.Anything, linkGenerated).Return(false, nil)
				urlRepo.On("Create", mock.Anything, mock.Anything).Return(domain.Shortening{}, errors.New("create error"))
			},
			expectedURL:   "",
			expectedError: errors.New("create url via repo: create error"),
		},
		{
			name: "regenerate link on duplicate",
			input: dto.ShortenURLIn{
				URL: "http://example.com",
			},
			mockBehavior: func(urlRepo *mocks.UrlRepository, shortener *mocks.Shortener, l *mocks.ILogger, linkGenerated string) {
				l.On("Info", mock.AnythingOfType("string")).Once()
				shortener.
					On("ShortenURL", 6, mock.Anything).
					Return(linkGenerated)
				urlRepo.On("Exists", mock.Anything, linkGenerated).Once().Return(true, nil)
				urlRepo.On("Exists", mock.Anything, linkGenerated).Return(false, nil)
				urlRepo.On("Create", mock.Anything, mock.Anything).Return(domain.Shortening{Link: linkGenerated}, nil)
			},
			expectedURL:   "http://localhost:8888/v1/" + "JLywKv",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlRepo := new(mocks.UrlRepository)
			shortener := new(mocks.Shortener)
			l := new(mocks.ILogger)

			linkGenerated := "JLywKv"
			tt.mockBehavior(urlRepo, shortener, l, linkGenerated)

			svc := &Service{
				l:          l,
				shortener:  shortener,
				expiration: 24 * time.Hour,
				host:       "http://localhost:8888/v1/",
				linkLen:    6,
				url:        urlRepo,
			}

			resultURL, err := svc.generateShortenUrl(context.Background(), tt.input)

			// Проверка результата
			assert.Equal(t, tt.expectedURL, resultURL)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверка выполнения мока
			urlRepo.AssertExpectations(t)
			shortener.AssertExpectations(t)
		})
	}
}
