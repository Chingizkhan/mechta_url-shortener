package app

import (
	"context"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/Chingizkhan/url-shortener/internal/repository/url_repo"
	"log"
	"sync"
	"time"
)

const (
	_defaultTickerTime = time.Minute * 30
	_limit             = 100
	_workersSize       = 3
)

type (
	UrlRepo interface {
		FindExpired(ctx context.Context, in url_repo.FindIn) (out []domain.Shortening, err error)
		Delete(ctx context.Context, link string) error
	}

	WorkerPool struct {
		ticker         *time.Ticker
		ch             chan domain.Shortening
		size           int
		loggerSwitcher bool
		wg             *sync.WaitGroup
		repo           UrlRepo
	}

	params struct {
		repo       UrlRepo
		logs       bool
		tickerTime *time.Duration
	}
)

func newWorkerPool(p params) *WorkerPool {
	tt := _defaultTickerTime

	if p.tickerTime != nil {
		tt = *p.tickerTime
	}
	return &WorkerPool{
		size:           _workersSize,
		ch:             make(chan domain.Shortening, _limit),
		ticker:         time.NewTicker(tt),
		wg:             &sync.WaitGroup{},
		loggerSwitcher: p.logs,
		repo:           p.repo,
	}
}

func (wp *WorkerPool) start() {
	for range wp.size {
		go wp.worker()
	}

	wp.distribute()

	for {
		select {
		case <-wp.ticker.C:
			wp.distribute()
		}
	}
}

func (wp *WorkerPool) worker() {
	for val := range wp.ch {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

		if err := wp.repo.Delete(ctx, val.Link); err != nil {
			wp.log("error on deleting shortening: %s", val.Link)
		}
		cancel()
	}
}

func (wp *WorkerPool) distribute() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	shortenings, err := wp.repo.FindExpired(ctx, url_repo.FindIn{
		Limit: _limit,
	})
	if err != nil {
		wp.log("error on finding shortenings:", err)
		return
	}

	for _, shortening := range shortenings {
		wp.ch <- shortening
	}
}

func (wp *WorkerPool) log(msg string, v ...any) {
	if wp.loggerSwitcher {
		log.Printf("[WORKER_POOL] "+msg, v...)
	}
}

func (wp *WorkerPool) stop() {
	wp.ticker.Stop()
	close(wp.ch)
}
