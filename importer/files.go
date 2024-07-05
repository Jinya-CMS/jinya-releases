package importer

import (
	"github.com/gabriel-vasile/mimetype"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/storage"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ImportFromFileSystem(args []string) {
	appSlug := args[0]
	trackSlug := args[1]
	directory := args[2]

	app, err := models.GetApplicationBySlug(appSlug)
	if err != nil {
		log.Fatalf("Failed to load the application with slug %s: %s", appSlug, err.Error())
		return
	}

	track, err := models.GetTrackBySlug(trackSlug, app.Id)
	if err != nil {
		log.Fatalf("Failed to load the track with slug %s: %s", trackSlug, err.Error())
		return
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Failed to read directory: %s", err.Error())
		return
	}

	readers := make([]io.ReadCloser, 0)

	for _, file := range files {
		if !file.IsDir() {
			log.Printf("Importing file %s", file.Name())
			versionNumber := strings.TrimSuffix(file.Name(), ".zip")
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			versionToUploadBinaryFor, err := models.GetVersionByTrackAndVersion(track.Id, versionNumber)
			if err != nil {
				log.Printf("Could not find version %s, lets create it", versionNumber)
				versionToUploadBinaryFor, err = models.CreateVersion(models.Version{
					ApplicationId: app.Id,
					TrackId:       track.Id,
					Version:       versionNumber,
					UploadDate:    fileInfo.ModTime(),
				})

				if err != nil {
					log.Println("Version not found and cannot be created")
					continue
				}
			}

			path := filepath.Join(directory, fileInfo.Name())
			fileReader, err := os.Open(path)
			if err != nil {
				log.Printf("Failed to open file %s: %s", file.Name(), err.Error())
				continue
			}

			mtype, err := mimetype.DetectFile(path)
			if err != nil {
				mtype = mimetype.Lookup("application/octet-stream")
			}

			log.Printf("File is of type %s", mtype.String())
			readers = append(readers, fileReader)
			err = storage.UploadVersion(fileReader, fileInfo.Size(), mtype.String(), versionToUploadBinaryFor)
			if err != nil {
				log.Printf("Failed to upload version %s: %s", versionToUploadBinaryFor.Version, err.Error())
				continue
			}

			log.Printf("Uploaded version %s", versionToUploadBinaryFor.Version)
		}
	}

	log.Println("Closing all file readers")
	for _, reader := range readers {
		err = reader.Close()
		if err != nil {
			log.Printf("Failed to close file reader: %s", err.Error())
		}
	}
}
