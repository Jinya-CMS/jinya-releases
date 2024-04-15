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

	if _, err = conn.Exec("DELETE FROM application"); err != nil {
		log.Printf("Failed to clean table application %v", err)
	}
}
