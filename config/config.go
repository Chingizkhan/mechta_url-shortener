package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"net/url"
	"time"
)

type (
	Config struct {
		App               `json:"app" yaml:"app"`
		PG                `json:"pg" yaml:"postgres"`
		HTTP              `json:"http" yaml:"http"`
		Log               `json:"log" yaml:"logger"`
		ExpiringProcessor `json:"expiring_processor" yaml:"expiring_processor"`
	}

	App struct {
		Name       string `env-required:"true" yaml:"name" env:"APP_NAME" json:"name"`
		Version    string `env-required:"true" yaml:"version" env:"APP_VERSION" json:"version"`
		Level      Level  `yaml:"level" json:"level"`
		Host       string `env-required:"true" yaml:"host" env:"APP_HOST" json:"host"`
		LinkLength int
	}

	PG struct {
		User            string   `env-required:"true" yaml:"user" env:"PG_USER" json:"user"`
		Password        string   `env-required:"true" yaml:"password" env:"PG_PASSWORD" json:"password"`
		Host            string   `env-required:"true" yaml:"host" env:"PG_HOST" json:"host"`
		Port            string   `env-required:"true" yaml:"port" env:"PG_PORT" json:"port"`
		Name            string   `env-required:"true" yaml:"name" env:"PG_NAME" json:"name"`
		PoolMax         int      `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX" json:"pool_max"`
		MaxConnLifetime Duration `env-required:"true" yaml:"max_conn_lifetime" env:"PG_MAX_CONN_LIFETIME" json:"max_conn_lifetime"`
		MaxConnIdleTime Duration `env-required:"true" yaml:"max_conn_idle_time" env:"PG_MAX_CONN_IDLE_TIME" json:"max_conn_idle_time"`
		SSLMode         string   `env-required:"true" yaml:"ssl_mode" env:"PG_SSL_MODE" json:"ssl_mode"`
	}

	HTTP struct {
		Port    string   `env-required:"true" yaml:"port" env:"HTTP_PORT" json:"port"`
		Timeout Duration `env-required:"true" yaml:"timeout" env:"HTTP_TIMEOUT" json:"timeout"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL" json:"level"`
	}

	ExpiringProcessor struct {
		TickerTimeout time.Duration `yaml:"ticker_timeout" json:"ticker_timeout"`
		ShowLogs      bool          `yaml:"show_logs" yaml:"show_logs"`
		LinkExpire    time.Duration `json:"link_expire" yaml:"link_expire"`
	}
)

func New(url string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(url, cfg); err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read env error: %w", err)
	}

	return cfg, nil
}

func (pg *PG) DSN() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(pg.User, pg.Password),
		Host:   pg.Host + ":" + pg.Port,
		Path:   pg.Name,
	}

	q := dsn.Query()

	q.Add("sslmode", pg.SSLMode)

	dsn.RawQuery = q.Encode()

	return dsn.String()
}
