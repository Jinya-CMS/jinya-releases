package service

import (
	"bytes"
	"encoding/json"
	"jinya-releases/database/models"
	"jinya-releases/test"
	"net/http"
	"testing"
)

func TestCreatePushtoken(t *testing.T) {
	type args struct {
		pushtoken createPushtokenRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "CreatePushtokenSuccessfull",
			args: args{
				pushtoken: createPushtokenRequest{
					AllowedApps: nil,
				},
			},
			wantErr:    false,
			wantStatus: http.StatusCreated,
		},
		{
			name: "CreatePushtokenAllowedAppsEmpty",
			args: args{
				pushtoken: createPushtokenRequest{
					AllowedApps: nil,
				},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CreatePushtokenAllowedAppsDoesNotExist",
			args: args{
				pushtoken: createPushtokenRequest{
					AllowedApps: []string{"c96a62cc-804e-4652-bee9-aa1c267fefc9"},
				},
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "testCreatePushtoken",
				Slug:              "testCreatePushtoken",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare CreatePushtoken() error = %v", err)
				test.CleanTables()
				return
			}

			if len(tt.args.pushtoken.AllowedApps) == 0 && !tt.wantErr {
				tt.args.pushtoken.AllowedApps = []string{app.Id}
			}

			var buffer bytes.Buffer
			_ = json.NewEncoder(&buffer).Encode(tt.args.pushtoken)
			gotPushtoken, gotErrDetails, gotStatus := CreatePushtoken(&buffer)
			test.CleanTables()

			if (gotErrDetails != nil) != tt.wantErr {
				t.Errorf("CreatePushtoken() error = %v, wantErr %v", gotErrDetails, tt.wantErr)
				return
			}

			if tt.wantStatus != gotStatus {
				t.Errorf("CreatePushtoken() status = %v, wantStatus %v", gotStatus, tt.wantStatus)
				return
			}

			if gotPushtoken != nil && gotPushtoken.Id == "" {
				t.Errorf("CreatePushtoken() id is empty string")
			}
		})
	}
}

func TestDeletePushtoken(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "DeletePushtokenSuccess",
			args: args{
				id: "",
			},
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
		{
			name: "DeletePushtokenDoesNotExist",
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
				Name:              "testDeletePushtoken",
				Slug:              "testDeletePushtoken",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare DeletePushtoken() error = %v", err)
				test.CleanTables()
				return
			}
			token, err := models.CreatePushtoken([]string{app.Id})
			if err != nil {
				t.Errorf("Prepare DeletePushtoken() error = %v", err)
				test.CleanTables()
				return
			}
			if len(tt.args.id) == 0 {
				tt.args.id = token.Id
			}

			errDetails, status := DeletePushtoken(tt.args.id)
			test.CleanTables()
			if (errDetails != nil) != tt.wantErr {
				t.Errorf("DeletePushtoken() error = %v, wantErr %v", errDetails, tt.wantErr)
				test.CleanTables()
				return
			} else if tt.wantStatus != status {
				t.Errorf("DeletePushtoken() status = %v, wantStatus %v", status, tt.wantStatus)
				return
			}
		})
	}
}

func TestGetAllPushtokens(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{
			name:    "GetAllPushtokensNone",
			count:   0,
			wantErr: false,
		},
		{
			name:    "GetAllPushtokensOne",
			count:   1,
			wantErr: false,
		},
		{
			name:    "GetAllPushtokensMany",
			count:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "testGetAllPushtokens",
				Slug:              "testGetAllPushtokens",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare GetAllPushtokens() error = %v", err)
				test.CleanTables()
				return
			}

			for i := 0; i < tt.count; i++ {
				models.CreatePushtoken([]string{app.Id})
			}

			got, errDetails := GetAllPushtokens()
			test.CleanTables()
			if (errDetails != nil) != tt.wantErr {
				t.Errorf("GetAllPushtokens() error = %v, wantErr %v", errDetails, tt.wantErr)
				return
			}
			if len(got) != tt.count {
				t.Errorf("GetAllPushtokens() got = %v, want %v", got, tt.count)
			}
		})
	}
}

