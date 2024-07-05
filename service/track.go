package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/utils"
	"log"
	"net/http"
)

type createTrackRequest struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
}

type updateTrackRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
}

func GetAllTracks(applicationId string) (tracks []models.Track, errDetails *utils.ErrorDetails, status int) {
	tracks, err := models.GetAllTracks(applicationId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			_, err := models.GetApplicationById(applicationId)
			if errors.Is(err, models.ErrApplicationNotFound) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
				errDetails = &utils.ErrorDetails{
					EntityType: "track",
					ErrorType:  "database",
					Message:    "Could not find application",
				}
				status = http.StatusNotFound
			}
		} else {
			errDetails = &utils.ErrorDetails{
				EntityType: "track",
				ErrorType:  "database",
				Message:    "Could not get all tracks",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}
	}

	return
}

func GetTrackById(trackId, applicationId string) (track *models.Track, errDetails *utils.ErrorDetails, status int) {
	status = http.StatusOK
	track, err := models.GetTrackById(trackId, applicationId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			errDetails = &utils.ErrorDetails{
				EntityType: "track",
				ErrorType:  "database",
				Message:    "Could not find track",
			}
			status = http.StatusNotFound
		} else if errors.Is(err, models.ErrApplicationNotFound) {
			errDetails = &utils.ErrorDetails{
				EntityType: "track",
				ErrorType:  "database",
				Message:    "Could not find application",
			}
			status = http.StatusNotFound
		} else {
			errDetails = &utils.ErrorDetails{
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

func CreateTrack(reader io.Reader, applicationId string) (track *models.Track, errDetails *utils.ErrorDetails, status int) {
	status = http.StatusCreated
	body := new(createTrackRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &utils.ErrorDetails{
				EntityType: "track",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &utils.ErrorDetails{
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
		ApplicationId: applicationId,
		Name:          body.Name,
		Slug:          body.Slug,
		IsDefault:     body.IsDefault,
	}

	track, err = models.CreateTrack(trk)
	if err != nil {
		errDetails = &utils.ErrorDetails{
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
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				status = http.StatusNotFound
				errDetails.Message = "Application not found"
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

func UpdateTrack(trackId, applicationId string, reader io.Reader) (track *models.Track, errDetails *utils.ErrorDetails, status int) {
	status = http.StatusNoContent

	body := new(updateTrackRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &utils.ErrorDetails{
				EntityType: "track",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &utils.ErrorDetails{
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
		errDetails = &utils.ErrorDetails{
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
		} else if errors.Is(err, models.ErrTrackNotFound) {
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

func DeleteTrack(trackId, applicationId string) (errDetails *utils.ErrorDetails, status int) {
	err := models.DeleteTrackById(trackId, applicationId)
	status = http.StatusNoContent

	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "track",
		}
		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrTrackNotFound) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Track not found"
		} else if errors.Is(err, models.ErrApplicationNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Application not found"
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
