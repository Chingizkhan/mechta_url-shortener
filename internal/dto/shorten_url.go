package dto

import (
	"encoding/json"
	"github.com/Chingizkhan/url-shortener/internal/domain"
	"io"
	"net/url"
	"regexp"
)

type (
	// swagger:parameters shortenUrlRequest
	shortenURLInRequest struct {
		// in:body
		// required: true
		Body ShortenURLIn `json:"body"`
	}

	// swagger:response shortenUrlResponse
	shortenUrlResponse struct {
		// in:body
		Body ShortenURLOut `json:"body"`
	}

	ShortenURLIn struct {
		URL string `json:"url"`
	}

	ShortenURLOut struct {
		Link string `json:"link"`
	}

	// swagger:response errorResponseSwagger
	errorResponseSwagger struct {
		// body
		// in: body
		Body struct {
			// Код ошибки
			//
			// Example: already_exists
			Error string `json:"error"`
		}
	}
)

func (in *ShortenURLIn) Parse(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(in)
}

func (in *ShortenURLIn) Validate() error {
	if len(in.URL) == 0 {
		return domain.ErrTooShort
	}

	if !IsUrl(in.URL) {
		return domain.ErrInvalidURL
	}

	return nil
}

func IsUrl(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// URL содержит корректный хост (доменное имя или IP-адрес)
	hostRegex := regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(:[0-9]{1,5})?$|^(\d{1,3}\.){3}\d{1,3}(:[0-9]{1,5})?$`)
	return hostRegex.MatchString(parsedURL.Host)
}
