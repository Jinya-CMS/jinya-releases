package service

import (
	"bytes"
	"encoding/json"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/test"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestCreateVersion(t *testing.T) {
	type args struct {
		reader        io.Reader
		applicationId string
		trackId       string
		version       createVersionRequest
		versions      []models.Version
	}
	tests := []struct {
		name           string
		args           args
		wantVersion    *models.Version
		wantErrDetails *ErrorDetails
		wantStatus     int
		wantErr        bool
	}{
		{
			name: "CreateVersionPositive",
			args: args{
				version: createVersionRequest{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Time{},
				}},
			wantErr:    false,
			wantStatus: http.StatusCreated,
		},
		{
			name: "CreateVersionVersionMissing",
			args: args{
				version: createVersionRequest{
					Version:    "",
					Url:        "testurl",
					UploadDate: time.Now(),
				}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreateVersionApplicationNotFound",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				version: createVersionRequest{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				}},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "CreateVersionTrackNotFound",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				version: createVersionRequest{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				}},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "CreateVersionVersionMissing",
			args: args{
				versions: []models.Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
				version: createVersionRequest{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				}},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "version_test",
				Slug:              "version_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare CreateVersion() error = %v", err)
				test.CleanTables()
				return
			}

			track, err := models.CreateTrack(models.Track{
				ApplicationId: app.Id,
				Name:          "version_test",
				Slug:          "version_test",
				IsDefault:     true,
			})
			if err != nil {
				t.Errorf("Prepare CreateVersion() error = %v", err)
				test.CleanTables()
				return
			}

			for _, version := range tt.args.versions {
				version.ApplicationId = app.Id
				version.TrackId = track.Id
				if _, err := models.CreateVersion(version); err != nil {
					t.Errorf("Prepare CreateVersion() error = %v", err)
					test.CleanTables()
					return
				}
			}

			var buffer bytes.Buffer
			_ = json.NewEncoder(&buffer).Encode(tt.args.version)

			var applicationId string
			if len(tt.args.applicationId) > 0 {
				applicationId = tt.args.applicationId
			} else {
				applicationId = app.Id
			}
			var trackId string
			if len(tt.args.trackId) > 0 {
				trackId = tt.args.trackId
			} else {
				trackId = track.Id
			}

			_, gotErrDetails, gotStatus := CreateVersion(&buffer, applicationId, trackId)
			test.CleanTables()

			if (gotErrDetails != nil) != tt.wantErr {
				t.Errorf("CreateVersion() gotErrDetails = %v, want %v", gotErrDetails, tt.wantErr)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("CreateVersion() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestDeleteVersion(t *testing.T) {
	type args struct {
		applicationId string
		trackId       string
		id            string
	}
	tests := []struct {
		name           string
		args           args
		wantErrDetails *ErrorDetails
		wantStatus     int
		wantErr        bool
	}{
		{
			name:       "DeleteVersionSuccess",
			args:       args{},
			wantStatus: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name: "DeleteVersionApplicationNotFound",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "DeleteVersionTrackNotFound",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "DeleteVersionVersionNotFound",
			args: args{
				id: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "version_test",
				Slug:              "version_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare DeleteVersion() error = %v", err)
				test.CleanTables()
				return
			}

			track, err := models.CreateTrack(models.Track{
				ApplicationId: app.Id,
				Name:          "version_test",
				Slug:          "version_test",
				IsDefault:     true,
			})
			if err != nil {
				t.Errorf("Prepare DeleteVersion() error = %v", err)
				test.CleanTables()
				return
			}

			version, err := models.CreateVersion(models.Version{
				ApplicationId: app.Id,
				TrackId:       track.Id,
				Version:       "1.0",
				Url:           "test_url",
				UploadDate:    time.Now(),
			})
			if err != nil {
				t.Errorf("Prepare DeleteVersion() error = %v", err)
				test.CleanTables()
				return
			}

			var applicationId string
			if len(tt.args.applicationId) > 0 {
				applicationId = tt.args.applicationId
			} else {
				applicationId = app.Id
			}
			var trackId string
			if len(tt.args.trackId) > 0 {
				trackId = tt.args.trackId
			} else {
				trackId = track.Id
			}
			var versionId string
			if len(tt.args.id) > 0 {
				versionId = tt.args.id
			} else {
				versionId = version.Id
			}

			gotErrDetails, gotStatus := DeleteVersion(applicationId, trackId, versionId)
			test.CleanTables()

			if (gotErrDetails != nil) != tt.wantErr {
				t.Errorf("DeleteVersion() gotErrDetails = %v, want %v", gotErrDetails, tt.wantErr)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("DeleteVersion() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestGetAllVersions(t *testing.T) {
	type args struct {
		applicationId string
		trackId       string
		versions      []models.Version
	}
	tests := []struct {
		name           string
		args           args
		wantVersions   []models.Version
		wantErrDetails *ErrorDetails
		wantStatus     int
		wantErr        bool
	}{
		{
			name: "GetAllVersionsNone",
			args: args{
				versions: make([]models.Version, 0),
			},
			wantVersions: make([]models.Version, 0),
			wantErr:      false,
		},
		{
			name: "GetAllVersionsOne",
			args: args{
				versions: []models.Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
			},
			wantVersions: []models.Version{
				{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllVersionsMany",
			args: args{
				versions: []models.Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
					{
						Version:    "2.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
			},
			wantVersions: []models.Version{
				{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
				{
					Version:    "2.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllVersionsApplicationNotFound",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				versions: []models.Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
			},
			wantVersions: []models.Version{
				{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name: "GetAllVersionsTrackNotFound",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				versions: []models.Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
			},
			wantVersions: []models.Version{
				{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "version_test",
				Slug:              "version_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare GetAllVersions() error = %v", err)
				test.CleanTables()
				return
			}

			track, err := models.CreateTrack(models.Track{
				ApplicationId: app.Id,
				Name:          "version_test",
				Slug:          "version_test",
				IsDefault:     true,
			})
			if err != nil {
				t.Errorf("Prepare GetAllVersions() error = %v", err)
				test.CleanTables()
				return
			}
			for _, version := range tt.args.versions {
				version.ApplicationId = app.Id
				version.TrackId = track.Id
				_, err = models.CreateVersion(version)
			}
			if err != nil {
				t.Errorf("Prepare GetAllVersions() error = %v", err)
				test.CleanTables()
				return
			}

			var applicationId string
			if len(tt.args.applicationId) > 0 {
				applicationId = tt.args.applicationId
			} else {
				applicationId = app.Id
			}
			var trackId string
			if len(tt.args.trackId) > 0 {
				trackId = tt.args.trackId
			} else {
				trackId = track.Id
			}

			gotVersions, errDetails, _ := GetAllVersions(applicationId, trackId)
			test.CleanTables()

			if (errDetails != nil) != tt.wantErr {
				t.Errorf("GetAllTracks() error = %v, wantErr %v", errDetails, tt.wantErr)
				return
			}
			if (len(gotVersions) != len(tt.wantVersions)) != tt.wantErr {
				t.Errorf("GetAllTracks() gotTracks = %v, want %v", gotVersions, tt.wantVersions)
				return
			}
		})
	}
}

func TestGetVersionById(t *testing.T) {
	type args struct {
		applicationId string
		trackId       string
		id            string
	}
	tests := []struct {
		name           string
		args           args
		wantVersion    *models.Version
		wantErrDetails *ErrorDetails
		wantStatus     int
		wantErr        bool
	}{
		{
			name:       "GetVersionByIdPositive",
			args:       args{},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "GetVersionByIdApplicationNotFound",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "GetVersionByIdTrackNotFound",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "GetVersionByIdVersionNotFound",
			args: args{
				id: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "version_test",
				Slug:              "version_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare DeleteVersion() error = %v", err)
				test.CleanTables()
				return
			}

			track, err := models.CreateTrack(models.Track{
				ApplicationId: app.Id,
				Name:          "version_test",
				Slug:          "version_test",
				IsDefault:     true,
			})
			if err != nil {
				t.Errorf("Prepare DeleteVersion() error = %v", err)
				test.CleanTables()
				return
			}

			version, err := models.CreateVersion(models.Version{
				ApplicationId: app.Id,
				TrackId:       track.Id,
				Version:       "1.0",
				Url:           "test_url",
				UploadDate:    time.Now(),
			})
			if err != nil {
				t.Errorf("Prepare DeleteVersion() error = %v", err)
				test.CleanTables()
				return
			}

			var applicationId string
			if len(tt.args.applicationId) > 0 {
				applicationId = tt.args.applicationId
			} else {
				applicationId = app.Id
			}
			var trackId string
			if len(tt.args.trackId) > 0 {
				trackId = tt.args.trackId
			} else {
				trackId = track.Id
			}
			var versionId string
			if len(tt.args.id) > 0 {
				versionId = tt.args.id
			} else {
				versionId = version.Id
			}
			gotVersion, gotErrDetails, gotStatus := GetVersionById(applicationId, trackId, versionId)
			test.CleanTables()

			if !reflect.DeepEqual(gotVersion, version) != tt.wantErr {
				t.Errorf("GetVersionById() gotVersion = %v, want %v", gotVersion, version)
			}
			if (gotErrDetails != nil) != tt.wantErr {
				t.Errorf("GetVersionById() gotErrDetails = %v, want %v", gotErrDetails, tt.wantErr)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetVersionById() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}
