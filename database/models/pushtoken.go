package models

import "jinya-releases/database"

type PushToken struct {
	Id          string   `json:"id" db:"id"`
	Token       string   `json:"token" db:"token"`
	AllowedApps []string `json:"allowedApps,omitempty"`
}

func CreatePushtoken(applications []Application) (*PushToken, error) {
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
		_, err = db.Exec("INSERT INTO pushtokenapplication (token, application) VALUES ($1, $2)", pushToken.Token, a.Id)

		if err != nil {
			return nil, err
		}
	}

	pushToken.AllowedApps = make([]string, len(applications))
	for i, application := range applications {
		pushToken.AllowedApps[i] = application.Id
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
		if err = db.Select(&pushTokens, "SELECT application FROM pushtokenapplication WHERE token = $1", pushtoken.Token); err != nil {
			return nil, err
		}
		pushTokens[i].AllowedApps = applicationIds

	}

	return pushTokens, nil
}
