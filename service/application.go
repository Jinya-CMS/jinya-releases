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

type createApplicationRequest struct {
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	HomepageTemplate     string `json:"homepageTemplate"`
	TrackpageTemplate    string `json:"trackpageTemplate"`
	AdditionalCss        string `json:"additionalCss,omitempty"`
	AdditionalJavaScript string `json:"additionalJavaScript,omitempty"`
}

type updateApplicationRequest struct {
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	HomepageTemplate     string `json:"homepageTemplate"`
	TrackpageTemplate    string `json:"trackpageTemplate"`
	AdditionalCss        string `json:"additionalCss,omitempty"`
	AdditionalJavaScript string `json:"additionalJavaScript,omitempty"`
}

func GetAllApplications() (applications []models.Application, errDetails *utils.ErrorDetails) {
	applications, err := models.GetAllApplications()

	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "application",
			ErrorType:  "database",
			Message:    "Could not get all applications",
		}
	}

	return
}

func GetApplicationById(id string) (application *models.Application, errDetails *utils.ErrorDetails, status int) {
	application, err := models.GetApplicationById(id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			errDetails = &utils.ErrorDetails{
				EntityType: "application",
				ErrorType:  "database",
				Message:    "Could not find application",
			}
			status = http.StatusNotFound
		} else {
			errDetails = &utils.ErrorDetails{
				EntityType: "application",
				ErrorType:  "server",
				Message:    "Unknown error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}
	}

	return
}

func CreateApplication(reader io.Reader) (application *models.Application, errDetails *utils.ErrorDetails, status int) {
	body := new(createApplicationRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &utils.ErrorDetails{
				EntityType: "application",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &utils.ErrorDetails{
				EntityType: "application",
				ErrorType:  "serialization",
				Message:    "Unknown serialization error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}

		return
	}

	app := models.Application{
		Name:                 body.Name,
		Slug:                 body.Slug,
		HomepageTemplate:     body.HomepageTemplate,
		TrackpageTemplate:    body.TrackpageTemplate,
		AdditionalCss:        &body.AdditionalCss,
		AdditionalJavaScript: &body.AdditionalJavaScript,
	}

	application, err = models.CreateApplication(app)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "application",
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
		} else if errors.Is(err, models.ErrHomepageTemplateEmpty) {
			status = http.StatusBadRequest
			errDetails.ErrorType = "request"
			errDetails.Message = "Homepage template missing"
		} else if errors.Is(err, models.ErrTrackpageTemplateEmpty) {
			status = http.StatusBadRequest
			errDetails.ErrorType = "request"
			errDetails.Message = "Trackpage template missing"
		} else if errors.As(err, &pgErr) {
			errDetails.ErrorType = "database"

			if pgErr.Code == pgerrcode.UniqueViolation {
				status = http.StatusConflict
				errDetails.Message = "Application already exists"
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

func UpdateApplication(id string, reader io.Reader) (application *models.Application, errDetails *utils.ErrorDetails, status int) {
	status = http.StatusNoContent

	body := new(updateApplicationRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &utils.ErrorDetails{
				EntityType: "application",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &utils.ErrorDetails{
				EntityType: "application",
				ErrorType:  "serialization",
				Message:    "Unknown serialization error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}

		return
	}

	app := models.Application{
		Id:                   id,
		Name:                 body.Name,
		Slug:                 body.Slug,
		HomepageTemplate:     body.HomepageTemplate,
		TrackpageTemplate:    body.TrackpageTemplate,
		AdditionalCss:        &body.AdditionalCss,
		AdditionalJavaScript: &body.AdditionalJavaScript,
	}

	application, err = models.UpdateApplication(app)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "application",
		}

		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrApplicationNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Application not found"
			return
		} else if errors.As(err, &pgErr) {
			errDetails.ErrorType = "database"

			if pgErr.Code == pgerrcode.UniqueViolation {
				status = http.StatusConflict
				errDetails.Message = "Application already exists"
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

func DeleteApplication(id string) (errDetails *utils.ErrorDetails, status int) {
	err := models.DeleteApplicationById(id)
	status = http.StatusNoContent

	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "application",
		}
		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrApplicationNotFound) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
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
