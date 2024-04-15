package models

import (
	"jinya-releases/test"
	"reflect"
	"testing"
)

func TestCreateTrack(t *testing.T) {
	type args struct {
		application Application
		track       Track
		tracks      []Track
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateTrackSuccessfully",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
			},
			wantErr: false,
		},
		{
			name: "CreateTrackNameMissing",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Slug:      "test",
					IsDefault: true,
				},
			},
			wantErr: true,
		},
		{
			name: "CreateTrackSlugMissing",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Name:      "test",
					IsDefault: true,
				},
			},
			wantErr: true,
		},
		{
			name: "CreateTrackNameNotUnique",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				tracks: []Track{
					{
						Name:      "test",
						Slug:      "test1",
						IsDefault: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "CreateTrackSlugNotUnique",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				tracks: []Track{
					{
						ApplicationId: "",
						Name:          "test1",
						Slug:          "test",
						IsDefault:     true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "CreateTrackApplicationIdWrong",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					ApplicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Name:          "test",
					Slug:          "test",
					IsDefault:     true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(tt.args.application)
			if err != nil {
				t.Errorf("Prepare CreateTrack() error = %v", err)
				test.CleanTables()
				return
			}

			for _, track := range tt.args.tracks {
				track.ApplicationId = app.Id
				if _, err := CreateTrack(track); err != nil {
					t.Errorf("Prepare create aditional track error = %v", err)
					test.CleanTables()
					return
				}
			}

			if len(tt.args.track.ApplicationId) == 0 {
				tt.args.track.ApplicationId = app.Id
			}

			got, err := CreateTrack(tt.args.track)
			test.CleanTables()

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTrack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Id == "" {
				t.Errorf("CreateTrack() id is empty string")
			}
		})
	}
}

func TestGetAllTracks(t *testing.T) {
	type args struct {
		application Application
		tracks      []Track
	}
	tests := []struct {
		name    string
		args    args
		want    []Track
		wantErr bool
	}{
		{
			name: "GetAllTracksNone",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				tracks: make([]Track, 0),
			},
			want:    make([]Track, 0),
			wantErr: false,
		},
		{
			name: "GetAllTracksOne",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				tracks: []Track{
					{
						ApplicationId: "",
						Name:          "test",
						Slug:          "test",
						IsDefault:     true,
					},
				},
			},
			want: []Track{
				{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllTracksMany",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				tracks: []Track{
					{
						Name:      "test",
						Slug:      "test",
						IsDefault: true,
					},
					{
						Name:      "test1",
						Slug:      "test1",
						IsDefault: true,
					},
				},
			},
			want: []Track{
				{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				{
					Name:      "test1",
					Slug:      "test1",
					IsDefault: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(tt.args.application)
			if err != nil {
				t.Errorf("Prepare CreateTrack() error = %v", err)
				test.CleanTables()
				return
			}

			for _, track := range tt.args.tracks {
				track.ApplicationId = app.Id
				if _, err := CreateTrack(track); err != nil {
					t.Errorf("Prepare create aditional track error = %v", err)
					test.CleanTables()
					return
				}
			}

			got, err := GetAllTracks(app.Id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTracks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetAllTracks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTrackById(t *testing.T) {
	type args struct {
		application Application
		track       Track
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetTrackByIdPositive",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
			},
			wantErr: false,
		},
		{
			name: "GetTrackByIdWrongApplicationId",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					ApplicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Name:          "test",
					Slug:          "test",
					IsDefault:     true,
				},
			},
			wantErr: true,
		},
		{
			name: "GetTrackByIdWrongId",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Id:        "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(tt.args.application)
			if err != nil {
				t.Errorf("Prepare GetTrackById() error = %v", err)
				test.CleanTables()
				return
			}

			track, err := CreateTrack(Track{
				ApplicationId: app.Id,
				Name:          tt.args.track.Name,
				Slug:          tt.args.track.Slug,
				IsDefault:     false,
			})
			if err != nil {
				t.Errorf("Prepare GetTrackById() error = %v", err)
				test.CleanTables()
				return
			}

			if len(tt.args.track.Id) > 0 {
				track.Id = "e2ebb12e-e77d-4618-ba79-3f26e8af239a"
			}

			if len(tt.args.track.ApplicationId) > 0 {
				track.ApplicationId = "e2ebb12e-e77d-4618-ba79-3f26e8af239a"
			}

			got, err := GetTrackById(track.Id, track.ApplicationId)
			test.CleanTables()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrackById() error = %v, wantErr %v, got %v", err, tt.wantErr, got)
				return
			}
			if !reflect.DeepEqual(got, track) && !tt.wantErr {
				t.Errorf("GetTrackById() got = %v, want %v", got, track)
			}
		})
	}
}

func TestGetTrackBySlug(t *testing.T) {
	type args struct {
		application Application
		track       Track
		slug        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetTrackBySlugPositive",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				slug: "test",
			},
			wantErr: false,
		},
		{
			name: "GetTrackBySlugWrongSlug",
			args: args{
				application: Application{
					Name:              "test",
					Slug:              "test",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				track: Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				slug: "wrongSlug",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(tt.args.application)
			if err != nil {
				t.Errorf("Prepare GetTrackById() error = %v", err)
				test.CleanTables()
				return
			}

			tt.args.track.ApplicationId = app.Id
			track, err := CreateTrack(tt.args.track)
			if err != nil {
				t.Errorf("Prepare GetTrackById() error = %v", err)
				test.CleanTables()
				return
			}

			got, err := GetTrackBySlug(tt.args.slug, app.Id)
			test.CleanTables()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrackById() error = %v, wantErr %v, got %v", err, tt.wantErr, got)
				return
			}
			if !reflect.DeepEqual(got, track) && !tt.wantErr {
				t.Errorf("GetTrackById() got = %v, want %v", got, track)
			}
		})
	}
}
