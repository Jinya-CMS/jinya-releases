package models

import (
	"jinya-releases/test"
	"reflect"
	"testing"
	"time"
)

func TestCreateVersion(t *testing.T) {
	type args struct {
		version  Version
		versions []Version
	}
	tests := []struct {
		name    string
		args    args
		want    *Version
		wantErr bool
	}{
		{
			name: "CreateVersionSuccessfully",
			args: args{
				version: Version{
					Version:    "1.0",
					Url:        "test_url",
					UploadDate: time.Now(),
				}},
			wantErr: false,
		},
		{
			name: "CreateVersionVersionMissing",
			args: args{
				version: Version{
					Url:        "test_url",
					UploadDate: time.Now(),
				}},
			wantErr: true,
		},
		{
			name: "CreateVersionWrongApplicationId",
			args: args{
				version: Version{
					ApplicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Version:       "1.0",
					Url:           "test_url",
					UploadDate:    time.Now(),
				}},
			wantErr: true,
		},
		{
			name: "CreateVersionWrongTrackId",
			args: args{
				version: Version{
					TrackId:    "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Version:    "1.0",
					Url:        "test_url",
					UploadDate: time.Now(),
				}},
			wantErr: true,
		},
		{
			name: "CreateVersionVersionExists",
			args: args{
				version: Version{
					Version:    "1.0",
					Url:        "test_url",
					UploadDate: time.Now(),
				},
				versions: []Version{
					{
						Version:    "1.0",
						Url:        "test_url",
						UploadDate: time.Now(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
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

			track, err := CreateTrack(Track{
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

			if len(tt.args.version.ApplicationId) == 0 {
				tt.args.version.ApplicationId = app.Id
			}

			if len(tt.args.version.TrackId) == 0 {
				tt.args.version.TrackId = track.Id
			}

			for _, version := range tt.args.versions {
				version.ApplicationId = app.Id
				version.TrackId = track.Id
				if _, err := CreateVersion(version); err != nil {
					t.Errorf("Prepare create aditional version error = %v", err)
					test.CleanTables()
					return
				}
			}

			got, err := CreateVersion(tt.args.version)
			test.CleanTables()

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Id == "" {
				t.Errorf("CreateVersion() id is empty string")
			}
		})
	}
}

func TestDeleteVersionById(t *testing.T) {
	type args struct {
		applicationId string
		trackId       string
		id            string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "DeleteVersionByIdExists",
			args:    args{},
			wantErr: false,
		},
		{
			name: "DeleteVersionByIdApplicationDoesNotExists",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr: true,
		},
		{
			name: "DeleteVersionByIdTrackDoesNotExists",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr: true,
		},
		{
			name: "DeleteVersionByIdVersionDoesNotExists",
			args: args{
				id: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
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

			track, err := CreateTrack(Track{
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

			version, err := CreateVersion(Version{
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

			if len(tt.args.applicationId) == 0 {
				tt.args.applicationId = app.Id
			}

			if len(tt.args.trackId) == 0 {
				tt.args.trackId = track.Id
			}

			if len(tt.args.id) == 0 {
				tt.args.id = version.Id
			}

			if err := DeleteVersionById(tt.args.applicationId, tt.args.trackId, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteVersionById() error = %v, wantErr %v", err, tt.wantErr)
				test.CleanTables()
			}
			test.CleanTables()
		})
	}
}

func TestGetAllVersions(t *testing.T) {
	type args struct {
		applicationId string
		trackId       string
		versions      []Version
	}
	tests := []struct {
		name    string
		args    args
		want    []Version
		wantErr bool
	}{
		{
			name: "GetAllVersionsNone",
			args: args{
				versions: make([]Version, 0),
			},
			want:    make([]Version, 0),
			wantErr: false,
		},
		{
			name: "GetAllTracksOne",
			args: args{
				versions: []Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
			},
			want: []Version{
				{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllTracksMany",
			args: args{
				versions: []Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
					{
						Version:    "1.1",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
			},
			want: []Version{
				{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
				{
					Version:    "1.1",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllTracksWrongApplication",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				versions: []Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
			},
			want: []Version{
				{
					Version:    "1.0",
					Url:        "testurl",
					UploadDate: time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name: "GetAllTracksWrongTrack",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				versions: []Version{
					{
						Version:    "1.0",
						Url:        "testurl",
						UploadDate: time.Now(),
					},
				},
			},
			want: []Version{
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
			app, err := CreateApplication(Application{
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

			track, err := CreateTrack(Track{
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
				if _, err := CreateVersion(version); err != nil {
					t.Errorf("Prepare create aditional track error = %v", err)
					test.CleanTables()
					return
				}
			}

			var appId string
			if len(tt.args.applicationId) > 0 {
				appId = tt.args.applicationId
			} else {
				appId = app.Id
			}
			var trkId string
			if len(tt.args.trackId) > 0 {
				trkId = tt.args.trackId
			} else {
				trkId = track.Id
			}

			got, err := GetAllVersions(appId, trkId)
			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) && !tt.wantErr {
				t.Errorf("GetAllVersions() got = %v, want %v", got, tt.want)
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
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "GetVersionByIdExists",
			args:    args{},
			wantErr: false,
		},
		{
			name: "GetVersionByIdApplicationDoesNotExist",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr: true,
		},
		{
			name: "GetVersionByIdTrackDoesNotExist",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr: true,
		},
		{
			name: "GetVersionByIdVersionDoesNotExist",
			args: args{
				id: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
				Name:              "version_test",
				Slug:              "version_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare GetVersionById() error = %v", err)
				test.CleanTables()
				return
			}

			track, err := CreateTrack(Track{
				ApplicationId: app.Id,
				Name:          "version_test",
				Slug:          "version_test",
				IsDefault:     true,
			})
			if err != nil {
				t.Errorf("Prepare GetVersionById() error = %v", err)
				test.CleanTables()
				return
			}

			version, err := CreateVersion(Version{
				ApplicationId: app.Id,
				TrackId:       track.Id,
				Version:       "1.0",
				Url:           "test_url",
				UploadDate:    time.Now(),
			})
			if err != nil {
				t.Errorf("Prepare GetVersionById() error = %v", err)
				test.CleanTables()
				return
			}

			var appId string
			if len(tt.args.applicationId) > 0 {
				appId = tt.args.applicationId
			} else {
				appId = app.Id
			}
			var trkId string
			if len(tt.args.trackId) > 0 {
				trkId = tt.args.trackId
			} else {
				trkId = track.Id
			}
			var versId string
			if len(tt.args.id) > 0 {
				versId = tt.args.id
			} else {
				versId = version.Id
			}

			got, err := GetVersionById(appId, trkId, versId)
			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersionById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, version) && !tt.wantErr {
				t.Errorf("GetVersionById() got = %v, want %v", got, version)
			}
		})
	}
}

func TestGetVersionByTrackAndVersion(t *testing.T) {
	type args struct {
		trackId       string
		versionString string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "GetVersionByTrackAndVersionExists",
			args:    args{},
			wantErr: false,
		},
		{
			name: "GetVersionByTrackAndVersionTrackDoesNotExist",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr: true,
		},
		{
			name: "GetVersionByTrackAndVersionVersionDoesNotExist",
			args: args{
				versionString: "2.0",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
				Name:              "version_test",
				Slug:              "version_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare GetVersionByTrackAndVersion() error = %v", err)
				test.CleanTables()
				return
			}

			track, err := CreateTrack(Track{
				ApplicationId: app.Id,
				Name:          "version_test",
				Slug:          "version_test",
				IsDefault:     true,
			})
			if err != nil {
				t.Errorf("Prepare GetVersionByTrackAndVersion() error = %v", err)
				test.CleanTables()
				return
			}

			version, err := CreateVersion(Version{
				ApplicationId: app.Id,
				TrackId:       track.Id,
				Version:       "1.0",
				Url:           "test_url",
				UploadDate:    time.Now(),
			})
			if err != nil {
				t.Errorf("Prepare GetVersionByTrackAndVersion() error = %v", err)
				test.CleanTables()
				return
			}

			var trkId string
			if len(tt.args.trackId) > 0 {
				trkId = tt.args.trackId
			} else {
				trkId = track.Id
			}
			var vers string
			if len(tt.args.versionString) > 0 {
				vers = tt.args.versionString
			} else {
				vers = version.Version
			}
			got, err := GetVersionByTrackAndVersion(trkId, vers)
			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersionByTrackAndVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, version) && !tt.wantErr {
				t.Errorf("GetVersionByTrackAndVersion() got = %v, want %v", got, version)
			}
		})
	}
}
