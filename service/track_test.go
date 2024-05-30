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
)

func TestGetAllTracks(t *testing.T) {
	type args struct {
		applicationId string
		tracks        []models.Track
		appId         string
	}
	tests := []struct {
		name       string
		args       args
		wantTracks []models.Track
		wantErr    bool
		wantStatus int
	}{
		{
			name: "GetAllTracksNone",
			args: args{
				tracks: make([]models.Track, 0),
			},
			wantTracks: make([]models.Track, 0),
			wantErr:    false,
		},
		{
			name: "GetAllTracksOne",
			args: args{
				tracks: []models.Track{
					{
						Name:      "test",
						Slug:      "test",
						IsDefault: true,
					},
				},
			},
			wantTracks: []models.Track{
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
				tracks: []models.Track{
					{
						Name:      "test",
						Slug:      "test",
						IsDefault: true,
					},
					{
						Name:      "test1",
						Slug:      "test1",
						IsDefault: false,
					},
				},
			},
			wantTracks: []models.Track{
				{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				{
					Name:      "test1",
					Slug:      "test1",
					IsDefault: false,
				},
			},
			wantErr: false,
		},
		{
			name: "GetAllTracksApplicationWrong",
			args: args{
				tracks: []models.Track{
					{
						Name:      "test",
						Slug:      "test",
						IsDefault: true,
					},
				},
				appId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantTracks: []models.Track{
				{
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
			app, err := models.CreateApplication(models.Application{
				Name:              "track_test",
				Slug:              "track_test",
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
				_, err = models.CreateTrack(track)
			}
			if err != nil {
				t.Errorf("Prepare GetAllTracks() error = %v", err)
				test.CleanTables()
				return
			}

			var appId string
			if len(tt.args.appId) > 0 {
				appId = tt.args.appId
			} else {
				appId = app.Id
			}

			gotTracks, errDetails, _ := GetAllTracks(appId)
			test.CleanTables()

			if (errDetails != nil) != tt.wantErr {
				t.Errorf("GetAllTracks() error = %v, wantErr %v", errDetails, tt.wantErr)
				return
			}
			if (len(gotTracks) != len(tt.wantTracks)) != tt.wantErr {
				t.Errorf("GetAllTracks() gotTracks = %v, want %v", gotTracks, tt.wantTracks)
				return
			}
		})
	}
}

func TestGetTrackById(t *testing.T) {
	type args struct {
		trackId       string
		applicationId string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "GetTrackByIdPositive",
			args:       args{},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "GetTrackByIdWorngTrackId",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "GetTrackByIdWorngApplicationId",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "track_test",
				Slug:              "track_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare GetTrackById() error = %v", err)
				test.CleanTables()
				return
			}
			track, err := models.CreateTrack(models.Track{
				ApplicationId: app.Id,
				Name:          "test",
				Slug:          "test",
				IsDefault:     false,
			})
			if err != nil {
				t.Errorf("Prepare GetTrackById() error = %v", err)
				test.CleanTables()
				return
			}

			var (
				trackId       string
				applicationId string
			)

			trackId = track.Id
			applicationId = app.Id

			if len(tt.args.trackId) > 0 {
				trackId = tt.args.trackId
			}
			if len(tt.args.applicationId) > 0 {
				applicationId = tt.args.applicationId
			}

			gotTrack, gotErrDetails, gotStatus := GetTrackById(trackId, applicationId)
			test.CleanTables()

			if !reflect.DeepEqual(gotTrack, track) != tt.wantErr {
				t.Errorf("GetTrackById() gotTrack = %v, want %v", gotTrack, track)
			}
			if (gotErrDetails != nil) != tt.wantErr {
				t.Errorf("GetTrackById() gotErrDetails = %v, want %v", gotErrDetails, tt.wantErr)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetTrackById() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestCreateTrack(t *testing.T) {
	type args struct {
		reader        io.Reader
		track         createTrackRequest
		applicationId string
		tracks        []models.Track
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "CreateTrackPositive",
			args: args{
				track: createTrackRequest{
					Name:      "test",
					Slug:      "test",
					IsDefault: false,
				}},
			wantErr:    false,
			wantStatus: http.StatusCreated,
		},
		{
			name: "CreateTrackNameMissing",
			args: args{
				track: createTrackRequest{
					Name:      "",
					Slug:      "test",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreateTrackSlugMissing",
			args: args{
				track: createTrackRequest{
					Name:      "test",
					Slug:      "",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreateTrackApplicationNotFound",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				track: createTrackRequest{
					Name:      "test",
					Slug:      "test",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "CreateTrackNameExists",
			args: args{
				tracks: []models.Track{
					{
						Name:      "test",
						Slug:      "test",
						IsDefault: true,
					},
				},
				track: createTrackRequest{
					Name:      "test",
					Slug:      "test1",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
		{
			name: "CreateTrackSlugExists",
			args: args{
				tracks: []models.Track{
					{
						Name:      "test",
						Slug:      "test",
						IsDefault: true,
					},
				},
				track: createTrackRequest{
					Name:      "test1",
					Slug:      "test",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "track_test",
				Slug:              "track_test",
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
				if _, err := models.CreateTrack(track); err != nil {
					t.Errorf("Prepare CreateTrack() error = %v", err)
					test.CleanTables()
					return
				}
			}

			var buffer bytes.Buffer
			_ = json.NewEncoder(&buffer).Encode(tt.args.track)

			var applicationId string
			if len(tt.args.applicationId) > 0 {
				applicationId = tt.args.applicationId
			} else {
				applicationId = app.Id
			}

			_, gotErrDetails, gotStatus := CreateTrack(&buffer, applicationId)
			test.CleanTables()

			if (gotErrDetails != nil) != tt.wantErr {
				t.Errorf("CreateTrack() gotErrDetails = %v, want %v", gotErrDetails, tt.wantErr)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("CreateTrack() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestUpdateTrack(t *testing.T) {
	type args struct {
		updateTrack      updateTrackRequest
		track            models.Track
		applicationId    string
		trackId          string
		additionalTracks []models.Track
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "UpdateTrack",
			args: args{
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				updateTrack: updateTrackRequest{
					Name:      "test1",
					Slug:      "test1",
					IsDefault: false,
				}},
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
		{
			name: "UpdateTrackNameMissing",
			args: args{
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				updateTrack: updateTrackRequest{
					Name:      "",
					Slug:      "test1",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UpdateTrackSlugMissing",
			args: args{
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				updateTrack: updateTrackRequest{
					Name:      "test1",
					Slug:      "",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UpdateTrackWrongTrackId",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				updateTrack: updateTrackRequest{
					Name:      "test1",
					Slug:      "test1",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "UpdateTrackWrongApplicationId",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				updateTrack: updateTrackRequest{
					Name:      "test1",
					Slug:      "test1",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "UpdateTrackNameExists",
			args: args{
				additionalTracks: []models.Track{
					{
						Name:      "test1",
						Slug:      "test2",
						IsDefault: true,
					},
				},
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				updateTrack: updateTrackRequest{
					Name:      "test1",
					Slug:      "test1",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
		{
			name: "UpdateTrackSlugExists",
			args: args{
				additionalTracks: []models.Track{
					{
						Name:      "test2",
						Slug:      "test1",
						IsDefault: true,
					},
				},
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: true,
				},
				updateTrack: updateTrackRequest{
					Name:      "test1",
					Slug:      "test1",
					IsDefault: false,
				}},
			wantErr:    true,
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "track_test",
				Slug:              "track_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare UpdateTrack() error = %v", err)
				test.CleanTables()
				return
			}

			tt.args.track.ApplicationId = app.Id
			track, err := models.CreateTrack(tt.args.track)
			if err != nil {
				t.Errorf("Prepare UpdateTrack() error = %v", err)
				test.CleanTables()
				return
			}

			for _, track := range tt.args.additionalTracks {
				track.ApplicationId = app.Id
				if _, err := models.CreateTrack(track); err != nil {
					t.Errorf("Prepare UpdateTrack() error = %v", err)
					test.CleanTables()
					return
				}
			}

			var buffer bytes.Buffer
			_ = json.NewEncoder(&buffer).Encode(tt.args.updateTrack)

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

			_, gotErrDetails, gotStatus := UpdateTrack(trackId, applicationId, &buffer)
			test.CleanTables()

			if (gotErrDetails != nil) != tt.wantErr {
				t.Errorf("UpdateTrack() gotErrDetails = %v, want %v", gotErrDetails, tt.wantErr)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("UpdateTrack() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestDeleteTrack(t *testing.T) {
	type args struct {
		trackId       string
		applicationId string
		track         models.Track
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "DeleteTrackSuccess",
			args: args{
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: false,
				},
			},
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
		{
			name: "DeleteTrackWrongTrackId",
			args: args{
				trackId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: false,
				},
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "DeleteTrackWrongApplicationId",
			args: args{
				applicationId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
				track: models.Track{
					Name:      "test",
					Slug:      "test",
					IsDefault: false,
				},
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "track_test",
				Slug:              "track_test",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare DeleteTrackById() error = %v", err)
				test.CleanTables()
				return
			}

			tt.args.track.ApplicationId = app.Id
			track, err := models.CreateTrack(tt.args.track)
			if err != nil {
				t.Errorf("Prepare UpdateTrack() error = %v", err)
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

			gotErrDetails, gotStatus := DeleteTrack(trackId, applicationId)
			test.CleanTables()

			if (gotErrDetails != nil) != tt.wantErr {
				t.Errorf("DeleteTrack() gotErrDetails = %v, want %v", gotErrDetails, tt.wantErr)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("DeleteTrack() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}
