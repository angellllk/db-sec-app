package main

import (
	"db-sec/core"
	"log"
)

func init() { log.SetFlags(log.Lshortfile | log.LstdFlags) }

func main() {
	orclDb, errDB := core.ConnectDB()
	if errDB != nil {
		return
	}
	defer orclDb.Close()

	// Create table
	errPrep := core.PrepareDB(orclDb)
	if errPrep != nil {
		return
	}

	errStart := core.StartApp(orclDb)
	if errStart != nil {
		log.Print(errStart)
	}
}
