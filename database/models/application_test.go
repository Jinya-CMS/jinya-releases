package models

import (
	"jinya-releases/config"
	"jinya-releases/database"
	migrator "jinya-releases/database/migrations"
	"jinya-releases/test"
	"log"
	"os"
	"reflect"
	"testing"
)

func cleanDatabase() error {
	conn, err := database.Connect()
	if err != nil {
		return err
	}

	_, _ = conn.Exec("ALTER TABLE track DROP CONSTRAINT track_applicationid_fkey")
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
	testPtr := "test"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateNewApplicationSuccessfulAllFields",
			args: args{
				application: Application{
					Name:                 "test",
					Slug:                 "test",
					HomepageTemplate:     "test",
					TrackpageTemplate:    "test",
					AdditionalCss:        &testPtr,
					AdditionalJavaScript: &testPtr,
					Logo:                 &testPtr,
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

func TestGetAllApplications(t *testing.T) {
	tests := []struct {
		name    string
		args    []Application
		want    []Application
		wantErr bool
	}{
		{
			name:    "GetAllApplicationsNone",
			args:    make([]Application, 0),
			want:    make([]Application, 0),
			wantErr: false,
		},
		{
			name: "GetAllApplicationsOne",
			args: []Application{
				{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			want: []Application{
				{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllApplicationsMany",
			args: []Application{
				{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				{
					Name:              "test1",
					Slug:              "test1",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
				},
			},
			want: []Application{
				{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				{
					Name:              "test1",
					Slug:              "test1",
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
				_, _ = CreateApplication(app)
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
		app Application
	}
	tests := []struct {
		name    string
		args    args
		want    *Application
		wantErr bool
	}{
		{
			name: "GetApplicationByIdPositiv",
			args: args{
				id: "",
				app: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "GetApplicationByIdNegative",
			args: args{
				id: "falseId",
				app: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			application, err := CreateApplication(tt.args.app)
			var got *Application
			if len(tt.args.id) > 0 {
				got, err = GetApplicationById(tt.args.id)
			} else {
				got, err = GetApplicationById(application.Id)
			}

			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetApplicationById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, application) && !tt.wantErr {
				t.Errorf("GetApplicationById() got = %v, want %v", got, application)
			}
		})
	}
}

func TestGetApplicationBySlug(t *testing.T) {
	type args struct {
		slug        string
		application Application
	}
	tests := []struct {
		name    string
		args    args
		want    *Application
		wantErr bool
	}{
		{
			name: "GetApplicationBySlugPositiv",
			args: args{
				slug: "test",
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "GetApplicationBySlugNegative",
			args: args{
				slug: "falseSlug",
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			application, err := CreateApplication(tt.args.application)
			if err != nil {
				t.Errorf("Prepare CreateApplication() error = %v", err)
				test.CleanTables()
				return
			}
			got, err := GetApplicationBySlug(tt.args.slug)

			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetApplicationById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, application) && !tt.wantErr {
				t.Errorf("GetApplicationById() got = %v, want %v", got, application)
			}
		})
	}
}

func TestUpdateApplication(t *testing.T) {
	type args struct {
		application            Application
		testApplication        Application
		additionalApplications []Application
	}
	tests := []struct {
		name    string
		args    args
		want    *Application
		wantErr bool
	}{
		{
			name: "UpdateApplicationFields",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				testApplication: Application{
					Name:              "test1",
					Slug:              "test1",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
				},
			},
			wantErr: false,
		},
		{
			name: "UpdateApplicationNameExists",
			args: args{
				additionalApplications: []Application{
					{
						Name:              "test2",
						Slug:              "test2",
						HomepageTemplate:  "test2",
						TrackpageTemplate: "test2",
					},
				},
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				testApplication: Application{
					Name:              "test2",
					Slug:              "test1",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
				},
			},
			wantErr: true,
		},
		{
			name: "UpdateApplicationSlugExists",
			args: args{
				additionalApplications: []Application{
					{
						Name:              "test2",
						Slug:              "test2",
						HomepageTemplate:  "test2",
						TrackpageTemplate: "test2",
					},
				},
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				testApplication: Application{
					Name:              "test1",
					Slug:              "test2",
					HomepageTemplate:  "test1",
					TrackpageTemplate: "test1",
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

			application, err := CreateApplication(tt.args.application)

			application.Name = tt.args.testApplication.Name
			application.Slug = tt.args.testApplication.Slug
			application.HomepageTemplate = tt.args.testApplication.HomepageTemplate
			application.TrackpageTemplate = tt.args.testApplication.TrackpageTemplate

			got, err := UpdateApplication(*application)
			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, application) && !tt.wantErr {
				t.Errorf("UpdateApplication() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteApplicationById(t *testing.T) {
	type args struct {
		id          string
		application Application
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "DeleteApplicationByIdExists",
			args: args{
				id: "",
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "DeleteApplicationByIdDoesNotExist",
			args: args{
				id: "falseId",
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CreateApplication(tt.args.application); err != nil {
				t.Errorf("Prepare CreateApplication() error = %v", err)
				test.CleanTables()
				return
			}

			got, err := GetApplicationBySlug(tt.args.application.Slug)
			if err != nil {
				t.Errorf("Prepare GetApplicationBySlug() error = %v", err)
				test.CleanTables()
				return
			}

			if len(tt.args.id) > 0 {
				if err := DeleteApplicationById(tt.args.id); (err != nil) != tt.wantErr {
					t.Errorf("DeleteApplicationById() error = %v, wantErr %v", err, tt.wantErr)
					test.CleanTables()
				}
			} else {
				if err := DeleteApplicationById(got.Id); (err != nil) != tt.wantErr {
					t.Errorf("DeleteApplicationById() error = %v, wantErr %v", err, tt.wantErr)
					test.CleanTables()
				}
			}
		})
	}
}
