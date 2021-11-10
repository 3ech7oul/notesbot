package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	notesbot "deni/notesbot"
)

const TELEGRAM_POST_URL = "https://api.telegram.org/bot777845702:AAFdPS_taJ3pTecEFv2jXkmbQfeOqVZGERw/sendMessage"

var notes []notesbot.Note

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

// https://core.telegram.org/bots/api#sendmessage
type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func Handler(res http.ResponseWriter, req *http.Request) {
	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	if !strings.Contains(strings.ToLower(body.Message.Text), "get list") {

		if err := sendList(body.Message.Chat.ID, notesbot.AllNotesTitleList(notes)); err != nil {
			fmt.Println("send list:", err)
			return
		}

		return
	}

	serchngNote := strings.ToLower(body.Message.Text)
	note, err := notesbot.FindNoteByAttribute(notes, serchngNote)

	if nil != err {
		fmt.Println("error in sending reply:", err)
		return
	}

	if err := sendResponce(body.Message.Chat.ID, note); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}

	fmt.Println("reply sent")
}

func sendResponce(chatID int64, note notesbot.Note) error {

	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   note.Body,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post(TELEGRAM_POST_URL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func sendList(chatID int64, titleListSlice []string) error {
	titlesList := strings.Join(titleListSlice, ",")
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   titlesList,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post(TELEGRAM_POST_URL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

// FInally, the main funtion starts our server on port 3000
func main() {
	var err error
	notes, err = notesbot.NewNotesFromFS("../knowledge", os.DirFS("../knowledge"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Posts indexed: %d\n", len(notes))

	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}
