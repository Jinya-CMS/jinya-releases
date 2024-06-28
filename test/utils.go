package test

import (
	"jinya-releases/database"
	"log"
)

func CleanTables() {
	conn, err := database.Connect()
	if err != nil {
		log.Printf("Failed to clean tables %v", err)
		return
	}

	if _, err = conn.Exec("DELETE FROM pushtokenapplication"); err != nil {
		log.Printf("Failed to clean table pushtokenapplication %v", err)
	}
	if _, err = conn.Exec("DELETE FROM pushtoken"); err != nil {
		log.Printf("Failed to clean table pushtoken %v", err)
	}
	if _, err = conn.Exec("DELETE FROM version"); err != nil {
		log.Printf("Failed to clean table version %v", err)
	}
	if _, err = conn.Exec("DELETE FROM track"); err != nil {
		log.Printf("Failed to clean table track %v", err)
	}
	if _, err = conn.Exec("DELETE FROM application"); err != nil {
		log.Printf("Failed to clean table application %v", err)
	}
}
