package domain

import (
	"time"
)

type (
	Shortening struct {
		Link      string    `json:"link"`
		SourceURL string    `json:"source_url"`
		Visits    int       `json:"visits"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		ExpireAt  time.Time `json:"expire_at"`
	}
)
