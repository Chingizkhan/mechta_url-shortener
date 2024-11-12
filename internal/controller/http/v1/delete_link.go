package v1

import (
	"github.com/Chingizkhan/url-shortener/internal/dto"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"net/http"
)

// swagger:route DELETE /{link} link deleteURLRequest
// Удаляет ссылку
//
// <i>Генерация короткого уникального идентификатора для ссылки в ответ</i> <br/><br/>
// Время жизни коротких ссылок 30 дней
//
// responses:
// 200: deleteURLResponse
// 500: errorResponseSwagger
func (h *Handler) linkDelete(w http.ResponseWriter, r *http.Request) {
	in := dto.DeleteURLIn{}
	if err := in.Parse(r); err != nil {
		h.l.Error("parse", logger.Err(err))
		h.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := h.shortening.Delete(r.Context(), in.ID); err != nil {
		h.l.Error("delete", logger.Err(err))
		h.Error(w, err, http.StatusInternalServerError)
		return
	}

	h.Resp(w, dto.DeleteURLOut{
		Status: "ok",
	}, http.StatusOK)
}
