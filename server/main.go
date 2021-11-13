package main

import (
	"log"
	"net/http"
	"os"

	notesbot "deni/notesbot"
)

const dbFileName = "notes.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := notesbot.NewFileSystemStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := notesbot.NewServer(store)
	log.Fatal(http.ListenAndServe(":5000", server))
}
