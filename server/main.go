package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	notesbot "deni/notesbot"

	"gopkg.in/yaml.v2"
)

type instanceConfig struct {
	ServerHost    string `yaml:"server_host"`
	ServerPort    string `yaml:"server_port"`
	TelegramToken string `yaml:"telegram_token"`
	DbFileName    string `yaml:"db_file_name"`
}

func (c *instanceConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func main() {
	data, err := ioutil.ReadFile("settings.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config instanceConfig

	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}

	db, err := os.OpenFile(config.DbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", config.DbFileName, err)
	}

	store, err := notesbot.NewFileSystemStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := notesbot.NewServer(store, config.TelegramToken)
	log.Fatal(http.ListenAndServe(":5000", server))
}
