package models

import (
	"jinya-releases/test"
	"reflect"
	"testing"
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
			name: "CreateVersionSuccess",
			args: args{
				version: Version{
					Version: "1.0 CreateVersionSuccess",
				}},
			wantErr: false,
		},
		{
			name: "CreateVersionVersionMissing",
			args: args{
				version: Version{}},
			wantErr: true,
		},
		{
			name: "CreateVersionWrongApplicationId",
			args: args{
				version: Version{
					ApplicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Version:       "1.0 CreateVersionWrongApplicationId",
				}},
			wantErr: true,
		},
		{
			name: "CreateVersionWrongTrackId",
			args: args{
				version: Version{
					TrackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Version: "1.0 CreateVersionWrongTrackId",
				}},
			wantErr: true,
		},
		{
			name: "CreateVersionVersionExists",
			args: args{
				version: Version{
					Version: "1.0 CreateVersionVersionExists",
				},
				versions: []Version{
					{
						Version: "1.0 CreateVersionVersionExists",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
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

			track, err := CreateTrack(Track{
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
			name:    "DeleteVersionByIdSuccess",
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
				Name:              "testDeleteVersionByIdApplication",
				Slug:              "testDeleteVersionByIdApplication",
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
				Name:          "testDeleteVersionByIdTrack",
				Slug:          "testDeleteVersionByIdTrack",
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
				Version:       "1.0 DeleteVersionById",
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
			name: "GetAllVersionsOne",
			args: args{
				versions: []Version{
					{
						Version: "1.0 GetAllVersionsOne",
					},
				},
			},
			want: []Version{
				{
					Version: "1.0 GetAllVersionsOne",
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllVersionsMany",
			args: args{
				versions: []Version{
					{
						Version: "1.0 GetAllVersionsMany",
					},
					{
						Version: "1.1 GetAllVersionsMany",
					},
				},
			},
			want: []Version{
				{
					Version: "1.0 GetAllVersionsMany",
				},
				{
					Version: "1.1 GetAllVersionsMany",
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllVersionsWrongApplication",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				versions: []Version{
					{
						Version: "1.0 GetAllVersionsWrongApplication",
					},
				},
			},
			want: []Version{
				{
					Version: "1.0 GetAllVersionsWrongApplication",
				},
			},
			wantErr: true,
		},
		{
			name: "GetAllVersionsWrongTrack",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				versions: []Version{
					{
						Version: "1.0 GetAllVersionsWrongTrack",
					},
				},
			},
			want: []Version{
				{
					Version: "1.0 GetAllVersionsWrongTrack",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
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

			track, err := CreateTrack(Track{
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
			name:    "GetVersionByIdSuccess",
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
				Name:              "testGetVersionByIdApplication",
				Slug:              "testGetVersionByIdApplication",
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
				Name:          "testGetVersionByIdTrack",
				Slug:          "testGetVersionByIdTrack",
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
				Version:       "1.0 testGetVersionByIdVersion",
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
			name:    "GetVersionByTrackAndVersionSuccess",
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
				Name:              "testGetVersionByTrackAndVersionApplication",
				Slug:              "testGetVersionByTrackAndVersionApplication",
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
				Name:          "testGetVersionByTrackAndVersionTrack",
				Slug:          "testGetVersionByTrackAndVersionTrack",
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
