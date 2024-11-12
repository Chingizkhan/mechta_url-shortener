package v1

import (
	"github.com/Chingizkhan/url-shortener/internal/dto"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"net/http"
)

// swagger:route POST /shortener shortener shortenUrlInRequest
// Генерация короткой ссылки
//
// <i>Генерация короткого уникального идентификатора для ссылки в ответ</i> <br/><br/>
// Время жизни коротких ссылок 30 дней
//
// responses:
// 200: shortenUrlResponse
// 500: errorResponseSwagger
func (h *Handler) generateShortenLink(w http.ResponseWriter, r *http.Request) {
	in := dto.ShortenURLIn{}
	if err := in.Parse(r.Body); err != nil {
		h.l.Error("parse", logger.Err(err))
		h.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := in.Validate(); err != nil {
		h.l.Error("validate", logger.Err(err))
		h.Error(w, err, http.StatusBadRequest)
		return
	}

	response, err := h.shortening.GenerateShortenUrl(r.Context(), in)
	if err != nil {
		h.l.Error("GenerateShortenUrl", logger.Err(err))
		h.Error(w, err, http.StatusInternalServerError)
		return
	}

	h.Resp(w, dto.ShortenURLOut{
		Link: response,
	}, http.StatusOK)
}
