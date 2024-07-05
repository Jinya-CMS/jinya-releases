package storage

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/utils"
	"net/http"
	"slices"
)

const appLogoFormat = "application/%s/logo"

func UploadLogo(r *http.Request) (errDetails *utils.ErrorDetails, status int) {
	id := mux.Vars(r)["id"]

	_, err := models.GetApplicationById(id)
	if err != nil {
		status = http.StatusNotFound
		errDetails = &utils.ErrorDetails{
			EntityType: "application",
			Message:    "Application not found",
			ErrorType:  "request",
		}

		return
	}

	status = http.StatusNoContent

	contentType := r.Header.Get("Content-Type")
	if !slices.Contains([]string{
		"image/apng",
		"image/avif",
		"image/gif",
		"image/jpeg",
		"image/png",
		"image/svg+xml",
		"image/webp",
	}, contentType) {
		errDetails = &utils.ErrorDetails{
			EntityType: "application",
			Message:    "file format not supported",
			ErrorType:  "storage",
		}
		status = http.StatusUnsupportedMediaType
		return
	}

	err = SaveFile(fmt.Sprintf(appLogoFormat, id), r.Body, r.ContentLength, contentType)
	if err != nil {
		errDetails = &utils.ErrorDetails{
			EntityType: "application",
			Message:    "upload failed",
			ErrorType:  "storage",
		}
		status = http.StatusConflict
	}

	return
}

func DownloadLogo(id string) (io.ReadCloser, string, int64, error) {
	return GetFile(fmt.Sprintf(appLogoFormat, id))
}
