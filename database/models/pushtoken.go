package models

import (
	"errors"
	"jinya-releases/database"
)

type PushToken struct {
	Id          string   `json:"id" db:"id"`
	Token       string   `json:"token" db:"token"`
	AllowedApps []string `json:"allowedApps,omitempty"`
}

var (
	ErrPushtokenNotFound = errors.New("token not found")
)

func CreatePushtoken(applications []string) (*PushToken, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	var tokenId string

	if err = db.Select(&tokenId, "INSERT INTO pushtoken (id, token) DEFAULT VALUES RETURNING id"); err != nil {
		return nil, err
	}

	pushToken := new(PushToken)
	if err = db.Get(pushToken, "SELECT id, token FROM pushtoken WHERE id = $1", tokenId); err != nil {
		return nil, err
	}

	for _, a := range applications {
		_, err = db.Exec("INSERT INTO pushtokenapplication (token, application) VALUES ($1, $2)", pushToken.Token, a)

		if err != nil {
			return nil, err
		}
	}

	pushToken.AllowedApps = make([]string, len(applications))
	for i, application := range applications {
		pushToken.AllowedApps[i] = application
	}

	return pushToken, nil
}

func GetAllPushTokens() ([]PushToken, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	pushTokens := make([]PushToken, 0)

	if err = db.Select(&pushTokens, "SELECT id, token FROM pushtoken"); err != nil {
		return nil, err
	}

	for i, pushtoken := range pushTokens {
		applicationIds := make([]string, 0)
		if err = db.Select(&applicationIds, "SELECT application FROM pushtokenapplication WHERE token = $1", pushtoken.Token); err != nil {
			return nil, err
		}
		pushTokens[i].AllowedApps = applicationIds

	}

	return pushTokens, nil
}

func GetPushTokenById(id string) (*PushToken, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	pushToken := new(PushToken)

	if err = db.Select(&pushToken, "SELECT id, token FROM pushtoken where id = $1", id); err != nil {
		return nil, err
	}

	applicationIds := make([]string, 0)
	if err = db.Select(&applicationIds, "SELECT application FROM pushtokenapplication WHERE token = $1", pushToken.Token); err != nil {
		return nil, err
	}
	pushToken.AllowedApps = applicationIds

	return pushToken, nil
}

func UpdatePushtoken(id string, applications []string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()
	var token string

	if err = db.Select(&token, "SELECT  token FROM pushtoken where id = $1", id); err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM pushtokenapplication WHERE token = $1", token)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrPushtokenNotFound
	}

	for _, a := range applications {
		_, err = db.Exec("INSERT INTO pushtokenapplication (token, application) VALUES ($1, $2)", token, a)

		if err != nil {
			return err
		}
	}

	return nil
}

func DeletePushtoken(id string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()
	var token string

	if err = db.Select(&token, "SELECT  token FROM pushtoken where id = $1", id); err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM pushtokenapplication WHERE token = $1", token)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrPushtokenNotFound
	}

	result, err = db.Exec("DELETE FROM pushtoken WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
