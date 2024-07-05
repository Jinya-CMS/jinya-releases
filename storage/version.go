package storage

import (
	"fmt"
	"io"
	"jinya-releases/database/models"
)

const versionBinaryFormat = "%s/%s/%s"

func UploadVersion(reader io.Reader, contentLength int64, contentType string, version *models.Version) error {
	return SaveFile(fmt.Sprintf(versionBinaryFormat, version.ApplicationId, version.TrackId, version.Id), reader, contentLength, contentType)
}

func DownloadVersion(applicationId, trackId, id string) (io.ReadCloser, string, int64, error) {
	return GetFile(fmt.Sprintf(versionBinaryFormat, applicationId, trackId, id))
}

func DeleteVersion(applicationId, trackId, id string) error {
	return DeleteFile(fmt.Sprintf(versionBinaryFormat, applicationId, trackId, id))
}
