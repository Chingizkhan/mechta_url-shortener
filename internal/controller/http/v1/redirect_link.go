package v1

import (
	"github.com/Chingizkhan/url-shortener/internal/dto"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"net/http"
)

// swagger:route POST /{link} link linkRedirectRequest
// редирект
//
// responses:
// 307:
// 500: errorResponseSwagger
func (h *Handler) linkRedirect(w http.ResponseWriter, r *http.Request) {
	in := dto.RedirectURLIn{}
	if err := in.Parse(r); err != nil {
		h.l.Error("parse", logger.Err(err))
		h.Error(w, err, http.StatusBadRequest)
		return
	}

	response, err := h.shortening.GetRedirectLink(r.Context(), in.ID)
	if err != nil {
		h.l.Error("get", logger.Err(err))
		h.Error(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, response.Link, http.StatusTemporaryRedirect)
}
