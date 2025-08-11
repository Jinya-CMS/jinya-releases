package models

import (
	"errors"
	"jinya-releases/database"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
)

type PushToken struct {
	Id            string `json:"id" db:"id"`
	Token         string `json:"token" db:"token"`
	ApplicationId string `json:"-" db:"application_id"`
}

var (
	ErrPushtokenNotFound    = errors.New("token not found")
	ErrApplicationlistEmpty = errors.New("application list empty")
)

func CreateToken(application string) (*PushToken, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	token, err := shortid.Generate()
	if err != nil {
		token = uuid.New().String()
	}

	if _, err = db.Exec("INSERT INTO push_token (token, application_id) VALUES ($1, $2)", token, application); err != nil {
		return nil, err
	}

	pushToken := new(PushToken)
	if err = db.Get(pushToken, "SELECT id, token FROM push_token WHERE token = $1", token); err != nil {
		return nil, err
	}

	return pushToken, nil
}

func CheckPushToken(token, applicationId string) bool {
	db, err := database.Connect()
	if err != nil {
		return false
	}

	defer db.Close()
	type tokenCount struct {
		Count int `db:"count"`
	}
	count := new(tokenCount)

	if err = db.Get(count, `select count(*) from push_token pt
         join application a on a.id = pt.application_id
where
    pt.token = $1 and a.slug = $2`, token, applicationId); err != nil {
		return false
	}

	return count.Count > 0
}

func ResetToken(applicationId string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	applicationCount := 0

	if err = db.Get(&applicationCount, "SELECT COUNT(*) FROM application WHERE id = $1", applicationId); err != nil {
		return err
	}

	if applicationCount == 0 {
		return ErrApplicationNotFound
	}

	_, err = db.Exec("DELETE FROM push_token WHERE application_id = $1", applicationId)

	return err
}
