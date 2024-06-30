package storage

import (
	"fmt"
	"jinya-releases/database/models"
	"net/http"
)

const versionBinaryFormat = "%s/%s/%s"

func UploadVersion(r *http.Request, version *models.Version) error {
	return SaveFile(fmt.Sprintf(versionBinaryFormat, version.ApplicationId, version.TrackId, version.Id), r.Body, r.ContentLength, r.Header.Get("Content-Type"))
}
