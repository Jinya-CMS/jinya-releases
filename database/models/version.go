package models

import "time"

type Version struct {
	Id         string    `json:"id"`
	Version    string    `json:"version"`
	Url        string    `json:"url,omitempty"`
	UploadDate time.Time `json:"uploadDate,omitempty"`
}
