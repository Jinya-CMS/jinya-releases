package models

import (
	"errors"
	"fmt"
	"jinya-releases/config"
	"jinya-releases/database"
	"time"
)

type Version struct {
	Id            string    `json:"id" db:"id"`
	ApplicationId string    `json:"-" db:"application_id"`
	TrackId       string    `json:"-" db:"track_id"`
	Version       string    `json:"version" db:"version"`
	Url           string    `json:"url,omitempty" db:"url"`
	UploadDate    time.Time `json:"uploadDate,omitempty" db:"upload_date"`
}

var (
	ErrVersionEmpty    = errors.New("version is empty")
	ErrVersionNotFound = errors.New("version not found")
)

const getVersionsSelectAndJoin = `select v.id,
       v.application_id,
       v.track_id,
       v.version,
       v.upload_date,
       '%s/content/version/' || a.slug || '/' || t.slug || '/' || v.version as url
from version v
         join track t on t.id = v.track_id
         join application a on a.id = t.application_id
 %s`

func CreateVersion(version Version) (*Version, error) {
	if version.Version == "" {
		return nil, ErrVersionEmpty
	}

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	applicationCount := 0
	trackCount := 0
	if err = db.Get(&applicationCount, "SELECT COUNT(*) FROM application WHERE id = $1", version.ApplicationId); err != nil {
		return nil, err
	}

	if applicationCount == 0 {
		return nil, ErrApplicationNotFound
	}

	if err = db.Get(&trackCount, "SELECT COUNT(*) FROM track WHERE id = $1", version.TrackId); err != nil {
		return nil, err
	}

	if applicationCount == 0 {
		return nil, ErrTrackNotFound
	}

	_, err = db.Exec("INSERT INTO version (application_id, track_id, version, upload_date) VALUES ($1, $2, $3, $4)", version.ApplicationId, version.TrackId, version.Version, time.Now())

	if err != nil {
		return nil, err
	}

	return GetVersionByTrackAndVersion(version.TrackId, version.Version)
}

func GetVersionByTrackAndVersion(trackId, versionString string) (*Version, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	version := new(Version)
	trackCount := 0

	if err = db.Get(&trackCount, "SELECT COUNT(*) FROM track WHERE id = $1", trackId); err != nil {
		return nil, err
	}

	if trackCount == 0 {
		return nil, ErrTrackNotFound
	}

	err = db.Get(version, fmt.Sprintf(getVersionsSelectAndJoin, config.LoadedConfiguration.ServerUrl, "where v.track_id = $1 and v.version = $2"), trackId, versionString)

	return version, err
}

func GetVersionById(applicationId, trackId, id string) (*Version, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	version := new(Version)
	applicationCount := 0
	trackCount := 0

	if err = db.Get(&applicationCount, "SELECT COUNT(*) FROM application WHERE id = $1", applicationId); err != nil {
		return nil, err
	}

	if applicationCount == 0 {
		return nil, ErrApplicationNotFound
	}

	if err = db.Get(&trackCount, "SELECT COUNT(*) FROM track WHERE id = $1", trackId); err != nil {
		return nil, err
	}

	if trackCount == 0 {
		return nil, ErrTrackNotFound
	}

	err = db.Get(version, fmt.Sprintf(getVersionsSelectAndJoin, config.LoadedConfiguration.ServerUrl, "where v.id = $1"), id)

	return version, err
}

func GetAllVersions(applicationId, trackId string) ([]Version, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	versions := make([]Version, 0)

	if _, err = GetApplicationById(applicationId); err != nil {
		return nil, err
	}

	if _, err = GetTrackById(trackId, applicationId); err != nil {
		return nil, err
	}

	err = db.Select(&versions, fmt.Sprintf(getVersionsSelectAndJoin, config.LoadedConfiguration.ServerUrl, "where v.application_id = $1 and v.track_id = $2 order by v.version desc"), applicationId, trackId)

	return versions, err
}

func DeleteVersionById(applicationId, trackId, id string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	applicationCount := 0
	trackCount := 0

	if err = db.Get(&applicationCount, "SELECT COUNT(*) FROM application WHERE id = $1", applicationId); err != nil {
		return err
	}

	if applicationCount == 0 {
		return ErrApplicationNotFound
	}

	if err = db.Get(&trackCount, "SELECT COUNT(*) FROM track WHERE id = $1", trackId); err != nil {
		return err
	}

	if trackCount == 0 {
		return ErrTrackNotFound
	}

	result, err := db.Exec("DELETE FROM version WHERE id = $1 AND application_id = $2 AND track_id = $3", id, applicationId, trackId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrVersionNotFound
	}

	return nil
}

func GetVersionBySlugs(applicationSlug, trackSlug string) ([]Version, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	var versions []Version

	err = db.Select(&versions, fmt.Sprintf(getVersionsSelectAndJoin, config.LoadedConfiguration.ServerUrl, "where a.slug = $1 and t.slug = $2 order by v.version"), applicationSlug, trackSlug)

	return versions, err
}

func GetVersionBySlugsAndNumber(applicationSlug, trackSlug, versionNumber string) (*Version, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	version := new(Version)

	err = db.Get(version, fmt.Sprintf(getVersionsSelectAndJoin, config.LoadedConfiguration.ServerUrl, "where a.slug = $1 and t.slug = $2 and v.version = $3"), applicationSlug, trackSlug, versionNumber)

	return version, err
}
