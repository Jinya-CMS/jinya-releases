package database

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	Id   uuid.UUID `json:"id" db:"id,primarykey"`
	Name string    `json:"name" db:"name,unique"`
	Logo *string   `json:"logo" db:"logo"`
	Slug string    `json:"slug" db:"slug,unique"`
}

type PushToken struct {
	Id            uuid.UUID `json:"id" db:"id,primarykey"`
	Token         string    `json:"token" db:"token,unique"`
	ApplicationId string    `json:"-" db:"application_id"`
}

type Track struct {
	Id            uuid.UUID `json:"id" db:"id,primarykey"`
	ApplicationId string    `json:"-" db:"application_id"`
	Name          string    `json:"name" db:"name"`
	Slug          string    `json:"slug" db:"slug"`
	IsDefault     bool      `json:"isDefault" db:"is_default"`
}

type Version struct {
	Id            uuid.UUID `json:"id" db:"id,primarykey"`
	ApplicationId string    `json:"-" db:"application_id"`
	TrackId       string    `json:"-" db:"track_id"`
	Version       string    `json:"version" db:"version"`
	Url           string    `json:"url,omitempty" db:"url"`
	UploadDate    time.Time `json:"uploadDate,omitempty" db:"upload_date"`
}
