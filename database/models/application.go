package models

import (
	"database/sql"
	"jinya-releases/database"
)

type Application struct {
	Id                   string         `json:"id"`
	Name                 string         `json:"name"`
	Logo                 string         `json:"logo"`
	Slug                 string         `json:"slug"`
	HomepageTemplate     string         `json:"homepageTemplate"`
	TrackpageTemplate    string         `json:"trackpageTemplate"`
	AdditionalCss        sql.NullString `json:"additionalCss,omitempty"`
	AdditionalJavaScript sql.NullString `json:"additionalJavaScript,omitempty"`
}

func CreateApplication(application Application) (*Application, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO application (name, logo, slug, homepage_template, trackpage_template, additional_css, additional_javascript) VALUES ($1, $2, $3, $4, $5, $6, $7)", application.Name, application.Logo, application.Slug, application.HomepageTemplate, application.TrackpageTemplate, application.AdditionalCss, application.AdditionalJavaScript)

	if err != nil {
		return nil, err
	}

	return GetApplicationBySlug(application.Slug)
}

func GetAllApplications() ([]Application, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	applications := make([]Application, 0)

	if err = db.Select(&applications, "SELECT id, name, logo, slug, homepage_template, trackpage_template, additional_css, additional_javascript FROM application ORDER BY name"); err != nil {
		return nil, err
	}

	return applications, nil
}

func GetApplicationById(id int) (*Application, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	application := new(Application)

	if err = db.Get(application, "SELECT id, name, logo, slug, homepage_template, trackpage_template, additional_css, additional_javascript FROM application WHERE id = $1", id); err != nil {
		return nil, err
	}

	return application, nil
}

func GetApplicationBySlug(slug string) (*Application, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	application := new(Application)

	if err = db.Get(application, "SELECT id, name, logo, slug, homepage_template, trackpage_template, additional_css, additional_javascript FROM application WHERE slug = $1", slug); err != nil {
		return nil, err
	}

	return application, nil
}

func UpdateApplication(application Application) (*Application, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE application SET name = $1, logo = $2, slug = $3, homepage_template = $4, trackpage_template = $5, additional_css = $6, additional_javascript = $7", application.Name, application.Logo, application.Slug, application.HomepageTemplate, application.TrackpageTemplate, application.AdditionalCss, application.AdditionalJavaScript)

	if err != nil {
		return nil, err
	}

	return GetApplicationBySlug(application.Slug)
}

func DeleteApplicationById(id int) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM application WHERE id = $1", id)
	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}
