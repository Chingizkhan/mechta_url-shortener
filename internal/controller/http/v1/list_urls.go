package v1

import (
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"net/http"
)

// swagger:route GET /shortener shortener linkListRequest
// Выводит список ссылок пользователя
//
// <i>Генерация короткого уникального идентификатора для ссылки в ответ</i> <br/><br/>
// Время жизни коротких ссылок 30 дней
//
// responses:
// 200: shortenUrlResponse
// 500: errorResponseSwagger
func (h *Handler) linkList(w http.ResponseWriter, r *http.Request) {
	response, err := h.shortening.List(r.Context())
	if err != nil {
		h.l.Error("shortening list", logger.Err(err))
		h.Error(w, err, http.StatusInternalServerError)
		return
	}

	h.Resp(w, response, http.StatusOK)
}
