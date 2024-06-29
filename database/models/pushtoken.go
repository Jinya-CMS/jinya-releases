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
	ErrPushtokenNotFound    = errors.New("token not found")
	ErrApplicationlistEmpty = errors.New("applicationlist empty")
)

func CreatePushtoken(applications []string) (*PushToken, error) {
	if len(applications) == 0 {
		return nil, ErrApplicationlistEmpty
	}
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	for _, a := range applications {
		count := 0
		err = db.Get(&count, "SELECT COUNT(*) FROM application WHERE id = $1", a)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, ErrApplicationNotFound
		}
	}

	var tokenId string

	if err = db.Get(&tokenId, "INSERT INTO pushtoken DEFAULT VALUES RETURNING id"); err != nil {
		return nil, err
	}

	pushToken := new(PushToken)
	if err = db.Get(pushToken, "SELECT id, token FROM pushtoken WHERE id = $1", tokenId); err != nil {
		return nil, err
	}

	for _, app := range applications {
		_, err = db.Exec("INSERT INTO pushtoken_application (token_id, application_id) VALUES ($1, $2)", pushToken.Token, app)

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
		if err = db.Select(&applicationIds, "SELECT application_id FROM pushtoken_application WHERE token_id = $1", pushtoken.Token); err != nil {
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
	type tempToken struct {
		Id    string `db:"id"`
		Token string `db:"token"`
	}
	tToken := new(tempToken)

	if err = db.Get(tToken, "SELECT id, token FROM pushtoken where id = $1", id); err != nil {
		return nil, err
	}

	pushToken := new(PushToken)
	pushToken.Id = tToken.Id
	pushToken.Token = tToken.Token

	applicationIds := make([]string, 0)
	if err = db.Select(&applicationIds, "SELECT application_id FROM pushtoken_application WHERE token_id = $1", pushToken.Token); err != nil {
		return nil, err
	}
	pushToken.AllowedApps = applicationIds

	return pushToken, nil
}

func UpdatePushtoken(id string, applications []string) (*PushToken, error) {
	if len(applications) == 0 {
		return nil, ErrApplicationlistEmpty
	}

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	for _, a := range applications {
		count := 0
		err = db.Get(&count, "SELECT COUNT(*) FROM application WHERE id = $1", a)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, ErrApplicationNotFound
		}
	}

	type tempToken struct {
		Id    string `db:"id"`
		Token string `db:"token"`
	}
	tToken := new(tempToken)

	if err = db.Get(tToken, "SELECT id, token FROM pushtoken where id = $1", id); err != nil {
		return nil, ErrPushtokenNotFound
	}

	result, err := db.Exec("DELETE FROM pushtoken_application WHERE token_id = $1", tToken.Token)
	if err != nil {
		return nil, ErrPushtokenNotFound
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, ErrPushtokenNotFound
	}

	if affected == 0 {
		return nil, ErrPushtokenNotFound
	}

	for _, a := range applications {
		_, err = db.Exec("INSERT INTO pushtoken_application (token_id, application_id) VALUES ($1, $2)", tToken.Token, a)

		if err != nil {
			return nil, err
		}
	}

	pushToken := new(PushToken)
	pushToken.Id = tToken.Id
	pushToken.Token = tToken.Token
	pushToken.AllowedApps = applications

	return pushToken, nil
}

func DeletePushtoken(id string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()
	var token string

	if err = db.Get(&token, "SELECT  token FROM pushtoken where id = $1", id); err != nil {
		return ErrPushtokenNotFound
	}

	result, err := db.Exec("DELETE FROM pushtoken_application WHERE token_id = $1", token)
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

	affected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
