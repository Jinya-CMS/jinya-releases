package storage

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"jinya-releases/service"
	"net/http"
	"slices"
)

const appLogoFormat = "application/%s/logo"

func UploadLogo(r *http.Request) (errDetails *service.ErrorDetails, status int) {
	id := mux.Vars(r)["id"]

	_, errDetails, status = service.GetApplicationById(id)
	if errDetails != nil {
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
		errDetails = &service.ErrorDetails{
			EntityType: "application",
			Message:    "file format not supported",
			ErrorType:  "storage",
		}
		status = http.StatusUnsupportedMediaType
		return
	}

	err := SaveFile(fmt.Sprintf(appLogoFormat, id), r.Body, r.ContentLength, contentType)
	if err != nil {
		errDetails = &service.ErrorDetails{
			EntityType: "application",
			Message:    "upload failed",
			ErrorType:  "storage",
		}
		status = http.StatusConflict
	}

	return
}

func DownloadLogo(id string) (io.ReadCloser, string, error) {
	return GetFile(fmt.Sprintf(appLogoFormat, id))
}
