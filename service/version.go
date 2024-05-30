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
	"time"
)

type createVersionRequest struct {
	Id         string    `json:"id"`
	Version    string    `json:"version"`
	Url        string    `json:"url"`
	UploadDate time.Time `json:"uploadDate"`
}

func CreateVersion(reader io.Reader, applicationId string, trackId string) (version *models.Version, errDetails *ErrorDetails, status int) {
	status = http.StatusCreated
	body := new(createVersionRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &ErrorDetails{
				EntityType: "version",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &ErrorDetails{
				EntityType: "version",
				ErrorType:  "serialization",
				Message:    "Unknown serialization error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}

		return
	}

	vsn := models.Version{
		ApplicationId: applicationId,
		TrackId:       trackId,
		Version:       body.Version,
		Url:           body.Url,
		UploadDate:    body.UploadDate,
	}

	version, err = models.CreateVersion(vsn)
	if err != nil {
		errDetails = &ErrorDetails{
			EntityType: "version",
		}

		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrVersionEmpty) {
			status = http.StatusBadRequest
			errDetails.ErrorType = "request"
			errDetails.Message = "Version missing"
		} else if errors.As(err, &pgErr) {
			errDetails.ErrorType = "database"

			if pgErr.Code == pgerrcode.UniqueViolation {
				status = http.StatusConflict
				errDetails.Message = "Version already exists"
			} else if pgErr.Code == pgerrcode.ForeignKeyViolation {
				status = http.StatusNotFound
				errDetails.Message = "Application or Track not found"
			} else {
				status = http.StatusInternalServerError
				errDetails.Message = "Unknown database error"
				log.Println(err.Error())
			}
		} else if errors.Is(err, models.ErrApplicationNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Application not found"
		} else if errors.Is(err, models.ErrTrackNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Track not found"
		} else {
			status = http.StatusInternalServerError
			errDetails.Message = "Unknown error"
			errDetails.ErrorType = "server"
			log.Println(err.Error())
		}
	}

	return
}

func GetAllVersions(applicationId string, trackId string) (versions []models.Version, errDetails *ErrorDetails, status int) {
	versions, err := models.GetAllVersions(applicationId, trackId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			_, err := models.GetApplicationById(applicationId)
			if errors.Is(err, models.ErrApplicationNotFound) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
				errDetails = &ErrorDetails{
					EntityType: "version",
					ErrorType:  "database",
					Message:    "Could not find application",
				}
				status = http.StatusNotFound
			} else {
				_, err := models.GetTrackById(trackId, applicationId)
				if errors.Is(err, models.ErrTrackNotFound) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
					errDetails = &ErrorDetails{
						EntityType: "version",
						ErrorType:  "database",
						Message:    "Could not find track",
					}
					status = http.StatusNotFound
				}
			}
		} else {
			errDetails = &ErrorDetails{
				EntityType: "version",
				ErrorType:  "database",
				Message:    "Could not get all versions",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}
	}

	return
}

func GetVersionById(applicationId string, trackId string, id string) (version *models.Version, errDetails *ErrorDetails, status int) {
	status = http.StatusOK
	version, err := models.GetVersionById(applicationId, trackId, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			errDetails = &ErrorDetails{
				EntityType: "version",
				ErrorType:  "database",
				Message:    "Could not find version",
			}
			status = http.StatusNotFound
		} else if errors.Is(err, models.ErrApplicationNotFound) {
			errDetails = &ErrorDetails{
				EntityType: "version",
				ErrorType:  "database",
				Message:    "Could not find application",
			}
			status = http.StatusNotFound
		} else if errors.Is(err, models.ErrTrackNotFound) {
			errDetails = &ErrorDetails{
				EntityType: "version",
				ErrorType:  "database",
				Message:    "Could not find track",
			}
			status = http.StatusNotFound
		} else {
			errDetails = &ErrorDetails{
				EntityType: "version",
				ErrorType:  "server",
				Message:    "Unknown error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}
	}

	return
}

func DeleteVersion(applicationId string, trackId string, id string) (errDetails *ErrorDetails, status int) {
	err := models.DeleteVersionById(applicationId, trackId, id)
	status = http.StatusNoContent

	if err != nil {
		errDetails = &ErrorDetails{
			EntityType: "version",
		}
		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrApplicationNotFound) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Application not found"
		} else if errors.Is(err, models.ErrTrackNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Track not found"
		} else if errors.Is(err, models.ErrVersionNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Version not found"
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
