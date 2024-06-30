package service

import (
	"bytes"
	"encoding/json"
	"jinya-releases/config"
	"jinya-releases/database"
	migrator "jinya-releases/database/migrations"
	"jinya-releases/database/models"
	"jinya-releases/test"
	"jinya-releases/utils"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func cleanDatabase() error {
	conn, err := database.Connect()
	if err != nil {
		return err
	}

	_, _ = conn.Exec("DROP TABLE pushtokenapplication CASCADE")
	_, _ = conn.Exec("DROP TABLE pushtoken CASCADE")
	_, _ = conn.Exec("DROP TABLE version CASCADE")
	_, _ = conn.Exec("DROP TABLE track CASCADE")
	_, _ = conn.Exec("DROP TABLE application CASCADE")
	_, _ = conn.Exec("DROP TABLE migrations CASCADE")

	return nil
}

func TestMain(m *testing.M) {
	if err := config.LoadConfiguration(); err != nil {
		log.Fatalln("Failed to load configuration")
	}

	if err := migrator.Migrate(); err != nil {
		log.Println("Failed to migrate database")
		log.Println("Running tests anyway, database might be migrated already")
	}

	code := m.Run()
	if err := cleanDatabase(); err != nil {
		log.Printf("Failed to clean database %v", err)
	}

	os.Exit(code)
}

func TestGetAllApplications(t *testing.T) {
	tests := []struct {
		name    string
		args    []models.Application
		want    []models.Application
		wantErr bool
	}{
		{
			name:    "GetAllApplicationsNone",
			args:    make([]models.Application, 0),
			want:    make([]models.Application, 0),
			wantErr: false,
		},
		{
			name: "GetAllApplicationsOne",
			args: []models.Application{
				{
					Name:              "testGetAllApplicationsOne",
					Slug:              "testGetAllApplicationsOne",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			want: []models.Application{
				{
					Name:              "testGetAllApplicationsOne",
					Slug:              "testGetAllApplicationsOne",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllApplicationsMany",
			args: []models.Application{
				{
					Name:              "testGetAllApplicationsMany",
					Slug:              "testGetAllApplicationsMany",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				{
					Name:              "testGetAllApplicationsMany1",
					Slug:              "testGetAllApplicationsMany1",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
				},
			},
			want: []models.Application{
				{
					Name:              "testGetAllApplicationsMany",
					Slug:              "testGetAllApplicationsMany",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				{
					Name:              "testGetAllApplicationsMany1",
					Slug:              "testGetAllApplicationsMany1",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, app := range tt.args {
				_, _ = models.CreateApplication(app)
			}
			got, err := GetAllApplications()
			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllApplications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetAllApplications() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetApplicationById(t *testing.T) {
	type args struct {
		id  string
		app models.Application
	}
	tests := []struct {
		name       string
		args       args
		want       *models.Application
		wantErr    bool
		wantStatus int
	}{
		{
			name: "GetApplicationByIdSuccess",
			args: args{
				id: "",
				app: models.Application{
					Name:              "testGetApplicationByIdSuccess",
					Slug:              "testGetApplicationByIdSuccess",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "GetApplicationByIdFalseId",
			args: args{
				id: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				app: models.Application{
					Name:              "GetApplicationByIdFalseId",
					Slug:              "GetApplicationByIdFalseId",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			application, _ := models.CreateApplication(tt.args.app)
			var (
				got        *models.Application
				status     int
				errDetails *utils.ErrorDetails
			)
			if len(tt.args.id) > 0 {
				got, errDetails, status = GetApplicationById(tt.args.id)
			} else {
				got, errDetails, status = GetApplicationById(application.Id)
			}

			test.CleanTables()
			if (errDetails != nil) != tt.wantErr {
				t.Errorf("GetApplicationById() error = %v, wantErr %v", errDetails, tt.wantErr)
				return
			}
			if tt.wantStatus != status {
				t.Errorf("GetApplicationById() status = %v, wantStatus %v", status, tt.wantStatus)
				return
			}
			if !reflect.DeepEqual(got, application) && !tt.wantErr {
				t.Errorf("GetApplicationById() got = %v, want %v", got, application)
			}
		})
	}
}

func TestCreateApplication(t *testing.T) {
	type args struct {
		additionalApplications []models.Application
		application            createApplicationRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "CreateNewApplicationSuccessfulAllFields",
			args: args{
				application: createApplicationRequest{
					Name:                 "testCreateNewApplicationSuccessfulAllFields",
					Slug:                 "testCreateNewApplicationSuccessfulAllFields",
					HomepageTemplate:     "test",
					TrackpageTemplate:    "test",
					AdditionalCss:        "test",
					AdditionalJavaScript: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "CreateNewApplicationSuccessfulRequiredFields",
			args: args{
				application: createApplicationRequest{
					Name:              "testCreateNewApplicationSuccessfulRequiredFields",
					Slug:              "testCreateNewApplicationSuccessfulRequiredFields",
					TrackpageTemplate: "test",
					HomepageTemplate:  "test",
				},
			},
			wantErr: false,
		},
		{
			name: "CreateNewApplicationNameMissing",
			args: args{
				application: createApplicationRequest{
					Name:              "testCreateNewApplicationNameMissing",
					Slug:              "testCreateNewApplicationNameMissing",
					TrackpageTemplate: "test",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreateNewApplicationSlugMissing",
			args: args{
				application: createApplicationRequest{
					Name:              "testCreateNewApplicationSlugMissing",
					Slug:              "testCreateNewApplicationSlugMissing",
					TrackpageTemplate: "test",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreateNewApplicationHomepageTemplateMissing",
			args: args{
				application: createApplicationRequest{
					Name:              "testCreateNewApplicationHomepageTemplateMissing",
					Slug:              "testCreateNewApplicationHomepageTemplateMissing",
					TrackpageTemplate: "test",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreateNewApplicationTrackpageTemplateMissing",
			args: args{
				application: createApplicationRequest{
					Name:              "testCreateNewApplicationTrackpageTemplateMissing",
					Slug:              "testCreateNewApplicationTrackpageTemplateMissing",
					TrackpageTemplate: "test",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreateNewApplicationNonUniqueName",
			args: args{
				additionalApplications: []models.Application{
					{
						Name:              "testCreateNewApplicationNonUniqueName",
						Slug:              "testCreateNewApplicationNonUniqueName",
						HomepageTemplate:  "test",
						TrackpageTemplate: "test",
					},
				},
				application: createApplicationRequest{
					Name:              "testCreateNewApplicationNonUniqueName",
					Slug:              "testCreateNewApplicationNonUniqueName2",
					HomepageTemplate:  "test2",
					TrackpageTemplate: "test2",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
		{
			name: "CreateNewApplicationNonUniqueSlug",
			args: args{
				additionalApplications: []models.Application{
					{
						Name:              "testCreateNewApplicationNonUniqueSlug",
						Slug:              "testCreateNewApplicationNonUniqueSlug",
						HomepageTemplate:  "test",
						TrackpageTemplate: "test",
					},
				},
				application: createApplicationRequest{
					Name:              "testCreateNewApplicationNonUniqueSlug2",
					Slug:              "testCreateNewApplicationNonUniqueSlug",
					HomepageTemplate:  "test2",
					TrackpageTemplate: "test2",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, app := range tt.args.additionalApplications {
				if _, err := models.CreateApplication(app); err != nil {
					t.Errorf("Prepare CreateApplication() error = %v", err)
					test.CleanTables()
					return
				}
			}

			var buffer bytes.Buffer
			_ = json.NewEncoder(&buffer).Encode(tt.args.application)

			got, errDetails, status := CreateApplication(&buffer)
			test.CleanTables()
			if (errDetails != nil) != tt.wantErr {
				t.Errorf("CreateApplication() error = %v, wantErr %v", errDetails, tt.wantErr)
				return
			}
			if tt.wantStatus != status {
				t.Errorf("CreateApplication() status = %v, wantStatus %v", status, tt.wantStatus)
				return
			}
			if got != nil && got.Id == "" {
				t.Errorf("CreateApplication() id is empty string")
			}
		})
	}
}

func TestUpdateApplication(t *testing.T) {
	type args struct {
		additionalApplications []models.Application
		application            models.Application
		id                     string
		testApplication        updateApplicationRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "UpdateApplicationFields",
			args: args{
				application: models.Application{
					Name:              "testUpdateApplicationFields",
					Slug:              "testUpdateApplicationFields",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				testApplication: updateApplicationRequest{
					Name:              "testUpdateApplicationFields1",
					Slug:              "testUpdateApplicationFields1",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
				},
			},
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
		{
			name: "UpdateApplicationNameExists",
			args: args{
				additionalApplications: []models.Application{
					{
						Name:              "testUpdateApplicationNameExists2",
						Slug:              "testUpdateApplicationNameExists2",
						HomepageTemplate:  "test2",
						TrackpageTemplate: "test2",
					},
				},
				application: models.Application{
					Name:              "testUpdateApplicationNameExists",
					Slug:              "testUpdateApplicationNameExists",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				testApplication: updateApplicationRequest{
					Name:              "testUpdateApplicationNameExists2",
					Slug:              "testUpdateApplicationNameExists1",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
		{
			name: "UpdateApplicationSlugExists",
			args: args{
				additionalApplications: []models.Application{
					{
						Name:              "testUpdateApplicationSlugExists2",
						Slug:              "testUpdateApplicationSlugExists2",
						HomepageTemplate:  "test2",
						TrackpageTemplate: "test2",
					},
				},
				application: models.Application{
					Name:              "testUpdateApplicationSlugExists",
					Slug:              "testUpdateApplicationSlugExists",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				testApplication: updateApplicationRequest{
					Name:              "testUpdateApplicationSlugExists1",
					Slug:              "testUpdateApplicationSlugExists2",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, app := range tt.args.additionalApplications {
				if _, err := models.CreateApplication(app); err != nil {
					t.Errorf("Prepare UpdateApplication() error = %v", err)
					test.CleanTables()
					return
				}
			}

			app, err := models.CreateApplication(tt.args.application)
			if err != nil {
				t.Errorf("Prepare UpdateApplication() error = %v", err)
				test.CleanTables()
				return
			}

			var buffer bytes.Buffer
			_ = json.NewEncoder(&buffer).Encode(tt.args.testApplication)
			got, errDetails, status := UpdateApplication(app.Id, &buffer)
			test.CleanTables()
			if (errDetails != nil) != tt.wantErr {
				t.Errorf("UpdateApplication() error = %v, wantErr %v", errDetails, tt.wantErr)
				return
			}
			if tt.wantStatus != status {
				t.Errorf("UpdateApplication() status = %v, wantStatus %v", status, tt.wantStatus)
				return
			}
			if got != nil && got.Id == "" {
				t.Errorf("UpdateApplication() id is empty string")
			}
		})
	}
}

func TestDeleteApplication(t *testing.T) {
	type args struct {
		id          string
		application models.Application
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "DeleteApplicationByIdSuccess",
			args: args{
				id: "",
				application: models.Application{
					Name:              "testDeleteApplicationByIdSuccess",
					Slug:              "testDeleteApplicationByIdSuccess",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
		{
			name: "DeleteApplicationByIdDoesNotExist",
			args: args{
				id: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				application: models.Application{
					Name:              "testDeleteApplicationByIdDoesNotExist",
					Slug:              "testDeleteApplicationByIdDoesNotExist",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := models.CreateApplication(tt.args.application); err != nil {
				t.Errorf("Prepare CreateApplication() error = %v", err)
				test.CleanTables()
				return
			}

			got, err := models.GetApplicationBySlug(tt.args.application.Slug)
			if err != nil {
				t.Errorf("Prepare GetApplicationBySlug() error = %v", err)
				test.CleanTables()
				return
			}

			if len(tt.args.id) > 0 {
				err, status := DeleteApplication(tt.args.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("DeleteApplication() error = %v, wantErr %v", err, tt.wantErr)
					test.CleanTables()
					return
				} else if tt.wantStatus != status {
					t.Errorf("DeleteApplication() status = %v, wantStatus %v", status, tt.wantStatus)
					return
				}
			} else {
				err, status := DeleteApplication(got.Id)
				if (err != nil) != tt.wantErr {
					t.Errorf("DeleteApplication() error = %v, wantErr %v", err, tt.wantErr)
					test.CleanTables()
					return
				} else if tt.wantStatus != status {
					t.Errorf("DeleteApplication() status = %v, wantStatus %v", status, tt.wantStatus)
					return
				}
			}

		})
	}
}
