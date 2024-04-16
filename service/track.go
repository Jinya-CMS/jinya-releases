package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"io"
	"jinya-releases/database/models"
	"log"
	"net/http"
)

type createTrackRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
}

type updateTrackRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
}

func GetAllTracks(applicationId string) (tracks []models.Track, errDetails *ErrorDetails) {
	tracks, err := models.GetAllTracks(applicationId)

	if err != nil {
		errDetails = &ErrorDetails{
			EntityType: "track",
			ErrorType:  "database",
			Message:    "Could not get all tracks",
		}
	}

	return
}

func GetTrackById(trackId string, applicationId string) (track *models.Track, errDetails *ErrorDetails, status int) {
	track, err := models.GetTrackById(trackId, applicationId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			errDetails = &ErrorDetails{
				EntityType: "track",
				ErrorType:  "database",
				Message:    "Could not find track",
			}
			status = http.StatusNotFound
		} else {
			errDetails = &ErrorDetails{
				EntityType: "track",
				ErrorType:  "server",
				Message:    "Unknown error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}
	}

	return
}

func CreateTrack(reader io.Reader) (track *models.Track, errDetails *ErrorDetails, status int) {
	body := new(createTrackRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &ErrorDetails{
				EntityType: "track",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &ErrorDetails{
				EntityType: "track",
				ErrorType:  "serialization",
				Message:    "Unknown serialization error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}

		return
	}

	trk := models.Track{
		Name:      body.Name,
		Slug:      body.Slug,
		IsDefault: body.IsDefault,
	}

	track, err = models.CreateTrack(trk)
	if err != nil {
		errDetails = &ErrorDetails{
			EntityType: "track",
		}

		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrSlugEmpty) {
			status = http.StatusBadRequest
			errDetails.ErrorType = "request"
			errDetails.Message = "Slug missing"
		} else if errors.Is(err, models.ErrNameEmpty) {
			status = http.StatusBadRequest
			errDetails.ErrorType = "request"
			errDetails.Message = "Name missing"
		} else if errors.As(err, &pgErr) {
			errDetails.ErrorType = "database"

			if pgErr.Code == pgerrcode.UniqueViolation {
				status = http.StatusConflict
				errDetails.Message = "Track already exists"
			} else {
				status = http.StatusInternalServerError
				errDetails.Message = "Unknown database error"
				log.Println(err.Error())
			}
		} else {
			status = http.StatusInternalServerError
			errDetails.Message = "Unknown error"
			errDetails.ErrorType = "server"
			log.Println(err.Error())
		}
	}

	return
}

func UpdateTrack(trackId string, applicationId string, reader io.Reader) (track *models.Track, errDetails *ErrorDetails, status int) {
	status = http.StatusNoContent

	body := new(updateTrackRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &ErrorDetails{
				EntityType: "track",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &ErrorDetails{
				EntityType: "track",
				ErrorType:  "serialization",
				Message:    "Unknown serialization error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}

		return
	}

	trk := models.Track{
		Id:            trackId,
		ApplicationId: applicationId,
		Name:          body.Name,
		Slug:          body.Slug,
		IsDefault:     body.IsDefault,
	}

	track, err = models.UpdateTrack(trk)
	if err != nil {
		errDetails = &ErrorDetails{
			EntityType: "track",
		}

		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrTrackNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Track not found"
			return
		} else if errors.As(err, &pgErr) {
			errDetails.ErrorType = "database"

			if pgErr.Code == pgerrcode.UniqueViolation {
				status = http.StatusConflict
				errDetails.Message = "Track already exists"
			} else {
				status = http.StatusInternalServerError
				errDetails.Message = "Unknown database error"
				log.Println(err.Error())
			}
		} else {
			status = http.StatusInternalServerError
			errDetails.Message = "Unknown error"
			errDetails.ErrorType = "server"
			log.Println(err.Error())
		}
	}
	return
}

func DeleteTrack(trackId string, applicationId string) (errDetails *ErrorDetails, status int) {
	err := models.DeleteTrackById(trackId, applicationId)
	status = http.StatusNoContent

	if err != nil {
		errDetails = &ErrorDetails{
			EntityType: "track",
		}
		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrTrackNotFound) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Track not found"
		} else if errors.As(err, &pgErr) {
			status = http.StatusInternalServerError
			errDetails.Message = "Unknown database error"
			log.Println(err.Error())
		} else {
			status = http.StatusInternalServerError
			errDetails.Message = "Unknown error"
			errDetails.ErrorType = "server"
			log.Println(err.Error())
		}
	}
	return
}
