package service

import (
	"bytes"
	"encoding/json"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/test"
	"jinya-releases/utils"
	"net/http"
	"reflect"
	"testing"
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
		wantErrDetails *utils.ErrorDetails
		wantStatus     int
		wantErr        bool
	}{
		{
			name: "CreateVersionSuccess",
			args: args{
				version: createVersionRequest{
					Version: "1.0 CreateVersionSuccess",
				}},
			wantErr:    false,
			wantStatus: http.StatusCreated,
		},
		{
			name: "CreateVersionVersionMissing",
			args: args{
				version: createVersionRequest{
					Version: "",
				}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreateVersionApplicationNotFound",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				version: createVersionRequest{
					Version: "1.0 CreateVersionApplicationNotFound",
				}},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "CreateVersionTrackNotFound",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				version: createVersionRequest{
					Version: "1.0 CreateVersionTrackNotFound",
				}},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "CreateVersionVersionExists",
			args: args{
				versions: []models.Version{
					{
						Version: "1.0 CreateVersionVersionExists",
					},
				},
				version: createVersionRequest{
					Version: "1.0 CreateVersionVersionExists",
				}},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "testCreateVersionApplication",
				Slug:              "testCreateVersionApplication",
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
				Name:          "testCreateVersionTrack",
				Slug:          "testCreateVersionTrack",
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
		wantErrDetails *utils.ErrorDetails
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
				Name:              "testDeleteVersionApplication",
				Slug:              "testDeleteVersionApplication",
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
				Name:          "testDeleteVersionTrack",
				Slug:          "testDeleteVersionTrack",
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
				Version:       "1.0 testDeleteVersionVersion",
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
		wantErrDetails *utils.ErrorDetails
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
						Version: "1.0 GetAllVersionsOne",
					},
				},
			},
			wantVersions: []models.Version{
				{
					Version: "1.0 GetAllVersionsOne",
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllVersionsMany",
			args: args{
				versions: []models.Version{
					{
						Version: "1.0 GetAllVersionsMany",
					},
					{
						Version: "2.0 GetAllVersionsMany",
					},
				},
			},
			wantVersions: []models.Version{
				{
					Version: "1.0 GetAllVersionsMany",
				},
				{
					Version: "2.0 GetAllVersionsMany",
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
						Version: "1.0 GetAllVersionsApplicationNotFound",
					},
				},
			},
			wantVersions: []models.Version{
				{
					Version: "1.0 GetAllVersionsApplicationNotFound",
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
						Version: "1.0 GetAllVersionsTrackNotFound",
					},
				},
			},
			wantVersions: []models.Version{
				{
					Version: "1.0 GetAllVersionsTrackNotFound",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "testGetAllVersionsApplication",
				Slug:              "testGetAllVersionsApplication",
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
				Name:          "testGetAllVersionsTrack",
				Slug:          "testGetAllVersionsTrack",
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
		wantErrDetails *utils.ErrorDetails
		wantStatus     int
		wantErr        bool
	}{
		{
			name:       "GetVersionByIdSuccess",
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
				Name:              "testGetVersionByIdApplication",
				Slug:              "testGetVersionByIdApplication",
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
				Name:          "testGetVersionByIdTrack",
				Slug:          "testGetVersionByIdTrack",
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
				Version:       "1.0 testGetVersionByIdVersion",
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
