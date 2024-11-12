package url_repo

import (
	repo "github.com/Chingizkhan/url-shortener/internal/repository"
	"github.com/Chingizkhan/url-shortener/pkg/postgres"
)

type Repository struct {
	*repo.DefaultRepo
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{repo.NewDefaultRepo(pg.Pool)}
}
