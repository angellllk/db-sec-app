package main

import (
	"db-sec/db"
	"log"
)

func init() { log.SetFlags(log.Lshortfile | log.LstdFlags) }

func main() {
	orclDb, errDB := db.ConnectDB()
	if errDB != nil {
		return
	}
	defer orclDb.Close()

	db.CreateTables(orclDb)
	db.DeleteTables(orclDb)
}
