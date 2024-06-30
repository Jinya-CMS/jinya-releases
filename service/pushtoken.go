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

type createPushtokenRequest struct {
	Id          string   `json:"id"`
	Token       string   `json:"token"`
	AllowedApps []string `json:"allowedApps,omitempty"`
}

type updatePushtokenRequest struct {
	Id          string   `json:"id"`
	Token       string   `json:"token"`
	AllowedApps []string `json:"allowedApps,omitempty"`
}

func CreatePushtoken(reader io.Reader) (pushtoken *models.PushToken, errDetails *utils.ErrorDetails, status int) {
	body := new(createPushtokenRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &utils.ErrorDetails{
				EntityType: "pushtoken",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &utils.ErrorDetails{
				EntityType: "pushtoken",
				ErrorType:  "serialization",
				Message:    "Unknown serialization error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}

		return
	}

	token := models.PushToken{
		AllowedApps: body.AllowedApps,
	}

	pushtoken, err = models.CreatePushtoken(token.AllowedApps)
	status = http.StatusCreated
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "pushtoken",
		}

		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrApplicationlistEmpty) {
			status = http.StatusBadRequest
			errDetails.ErrorType = "request"
			errDetails.Message = "Applicationlist empty"
		} else if errors.Is(err, models.ErrApplicationNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Application not found"
		} else if errors.As(err, &pgErr) {
			errDetails.ErrorType = "database"
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

func GetAllPushtokens() (pushtokens []models.PushToken, errDetails *utils.ErrorDetails) {
	pushtokens, err := models.GetAllPushTokens()

	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "pushtoken",
			ErrorType:  "database",
			Message:    "Could not get all pushtokens",
		}
	}

	return
}

func GetPushtokenById(id string) (pushtoken *models.PushToken, errDetails *utils.ErrorDetails, status int) {
	pushtoken, err := models.GetPushTokenById(id)
	status = http.StatusOK
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			errDetails = &utils.ErrorDetails{
				EntityType: "pushtoken",
				ErrorType:  "database",
				Message:    "Could not find pushtoken",
			}
			status = http.StatusNotFound
		} else {
			errDetails = &utils.ErrorDetails{
				EntityType: "pushtoken",
				ErrorType:  "server",
				Message:    "Unknown error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}
	}

	return
}

func UpdatePushtoken(id string, reader io.Reader) (pushtoken *models.PushToken, errDetails *utils.ErrorDetails, status int) {
	status = http.StatusNoContent

	body := new(updatePushtokenRequest)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(body)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			errDetails = &utils.ErrorDetails{
				EntityType: "pushtoken",
				ErrorType:  "request",
				Message:    "Json syntax error",
			}
			status = http.StatusBadRequest
		} else {
			errDetails = &utils.ErrorDetails{
				EntityType: "pushtoken",
				ErrorType:  "serialization",
				Message:    "Unknown serialization error",
			}
			status = http.StatusInternalServerError
			log.Println(err.Error())
		}

		return
	}

	token := models.PushToken{
		Id:          id,
		Token:       body.Token,
		AllowedApps: body.AllowedApps,
	}

	pushtoken, err = models.UpdatePushtoken(token.Id, token.AllowedApps)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "pushtoken",
		}

		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrPushtokenNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Pushtoken not found"
			return
		} else if errors.Is(err, models.ErrApplicationlistEmpty) {
			status = http.StatusBadRequest
			errDetails.ErrorType = "request"
			errDetails.Message = "Applicationlist empty"
			return
		} else if errors.Is(err, models.ErrApplicationNotFound) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Application not found"
			return
		} else if errors.As(err, &pgErr) {
			errDetails.ErrorType = "database"
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

func DeletePushtoken(id string) (errDetails *utils.ErrorDetails, status int) {
	err := models.DeletePushtoken(id)
	status = http.StatusNoContent

	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "pushtoken",
		}
		var pgErr *pgconn.PgError
		if errors.Is(err, models.ErrPushtokenNotFound) || (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation) {
			status = http.StatusNotFound
			errDetails.ErrorType = "request"
			errDetails.Message = "Pushtoken not found"
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
