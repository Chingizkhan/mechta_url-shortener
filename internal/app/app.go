package app

import (
	"github.com/Chingizkhan/url-shortener/config"
	v1 "github.com/Chingizkhan/url-shortener/internal/controller/http/v1"
	"github.com/Chingizkhan/url-shortener/internal/repository/url_repo"
	"github.com/Chingizkhan/url-shortener/internal/service/shortener"
	"github.com/Chingizkhan/url-shortener/internal/service/shortening"
	"github.com/Chingizkhan/url-shortener/internal/service/transactional"
	"github.com/Chingizkhan/url-shortener/pkg/httpserver"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"github.com/Chingizkhan/url-shortener/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config, l *logger.Logger) {
	pg, err := postgres.New(
		cfg.PG.DSN(),
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.MaxConnLifetime(time.Duration(cfg.PG.MaxConnLifetime)),
		postgres.MaxConnIdleTime(time.Duration(cfg.PG.MaxConnIdleTime)),
	)
	if err != nil {
		l.Error("app - Run - postgres.New:", logger.Err(err))
		os.Exit(1)
	}
	defer pg.Close()

	urlRepo := url_repo.New(pg)
	transactionalService := transactional.New(pg)
	shorteningService := shortening.New(
		l,
		shortener.New(),
		cfg.App.Host,
		&cfg.ExpiringProcessor.LinkExpire,
		cfg.App.LinkLength,
		urlRepo,
		transactionalService,
	)

	// start worker pool
	wp := newWorkerPool(params{
		repo:       urlRepo,
		logs:       cfg.ExpiringProcessor.ShowLogs,
		tickerTime: &cfg.ExpiringProcessor.TickerTimeout,
	})
	go wp.start()

	// start http server
	router := chi.NewRouter()
	handler := v1.NewHandler(&v1.HandlerParams{
		Logger:     l,
		Cfg:        cfg,
		Shortening: shorteningService,
	})
	handler.Register(router, time.Duration(cfg.HTTP.Timeout))
	httpServer := httpserver.New(
		router,
		httpserver.Port(cfg.HTTP.Port),
	)
	l.Info("http server started", slog.String("env", cfg.Log.Level), slog.String("port", cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal:", slog.String("signal", s.String()))
	case err = <-httpServer.Notify():
		l.Error("app - Run - http_server.Notify:", logger.Err(err))
	}

	// shutdown
	if err = httpServer.Shutdown(); err != nil {
		l.Error("app - Run - httpServer.Shutdown:", logger.Err(err))
		return
	}
}
