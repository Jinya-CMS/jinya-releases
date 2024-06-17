package models

import (
	"jinya-releases/test"
	"reflect"
	"testing"
)

func TestCreateTrack(t *testing.T) {
	type args struct {
		track  Track
		tracks []Track
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateTrackSuccess",
			args: args{
				track: Track{
					Name:      "testCreateTrackSuccess",
					Slug:      "testCreateTrackSuccess",
					IsDefault: true,
				},
			},
			wantErr: false,
		},
		{
			name: "CreateTrackNameMissing",
			args: args{
				track: Track{
					Slug:      "testCreateTrackNameMissing",
					IsDefault: true,
				},
			},
			wantErr: true,
		},
		{
			name: "CreateTrackSlugMissing",
			args: args{
				track: Track{
					Name:      "testCreateTrackSlugMissing",
					IsDefault: true,
				},
			},
			wantErr: true,
		},
		{
			name: "CreateTrackNameNotUnique",
			args: args{
				track: Track{
					Name:      "testCreateTrackNameNotUnique",
					Slug:      "testCreateTrackNameNotUnique",
					IsDefault: true,
				},
				tracks: []Track{
					{
						Name:      "testCreateTrackNameNotUnique",
						Slug:      "testCreateTrackNameNotUnique1",
						IsDefault: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "CreateTrackSlugNotUnique",
			args: args{
				track: Track{
					Name:      "testCreateTrackSlugNotUnique",
					Slug:      "testCreateTrackSlugNotUnique",
					IsDefault: true,
				},
				tracks: []Track{
					{
						ApplicationId: "",
						Name:          "testCreateTrackSlugNotUnique1",
						Slug:          "testCreateTrackSlugNotUnique",
						IsDefault:     true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "CreateTrackApplicationIdWrong",
			args: args{
				track: Track{
					ApplicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Name:          "testCreateTrackApplicationIdWrong",
					Slug:          "testCreateTrackApplicationIdWrong",
					IsDefault:     true,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
				Name:              "testCreateTrackApplication",
				Slug:              "testCreateTrackApplication",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
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
		tracks []Track
		appId  string
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
				tracks: make([]Track, 0),
			},
			want:    make([]Track, 0),
			wantErr: false,
		},
		{
			name: "GetAllTracksOne",
			args: args{
				tracks: []Track{
					{
						ApplicationId: "",
						Name:          "testGetAllTracksOne",
						Slug:          "testGetAllTracksOne",
						IsDefault:     true,
					},
				},
			},
			want: []Track{
				{
					Name:      "testGetAllTracksOne",
					Slug:      "testGetAllTracksOne",
					IsDefault: true,
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllTracksMany",
			args: args{
				tracks: []Track{
					{
						Name:      "testGetAllTracksMany",
						Slug:      "testGetAllTracksMany",
						IsDefault: true,
					},
					{
						Name:      "testGetAllTracksMany1",
						Slug:      "testGetAllTracksMany1",
						IsDefault: true,
					},
				},
			},
			want: []Track{
				{
					Name:      "testGetAllTracksMany",
					Slug:      "testGetAllTracksMany",
					IsDefault: true,
				},
				{
					Name:      "testGetAllTracksMany1",
					Slug:      "testGetAllTracksMany1",
					IsDefault: true,
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllTracksWrongApplication",
			args: args{
				tracks: []Track{
					{
						ApplicationId: "",
						Name:          "testGetAllTracksWrongApplication",
						Slug:          "testGetAllTracksWrongApplication",
						IsDefault:     true,
					},
				},
				appId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			want: []Track{
				{
					Name:      "testGetAllTracksWrongApplication",
					Slug:      "testGetAllTracksWrongApplication",
					IsDefault: true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
				Name:              "testGetAllTracksApplication",
				Slug:              "testGetAllTracksApplication",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare GetAllTracks() error = %v", err)
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

			var appId string
			if len(tt.args.appId) > 0 {
				appId = tt.args.appId
			} else {
				appId = app.Id
			}
			got, err := GetAllTracks(appId)
			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTracks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) && !tt.wantErr {
				t.Errorf("GetAllTracks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTrackById(t *testing.T) {
	type args struct {
		track Track
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetTrackByIdSuccess",
			args: args{
				track: Track{
					Name:      "testGetTrackByIdSuccess",
					Slug:      "testGetTrackByIdSuccess",
					IsDefault: true,
				},
			},
			wantErr: false,
		},
		{
			name: "GetTrackByIdWrongApplicationId",
			args: args{
				track: Track{
					ApplicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Name:          "testGetTrackByIdWrongApplicationId",
					Slug:          "testGetTrackByIdWrongApplicationId",
					IsDefault:     true,
				},
			},
			wantErr: true,
		},
		{
			name: "GetTrackByIdWrongId",
			args: args{
				track: Track{
					Id:        "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
					Name:      "testGetTrackByIdWrongId",
					Slug:      "testGetTrackByIdWrongId",
					IsDefault: true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
				Name:              "testGetTrackByIdApplication",
				Slug:              "testGetTrackByIdApplication",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
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
		track Track
		slug  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetTrackBySlugSuccess",
			args: args{
				track: Track{
					Name:      "testGetTrackBySlugSuccess",
					Slug:      "testGetTrackBySlugSuccess",
					IsDefault: true,
				},
				slug: "testGetTrackBySlugSuccess",
			},
			wantErr: false,
		},
		{
			name: "GetTrackBySlugWrongSlug",
			args: args{
				track: Track{
					Name:      "testGetTrackBySlugWrongSlug",
					Slug:      "testGetTrackBySlugWrongSlug",
					IsDefault: true,
				},
				slug: "wrongSlug",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
				Name:              "testGetTrackBySlugApplication",
				Slug:              "testGetTrackBySlugApplication",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare GetTrackBySlug() error = %v", err)
				test.CleanTables()
				return
			}

			tt.args.track.ApplicationId = app.Id
			track, err := CreateTrack(tt.args.track)
			if err != nil {
				t.Errorf("Prepare GetTrackBySlug() error = %v", err)
				test.CleanTables()
				return
			}

			got, err := GetTrackBySlug(tt.args.slug, app.Id)
			test.CleanTables()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrackBySlug() error = %v, wantErr %v, got %v", err, tt.wantErr, got)
				return
			}
			if !reflect.DeepEqual(got, track) && !tt.wantErr {
				t.Errorf("GetTrackBySlug() got = %v, want %v", got, track)
			}
		})
	}
}

func TestUpdateTrack(t *testing.T) {
	type args struct {
		track              Track
		testTrack          Track
		additionalTracks   []Track
		missingTrack       bool
		missingApplication bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Track
		wantErr bool
	}{
		{
			name: "UpdateTrackFields",
			args: args{
				track: Track{
					Name:      "testUpdateTrackFields",
					Slug:      "testUpdateTrackFields",
					IsDefault: false,
				},
				testTrack: Track{
					Name:      "testUpdateTrackFields1",
					Slug:      "testUpdateTrackFields1",
					IsDefault: true,
				},
				missingTrack:       false,
				missingApplication: false,
			},
			wantErr: false,
		},
		{
			name: "UpdateTrackNameExists",
			args: args{
				track: Track{
					Name:      "testUpdateTrackNameExists",
					Slug:      "testUpdateTrackNameExists",
					IsDefault: false,
				},
				testTrack: Track{
					Name:      "testUpdateTrackNameExists2",
					Slug:      "testUpdateTrackNameExists1",
					IsDefault: true,
				},
				additionalTracks: []Track{
					{
						Name:      "testUpdateTrackNameExists2",
						Slug:      "testUpdateTrackNameExists2",
						IsDefault: false,
					},
				},
				missingTrack:       false,
				missingApplication: false,
			},
			wantErr: true,
		},
		{
			name: "UpdateTrackSlugExists",
			args: args{
				track: Track{
					Name:      "testUpdateTrackSlugExists",
					Slug:      "testUpdateTrackSlugExists",
					IsDefault: false,
				},
				testTrack: Track{
					Name:      "testUpdateTrackSlugExists1",
					Slug:      "testUpdateTrackSlugExists2",
					IsDefault: true,
				},
				additionalTracks: []Track{
					{
						Name:      "testUpdateTrackSlugExists2",
						Slug:      "testUpdateTrackSlugExists2",
						IsDefault: false,
					},
				},
				missingTrack:       false,
				missingApplication: false,
			},
			wantErr: true,
		},
		{
			name: "UpdateTrackTrackNotFound",
			args: args{
				track: Track{
					Name:      "testUpdateTrackTrackNotFound",
					Slug:      "testUpdateTrackTrackNotFound",
					IsDefault: false,
				},
				testTrack: Track{
					Name:      "testUpdateTrackTrackNotFound1",
					Slug:      "testUpdateTrackTrackNotFound2",
					IsDefault: true,
				},
				missingTrack:       true,
				missingApplication: false,
			},
			wantErr: true,
		},
		{
			name: "UpdateTrackApplicationNotFound",
			args: args{
				track: Track{
					Name:      "testUpdateTrackApplicationNotFound",
					Slug:      "testUpdateTrackApplicationNotFound",
					IsDefault: false,
				},
				testTrack: Track{
					Name:      "testUpdateTrackApplicationNotFound1",
					Slug:      "testUpdateTrackApplicationNotFound2",
					IsDefault: true,
				},
				missingTrack:       true,
				missingApplication: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
				Name:              "testUpdateTrackApplication",
				Slug:              "testUpdateTrackApplication",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare UpdateTrack() error = %v", err)
				test.CleanTables()
				return
			}
			for _, track := range tt.args.additionalTracks {
				track.ApplicationId = app.Id
				if _, err := CreateTrack(track); err != nil {
					t.Errorf("Prepare UpdateTrack() error = %v", err)
					test.CleanTables()
					return
				}
			}

			tt.args.track.ApplicationId = app.Id
			track, err := CreateTrack(tt.args.track)
			if err != nil {
				t.Errorf("Prepare UpdateTrack() error = %v", err)
				test.CleanTables()
				return
			}
			if tt.args.missingTrack {
				track.Id = "e2ebb12e-e77d-4618-ba79-3f26e8af239a"
			}
			if tt.args.missingApplication {
				track.ApplicationId = "e2ebb12e-e77d-4618-ba79-3f26e8af239a"
			}
			track.Name = tt.args.testTrack.Name
			track.Slug = tt.args.testTrack.Slug
			track.IsDefault = tt.args.testTrack.IsDefault

			got, err := UpdateTrack(*track)
			test.CleanTables()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTrack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, track) != tt.wantErr {
				t.Errorf("UpdateTrack() got = %v, want %v", got, track)
			}
		})
	}
}

func TestDeleteTrackById(t *testing.T) {
	type args struct {
		track Track
		id    string
		AppId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "DeleteTrackByIdSuccess",
			args: args{
				id: "",
				track: Track{
					Name:      "testDeleteTrackByIdSuccess",
					Slug:      "testDeleteTrackByIdSuccess",
					IsDefault: false,
				},
			},
			wantErr: false,
		},
		{
			name: "DeleteTrackByIdDoesNotExist",
			args: args{
				id: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				track: Track{
					Name:      "testDeleteTrackByIdDoesNotExist",
					Slug:      "testDeleteTrackByIdDoesNotExist",
					IsDefault: false,
				},
			},
			wantErr: true,
		},
		{
			name: "DeleteTrackByIdApplicationDoesNotExist",
			args: args{
				AppId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				track: Track{
					Name:      "testDeleteTrackByIdApplicationDoesNotExist",
					Slug:      "testDeleteTrackByIdApplicationDoesNotExist",
					IsDefault: false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := CreateApplication(Application{
				Name:              "testDeleteTrackByIdApplication",
				Slug:              "testDeleteTrackByIdApplication",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare DeleteTrackById() error = %v", err)
				test.CleanTables()
				return
			}
			tt.args.track.ApplicationId = app.Id
			track, err := CreateTrack(tt.args.track)
			if err != nil {
				t.Errorf("Prepare DeleteTrackById() error = %v", err)
				test.CleanTables()
				return
			}
			if len(tt.args.AppId) > 0 {
				if err := DeleteTrackById(track.Id, tt.args.AppId); (err != nil) != tt.wantErr {
					t.Errorf("DeleteTrackById() error = %v, wantErr %v", err, tt.wantErr)
					test.CleanTables()
				}
			} else if len(tt.args.id) > 0 {
				if err := DeleteTrackById(tt.args.id, track.ApplicationId); (err != nil) != tt.wantErr {
					t.Errorf("DeleteTrackById() error = %v, wantErr %v", err, tt.wantErr)
					test.CleanTables()
				}
			} else {
				if err := DeleteTrackById(track.Id, track.ApplicationId); (err != nil) != tt.wantErr {
					t.Errorf("DeleteTrackById() error = %v, wantErr %v", err, tt.wantErr)
					test.CleanTables()
				}
			}
			test.CleanTables()
		})
	}
}
