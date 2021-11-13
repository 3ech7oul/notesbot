package main

import (
	notesbot "deni/notesbot"
	restclient "deni/notesbot/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	notes, err := notesbot.NewNotesFromFS("../knowledge", os.DirFS("../knowledge"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Posts indexed: %d\n", len(notes))

	client := restclient.Client
	response, error := restclient.SendNotes(client, "https://notesbot-r96tq.ondigitalocean.app/sync-notes", notes)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}
