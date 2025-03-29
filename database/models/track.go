package models

import (
	"database/sql"
	"errors"
	"jinya-releases/database"
)

type Track struct {
	Id            string `json:"id" db:"id"`
	ApplicationId string `json:"-" db:"application_id"`
	Name          string `json:"name" db:"name"`
	Slug          string `json:"slug" db:"slug"`
	IsDefault     bool   `json:"isDefault" db:"is_default"`
}

var (
	ErrTrackNotFound = errors.New("track not found")
)

func resetDefault(db *sql.Tx) error {
	_, err := db.Exec("UPDATE track SET is_default = false WHERE is_default = true")

	return err
}

func CreateTrack(track Track) (*Track, error) {
	if track.Name == "" {
		return nil, ErrNameEmpty
	}
	if track.Slug == "" {
		return nil, ErrSlugEmpty
	}

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	if track.IsDefault {
		err = resetDefault(tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	_, err = db.Exec("INSERT INTO track (application_id, name, slug, is_default) VALUES ($1, $2, $3, $4)", track.ApplicationId, track.Name, track.Slug, track.IsDefault)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return GetTrackBySlug(track.Slug, track.ApplicationId)
}

func GetAllTracks(applicationId string) ([]Track, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	tracks := make([]Track, 0)
	applicationCount := 0

	if err = db.Get(&applicationCount, "SELECT COUNT(*) FROM application WHERE id = $1", applicationId); err != nil {
		return nil, err
	}

	if applicationCount == 0 {
		return nil, ErrApplicationNotFound
	}

	if err = db.Select(&tracks, "SELECT id, application_id, name, slug, is_default FROM track WHERE application_id = $1 ORDER BY name", applicationId); err != nil {
		return nil, err
	}

	return tracks, nil
}

func GetTrackById(id, applicationId string) (*Track, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	track := new(Track)
	applicationCount := 0

	if err = db.Get(&applicationCount, "SELECT COUNT(*) FROM application WHERE id = $1", applicationId); err != nil {
		return nil, err
	}

	if applicationCount == 0 {
		return nil, ErrApplicationNotFound
	}

	if err = db.Get(track, "SELECT id, application_id, name, slug, is_default FROM track WHERE id = $1 AND application_id = $2", id, applicationId); err != nil {
		return nil, err
	}

	return track, nil
}

func GetTrackBySlug(slug, applicationId string) (*Track, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	track := new(Track)
	applicationCount := 0

	if err = db.Get(&applicationCount, "SELECT COUNT(*) FROM application WHERE id = $1", applicationId); err != nil {
		return nil, err
	}

	if applicationCount == 0 {
		return nil, ErrApplicationNotFound
	}

	if err = db.Get(track, "SELECT id, application_id, name, slug, is_default FROM track WHERE slug = $1 AND application_id = $2", slug, applicationId); err != nil {
		return nil, err
	}

	return track, nil
}

func UpdateTrack(track Track) (*Track, error) {
	if track.Name == "" {
		return nil, ErrNameEmpty
	}
	if track.Slug == "" {
		return nil, ErrSlugEmpty
	}

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	if track.IsDefault {
		err = resetDefault(tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	result, err := tx.Exec("UPDATE track SET name = $1, slug = $2, is_default = $3 WHERE id = $4 AND application_id = $5", track.Name, track.Slug, track.IsDefault, track.Id, track.ApplicationId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if affected == 0 {
		tx.Rollback()
		return nil, ErrTrackNotFound
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return GetTrackBySlug(track.Slug, track.ApplicationId)
}

func DeleteTrackById(id, applicationId string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	applicationCount := 0

	if err = db.Get(&applicationCount, "SELECT COUNT(*) FROM application WHERE id = $1", applicationId); err != nil {
		return err
	}

	if applicationCount == 0 {
		return ErrApplicationNotFound
	}

	defer db.Close()

	result, err := db.Exec("DELETE FROM track WHERE id = $1 AND application_id = $2", id, applicationId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrTrackNotFound
	}

	return nil
}
