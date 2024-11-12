// REST API для сервиса сокращателя ссылок.
//
//		REST API для сервиса сокращателя ссылок.
//
//		 Базовая структура <b>URL: {scheme}://{host}/{version}/{actor}/...</b><br/>
//		     &ensp; - <b>scheme</b> - протокол по которому работает запрос.<br/>
//		     &ensp; - <b>host</b> - домен/сервер на который идет запрос.<br/>
//		     &ensp; - <b>version</b> - версия апи.<br/>
//		     &ensp; - <b>actor</b> - тип клиента для которого данный эндпоинт предназначен.<br/>
//
//		 Типы actor:<br/>
//		     &ensp; - <b>shortener</b> - сам сокращатель ссылок<br/>
//		     &ensp; - <b>link</b> - сокращенная ссылка<br/><br/>
//
//		 <h3>По умолчанию ссылка создается с 6 символами, вы можете это конфигурировать передав флаг <i>-link_len</i> с необходимым количеством символов</h3>
//
//		     Schemes: http, https
//
//		     Version: 1.0.0
//
//		     Consumes:
//		     - application/json
//
//		     Produces:
//		     - application/json
//
//	          SecurityDefinitions:
//	          Bearer:
//	             type: apiKey
//	             name: Authorization
//	             in: header
//
// swagger:meta
package main

import (
	"flag"
	"github.com/Chingizkhan/url-shortener/config"
	"github.com/Chingizkhan/url-shortener/internal/app"
	"github.com/Chingizkhan/url-shortener/pkg/logger"
	"log"
)

// todo: embed config also as file
func main() {
	linkLen := flag.Int("link_len", 6, "used to specify length of output link")
	flag.Parse()

	cfg, err := config.New("./config/config.yml")
	if err != nil {
		log.Fatalf("new config: %s", err)
	}
	cfg.App.LinkLength = *linkLen
	l := logger.New(cfg.Log.Level)
	l.Debug("debug messages are enabled")
	app.Migrate(cfg, l)
	app.Run(cfg, l)
}
