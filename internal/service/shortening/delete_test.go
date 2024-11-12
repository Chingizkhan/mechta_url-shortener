package shortening

import (
	"context"
	"errors"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/Chingizkhan/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDelete(t *testing.T) {
	type mockBehavior func(urlRepo *mocks.UrlRepository, link string)

	tests := []struct {
		name          string
		link          string
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "successful delete",
			link: "some-link",
			mockBehavior: func(urlRepo *mocks.UrlRepository, link string) {
				urlRepo.On("Exists", mock.Anything, link).Once().Return(true, nil)
				urlRepo.On("Delete", mock.Anything, link).Once().Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "link does not exist",
			link: "non-existent-link",
			mockBehavior: func(urlRepo *mocks.UrlRepository, link string) {
				urlRepo.On("Exists", mock.Anything, link).Once().Return(false, nil)
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name: "exists method error",
			link: "some-link",
			mockBehavior: func(urlRepo *mocks.UrlRepository, link string) {
				urlRepo.On("Exists", mock.Anything, link).Once().Return(false, errors.New("exists error"))
			},
			expectedError: errors.New("exists: exists error"),
		},
		{
			name: "delete method error",
			link: "some-link",
			mockBehavior: func(urlRepo *mocks.UrlRepository, link string) {
				urlRepo.On("Exists", mock.Anything, link).Once().Return(true, nil)
				urlRepo.On("Delete", mock.Anything, link).Once().Return(errors.New("delete error"))
			},
			expectedError: errors.New("delete: delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlRepo := new(mocks.UrlRepository)
			txService := new(mocks.TransactionalService)

			tt.mockBehavior(urlRepo, tt.link)

			svc := &Service{
				url: urlRepo,
				tx:  txService,
			}

			err := svc.delete(context.Background(), tt.link)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			urlRepo.AssertExpectations(t)
		})
	}
}
