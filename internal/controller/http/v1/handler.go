package v1

import (
	"context"
	"github.com/Chingizkhan/url-shortener/config"
	customMiddleware "github.com/Chingizkhan/url-shortener/internal/controller/http/middleware"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"github.com/Chingizkhan/url-shortener/internal/dto"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"time"
)

type (
	ShorteningService interface {
		GenerateShortenUrl(ctx context.Context, in dto.ShortenURLIn) (link string, err error)
		GetRedirectLink(ctx context.Context, id string) (dto.RedirectURLOut, error)
		List(ctx context.Context) (out []domain.Shortening, err error)
		Get(ctx context.Context, id string) (out domain.Shortening, err error)
		Delete(ctx context.Context, id string) error
	}

	Handler struct {
		l          logger.ILogger
		cfg        *config.Config
		shortening ShorteningService
	}

	HandlerParams struct {
		Logger     logger.ILogger
		Cfg        *config.Config
		Shortening ShorteningService
	}
)

func NewHandler(p *HandlerParams) *Handler {
	return &Handler{
		l:          p.Logger,
		cfg:        p.Cfg,
		shortening: p.Shortening,
	}
}

func (h *Handler) Register(r *chi.Mux, timeout time.Duration) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(timeout))
	r.Use(customMiddleware.Cors)
	r.Use(customMiddleware.Logging(h.l))

	h.routes(r)
	h.swaggerDocs(r)
}
