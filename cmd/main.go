package main

import (
	notesbot "deni/notesbot"
	restclient "deni/notesbot/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type instanceConfig struct {
	ServerUrl string `yaml:"server_url"`
	NotesPath string `yaml:"notes_path"`
}

func (c *instanceConfig) parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func main() {
	var config instanceConfig
	dataconf, err := ioutil.ReadFile("settings.yaml")
	if err := config.parse(dataconf); err != nil {
		fmt.Println(err)
	}

	notes, err := notesbot.NewNotesFromFS(config.NotesPath, os.DirFS(config.NotesPath))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Posts indexed: %d\n", len(notes))

	client := restclient.Client
	response, error := restclient.SendNotes(client, config.ServerUrl, notes)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}