func TestGetPushtokenById(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "GetPushtokenByIdSuccess",
			args: args{
				id: "",
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "GetPushtokenByIdFalseId",
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
				Name:              "testGetPushtokenById",
				Slug:              "testGetPushtokenById",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare GetPushtokenById() error = %v", err)
				test.CleanTables()
				return
			}
			token, err := models.CreatePushtoken([]string{app.Id})
			if err != nil {
				t.Errorf("Prepare GetPushtokenById() error = %v", err)
				test.CleanTables()
				return
			}
			if len(tt.args.id) == 0 {
				tt.args.id = token.Id
			}

			_, errDetails, status := GetPushtokenById(tt.args.id)
			test.CleanTables()

			if (errDetails != nil) != tt.wantErr {
				t.Errorf("GetPushtokenById() error = %v, wantErr %v", errDetails, tt.wantErr)
				return
			}
			if tt.wantStatus != status {
				t.Errorf("GetPushtokenById() status = %v, wantStatus %v", status, tt.wantStatus)
				return
			}
		})
	}
}

func TestUpdatePushtoken(t *testing.T) {
	type args struct {
		app   models.Application
		id    string
		appId string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "UpdatePushtokenSuccess",
			args: args{
				app: models.Application{
					Name:              "testUpdatePushtoken2",
					Slug:              "testUpdatePushtoken2",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
			},
			wantErr:    false,
			wantStatus: http.StatusNoContent,
		},
		{
			name: "UpdatePushtokenFalseId",
			args: args{
				app: models.Application{
					Name:              "testUpdatePushtoken2",
					Slug:              "testUpdatePushtoken2",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				id: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "UpdatePushtokenNoApplist",
			args:       args{},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UpdatePushtokenWrongApplist",
			args: args{
				app: models.Application{
					Name:              "testUpdatePushtoken2",
					Slug:              "testUpdatePushtoken2",
					HomepageTemplate:  "test",
					TrackpageTemplate: "test",
				},
				appId: "e2ebb12e-e77d-4618-ba79-3f26e8af239a",
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := models.CreateApplication(models.Application{
				Name:              "testUpdatePushtoken",
				Slug:              "testUpdatePushtoken",
				HomepageTemplate:  "test",
				TrackpageTemplate: "test",
			})
			if err != nil {
				t.Errorf("Prepare UpdatePushtoken() error = %v", err)
				test.CleanTables()
				return
			}
			token, err := models.CreatePushtoken([]string{app.Id})
			if err != nil {
				t.Errorf("Prepare UpdatePushtoken() error = %v", err)
				test.CleanTables()
				return
			}

			newApp := new(models.Application)
			if len(tt.args.app.Name) != 0 {
				newApp, err = models.CreateApplication(tt.args.app)
				if err != nil {
					t.Errorf("Prepare UpdatePushtoken() error = %v", err)
					test.CleanTables()
					return
				}
			}

			id := ""
			if len(tt.args.id) == 0 {
				id = token.Id
			} else {
				id = tt.args.id
			}

			appId := ""
			if len(tt.args.appId) == 0 {
				appId = newApp.Id
			} else {
				appId = tt.args.appId
			}

			request1 := updatePushtokenRequest{
				Id:          id,
				Token:       token.Token,
				AllowedApps: []string{appId},
			}

			request2 := updatePushtokenRequest{
				Id:          id,
				Token:       token.Token,
				AllowedApps: []string{},
			}

			var buffer bytes.Buffer
			if len(appId) != 0 {
				_ = json.NewEncoder(&buffer).Encode(request1)
			} else {
				_ = json.NewEncoder(&buffer).Encode(request2)
			}

			got, errDetails, status := UpdatePushtoken(id, &buffer)
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
