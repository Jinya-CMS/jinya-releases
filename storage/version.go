package storage

import (
	"fmt"
	"io"
	"jinya-releases/database"

	"github.com/google/uuid"
)

const versionBinaryFormat = "%s/%s/%s"

func UploadVersion(reader io.Reader, contentLength int64, contentType string, version *database.Version) error {
	versionId, err := database.GetDbMap().SelectStr(`
insert into version (application_id, track_id, version, upload_date)
values ($1, $2, $3, now())
on conflict (application_id, track_id, version) do update set upload_date = now()
returning id
`, version.ApplicationId, version.TrackId, version.Version)
	if err != nil {
		return err
	}

	version.Id, err = uuid.Parse(versionId)
	if err != nil {
		return err
	}

	return saveFile(fmt.Sprintf(versionBinaryFormat, version.ApplicationId, version.TrackId, versionId), reader, contentLength, contentType)
}

func DownloadVersion(applicationId, trackId, id string) (io.ReadCloser, string, int64, error) {
	return getFile(fmt.Sprintf(versionBinaryFormat, applicationId, trackId, id))
}

func DeleteVersion(applicationId, trackId, id string) error {
	_, err := database.GetDbMap().Exec("delete from version where application_id = $1 and track_id = $2 and id = $3", applicationId, trackId, id)
	if err != nil {
		return err
	}

	return deleteFile(fmt.Sprintf(versionBinaryFormat, applicationId, trackId, id))
}
