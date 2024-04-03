package models

import (
	"database/sql"
	"jinya-releases/config"
	"jinya-releases/database"
	migrator "jinya-releases/database/migrations"
	"jinya-releases/test"
	"log"
	"os"
	"testing"
)

func cleanDatabase() error {
	conn, err := database.Connect()
	if err != nil {
		return err
	}

	_, _ = conn.Exec("DROP TABLE application CASCADE")
	_, _ = conn.Exec("DROP TABLE migrations CASCADE")

	return nil
}

func TestMain(m *testing.M) {
	if err := config.LoadConfiguration(); err != nil {
		log.Fatalln("Failed to load configuration")
	}

	if err := migrator.Migrate(); err != nil {
		log.Fatalln("Failed to migrate database")
	}

	code := m.Run()
	if err := cleanDatabase(); err != nil {
		log.Printf("Failed to clean database %v", err)
	}

	os.Exit(code)
}

func TestCreateApplication(t *testing.T) {
	type args struct {
		additionalApplications []Application
		application            Application
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateNewApplicationSuccessfulAllFields",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
					AdditionalCss: sql.NullString{
						String: "test",
						Valid:  true,
					},
					AdditionalJavaScript: sql.NullString{
						String: "test",
						Valid:  true,
					},
					Logo: sql.NullString{
						String: "test",
						Valid:  true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "CreateNewApplicationSuccessfulRequiredFields",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					TrackpageTemplate: "test",
					HomepageTemplate:  "test",
				},
			},
			wantErr: false,
		},
		{
			name: "CreateNewApplicationNameMissing",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "CreateNewApplicationSlugMissing",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "CreateNewApplicationHomepageTemplateMissing",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "CreateNewApplicationTrackpageTemplateMissing",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "CreateNewApplicationNonUniqueName",
			args: args{
				additionalApplications: []Application{
					{
						Name:              "test",
						Slug:              "test",
						HomepageTemplate:  "test",
						TrackpageTemplate: "test",
					},
				},
				application: Application{
					Name:              "test",
					Slug:              "test2",
					HomepageTemplate:  "test2",
					TrackpageTemplate: "test2",
				},
			},
			wantErr: true,
		},
		{
			name: "CreateNewApplicationNonUniqueSlug",
			args: args{
				additionalApplications: []Application{
					{
						Name:              "test",
						Slug:              "test",
						HomepageTemplate:  "test",
						TrackpageTemplate: "test",
					},
				},
				application: Application{
					Name:              "test2",
					Slug:              "test",
					HomepageTemplate:  "test2",
					TrackpageTemplate: "test2",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, app := range tt.args.additionalApplications {
				if _, err := CreateApplication(app); err != nil {
					t.Errorf("Prepare CreateApplication() error = %v", err)
					test.CleanTables()
					return
				}
			}

			got, err := CreateApplication(tt.args.application)
			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Id == "" {
				t.Errorf("CreateApplication() id is empty string")
			}
		})
	}
}
