package v1

import (
	"github.com/Chingizkhan/url-shortener/internal/dto"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"net/http"
)

// swagger:route GET /{link} link statsURLRequest
// дает информацию
//
// <i>Генерация короткого уникального идентификатора для ссылки в ответ</i> <br/><br/>
// Время жизни коротких ссылок 30 дней
//
// responses:
// 200: statsURLResponse
// 500: errorResponseSwagger
func (h *Handler) linkStats(w http.ResponseWriter, r *http.Request) {
	in := dto.StatsURLIn{}
	if err := in.Parse(r); err != nil {
		h.l.Error("parse", logger.Err(err))
		h.Error(w, err, http.StatusBadRequest)
		return
	}

	response, err := h.shortening.Get(r.Context(), in.ID)
	if err != nil {
		h.l.Error("get", logger.Err(err))
		h.Error(w, err, http.StatusInternalServerError)
		return
	}

	h.Resp(w, dto.StatsURLResponse{
		SourceURL:       response.SourceURL,
		Visits:          response.Visits,
		CreatedAt:       response.CreatedAt,
		LastTimeVisited: response.UpdatedAt,
	}, http.StatusOK)
}
