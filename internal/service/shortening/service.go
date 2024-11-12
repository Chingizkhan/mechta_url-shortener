package shortening

import (
	"context"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/Chingizkhan/url-shortener/internal/repository/url_repo"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"math/rand"
	"time"
)

const (
	_linkLength = 6
	_expiration = time.Hour * 24 * 30
)

type (
	Service struct {
		l          logger.ILogger
		shortener  Shortener
		expiration time.Duration
		host       string
		linkLen    int
		url        UrlRepository
		tx         TransactionalService
	}

	TransactionalService interface {
		Exec(ctx context.Context, fn func(txCtx context.Context) error) error
	}

	UrlRepository interface {
		Exists(ctx context.Context, link string) (out bool, err error)
		Get(ctx context.Context, id string) (out domain.Shortening, err error)
		Create(ctx context.Context, in url_repo.CreateIn) (out domain.Shortening, err error)
		Update(ctx context.Context, in url_repo.UpdateIn) (out domain.Shortening, err error)
		List(ctx context.Context) (out []domain.Shortening, err error)
		Delete(ctx context.Context, link string) error
	}

	Shortener interface {
		ShortenURL(n int, src rand.Source) string
	}
)

func New(
	l logger.ILogger,
	shortener Shortener,
	host string,
	expiration *time.Duration,
	linkLen int,
	url UrlRepository,
	tx TransactionalService,
) *Service {
	var (
		linkLength = _linkLength
		exp        = _expiration
	)

	if linkLen != 0 {
		linkLength = linkLen
	}

	if expiration != nil {
		exp = *expiration
	}

	return &Service{
		l:          l,
		shortener:  shortener,
		host:       host,
		expiration: exp,
		linkLen:    linkLength,
		url:        url,
		tx:         tx,
	}
}
