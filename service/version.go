package service

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"jinya-releases/database/models"
	"jinya-releases/storage"
	"jinya-releases/utils"
	"log"
	"net/http"
	"time"
)

func GetAllVersions(applicationId, trackId string) (versions []models.Version, errDetails *utils.ErrorDetails, status int) {
	versions, err := models.GetAllVersions(applicationId, trackId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			_, err := models.GetApplicationById(applicationId)
			if errors.Is(err, models.ErrApplicationNotFound) || errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
				errDetails = &utils.ErrorDetails{
					EntityType: "version",
					ErrorType:  "database",
					Message:    "Could not find application",
				}
				status = http.StatusNotFound
			} else {
				_, err := models.GetTrackById(trackId, applicationId)
				if errors.Is(err, models.ErrTrackNotFound) || errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
					errDetails = &utils.ErrorDetails{
						EntityType: "version",
						ErrorType:  "database",
						Message:    "Could not find track",
					}
					status = http.StatusNotFound
				}
			}
		} else {
			errDetails = &utils.ErrorDetails{
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

func GetVersionById(applicationId, trackId, id string) (version *models.Version, errDetails *utils.ErrorDetails, status int) {
	status = http.StatusOK
	version, err := models.GetVersionById(applicationId, trackId, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			errDetails = &utils.ErrorDetails{
				EntityType: "version",
				ErrorType:  "database",
				Message:    "Could not find version",
			}
			status = http.StatusNotFound
		} else if errors.Is(err, models.ErrApplicationNotFound) {
			errDetails = &utils.ErrorDetails{
				EntityType: "version",
				ErrorType:  "database",
				Message:    "Could not find application",
			}
			status = http.StatusNotFound
		} else if errors.Is(err, models.ErrTrackNotFound) {
			errDetails = &utils.ErrorDetails{
				EntityType: "version",
				ErrorType:  "database",
				Message:    "Could not find track",
			}
			status = http.StatusNotFound
		} else {
			errDetails = &utils.ErrorDetails{
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

func DeleteVersion(applicationId, trackId, id string) (errDetails *utils.ErrorDetails, status int) {
	err := models.DeleteVersionById(applicationId, trackId, id)
	status = http.StatusNoContent

	if err != nil {
		errDetails = &utils.ErrorDetails{
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

	err = storage.DeleteVersion(applicationId, trackId, id)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "version",
			ErrorType:  "storage",
			Message:    "Failed to delete binary",
		}
		log.Println(err.Error())
		status = http.StatusInternalServerError
		return
	}
	return
}

func UploadVersion(r *http.Request, applicationId, trackId, versionNumber string) (errDetails *utils.ErrorDetails, status int) {
	app, err := models.GetApplicationById(applicationId)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "version",
			Message:    "Application not found",
			ErrorType:  "request",
		}
		status = http.StatusNotFound

		return
	}

	track, err := models.GetTrackById(trackId, app.Id)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "version",
			Message:    "Track not found",
			ErrorType:  "request",
		}
		status = http.StatusNotFound

		return
	}

	return performUpload(r, app, track, versionNumber)
}

func PushVersion(r *http.Request, applicationSlug, trackSlug, versionNumber string) (errDetails *utils.ErrorDetails, status int) {
	app, err := models.GetApplicationBySlug(applicationSlug)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "version",
			Message:    "Application not found",
			ErrorType:  "request",
		}
		status = http.StatusNotFound

		return
	}

	track, err := models.GetTrackBySlug(trackSlug, app.Id)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "version",
			Message:    "Track not found",
			ErrorType:  "request",
		}
		status = http.StatusNotFound

		return
	}

	return performUpload(r, app, track, versionNumber)
}

func performUpload(r *http.Request, app *models.Application, track *models.Track, versionNumber string) (errDetails *utils.ErrorDetails, status int) {
	status = http.StatusNoContent
	versionToUploadBinaryFor, err := models.GetVersionByTrackAndVersion(track.Id, versionNumber)
	if err != nil {
		versionToUploadBinaryFor, err = models.CreateVersion(models.Version{
			ApplicationId: app.Id,
			TrackId:       track.Id,
			Version:       versionNumber,
			UploadDate:    time.Now(),
		})

		if err != nil {
			status = http.StatusNotFound
			errDetails = &utils.ErrorDetails{
				EntityType: "version",
				Message:    "Version not found and cannot be created",
				ErrorType:  "request",
			}

			return
		}
	}

	err = storage.UploadVersion(r.Body, r.ContentLength, r.Header.Get("Content-Type"), versionToUploadBinaryFor)
	if err != nil {
		status = http.StatusInternalServerError
		errDetails = &utils.ErrorDetails{
			EntityType: "version",
			Message:    "Version not found and cannot be created",
			ErrorType:  "request",
		}
	}

	return
}
