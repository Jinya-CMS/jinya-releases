package models

import (
	"errors"
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

func CreateVersion(version Version) (*Version, error) {
	if version.Version == "" {
		return nil, ErrVersionEmpty
	}

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO version (id, application_id, track_id, version, url, upload_date) VALUES ($1, $2, $3, $4)", version.Id, version.Version, version.Url, version.UploadDate)

	if err != nil {
		return nil, err
	}

	return GetVersionByTrackAndVersion(version.TrackId, version.Version)
}

func GetVersionByTrackAndVersion(trackId string, versionString string) (*Version, error) {
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

	if err = db.Get(version, "SELECT id, application_id, track_id, version.version, url, upload_date  FROM version WHERE track_id = $1 AND version = $2", trackId, versionString); err != nil {
		return nil, err
	}

	return version, nil
}

func GetVersionById(applicationId string, trackId string, id string) (*Version, error) {
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

	if err = db.Get(version, "SELECT id, application_id, track_id, version.version, url, upload_date  FROM version WHERE id = $1", id); err != nil {
		return nil, err
	}

	return version, nil
}

func GetAllVersions(applicationId string, trackId string) ([]Version, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	versions := make([]Version, 0)
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

	if err = db.Select(&versions, "SELECT id, application_id, track_id, version.version, url, upload_date FROM version WHERE application_id = $1 AND track_id = $2 ORDER BY upload_date", applicationId, trackId); err != nil {
		return nil, err
	}

	return versions, nil
}

func DeleteVersionById(id string, applicationId string, trackId string) error {
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
