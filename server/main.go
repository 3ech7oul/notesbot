package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	notesbot "deni/notesbot"

	"gopkg.in/yaml.v2"
)

type instanceConfig struct {
	ServerPort     string `yaml:"server_port"`
	TelegramToken  string `yaml:"telegram_token"`
	DbFileName     string `yaml:"db_file_name"`
	TelegramSecret string `yaml:"secret"`
}

func (c *instanceConfig) parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func main() {
	var config instanceConfig
	dataconf, err := ioutil.ReadFile("settings.yaml")
	if err != nil {
		fmt.Println(err)

		config.ServerPort = os.Getenv("s_PORT")
		config.TelegramToken = os.Getenv("TOKEN")
		config.DbFileName = os.Getenv("DB_FILE")
		config.TelegramSecret = os.Getenv("SECRET")

	} else {
		if err := config.parse(dataconf); err != nil {
			fmt.Println(err)
		}
	}

	db, err := os.OpenFile(config.DbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", config.DbFileName, err)
	}

	store, err := notesbot.NewFileSystemStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := notesbot.NewServer(store, config.TelegramToken, config.TelegramSecret)
	log.Fatal(http.ListenAndServe(":5000", server))
}
