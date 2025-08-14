package storage

import (
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

const appLogoFormat = "application/%s/logo"

func UploadLogo(r *http.Request) error {
	id := mux.Vars(r)["id"]

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
		return fmt.Errorf("file format not supported")
	}

	return saveFile(fmt.Sprintf(appLogoFormat, id), r.Body, r.ContentLength, contentType)
}

func DownloadLogo(id string) (io.ReadCloser, string, int64, error) {
	return getFile(fmt.Sprintf(appLogoFormat, id))
}
