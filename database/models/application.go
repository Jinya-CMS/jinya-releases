package models

import (
	"errors"
	"jinya-releases/database"
)

type Application struct {
	Id                   string  `json:"id" db:"id"`
	Name                 string  `json:"name" db:"name"`
	Logo                 *string `json:"logo" db:"logo"`
	Slug                 string  `json:"slug" db:"slug"`
	HomepageTemplate     string  `json:"homepageTemplate" db:"homepage_template"`
	TrackpageTemplate    string  `json:"trackpageTemplate" db:"trackpage_template"`
	AdditionalCss        *string `json:"additionalCss,omitempty" db:"additional_css"`
	AdditionalJavaScript *string `json:"additionalJavaScript,omitempty" db:"additional_javascript"`
}

var (
	ErrNameEmpty              = errors.New("name is empty")
	ErrSlugEmpty              = errors.New("slug is empty")
	ErrHomepageTemplateEmpty  = errors.New("homepage template is empty")
	ErrTrackpageTemplateEmpty = errors.New("trackpage template is empty")
	ErrApplicationNotFound    = errors.New("application not found")
)

func CreateApplication(application Application) (*Application, error) {
	if application.Name == "" {
		return nil, ErrNameEmpty
	}
	if application.Slug == "" {
		return nil, ErrSlugEmpty
	}
	if application.HomepageTemplate == "" {
		return nil, ErrHomepageTemplateEmpty
	}
	if application.TrackpageTemplate == "" {
		return nil, ErrTrackpageTemplateEmpty
	}

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

func GetApplicationById(id string) (*Application, error) {
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

	result, err := db.Exec("UPDATE application SET name = $1, logo = $2, slug = $3, homepage_template = $4, trackpage_template = $5, additional_css = $6, additional_javascript = $7 WHERE id = $8", application.Name, application.Logo, application.Slug, application.HomepageTemplate, application.TrackpageTemplate, application.AdditionalCss, application.AdditionalJavaScript, application.Id)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return nil, ErrApplicationNotFound
	}

	return GetApplicationBySlug(application.Slug)
}

func DeleteApplicationById(id string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	result, err := db.Exec("DELETE FROM application WHERE id = $1", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrApplicationNotFound
	}

	return nil
}
