package storage

import (
	"fmt"
	"io"
	"jinya-releases/database/models"
	"net/http"
)

const versionBinaryFormat = "%s/%s/%s"

func UploadVersion(r *http.Request, version *models.Version) error {
	return SaveFile(fmt.Sprintf(versionBinaryFormat, version.ApplicationId, version.TrackId, version.Id), r.Body, r.ContentLength, r.Header.Get("Content-Type"))
}

func DownloadVersion(applicationId, trackId, id string) (io.ReadCloser, string, error) {
	return GetFile(fmt.Sprintf(appLogoFormat, id))
}

func DeleteVersion(applicationId, trackId, id string) error {
	err := DeleteFile(fmt.Sprintf(versionBinaryFormat, applicationId, trackId, id))

	return err
}
