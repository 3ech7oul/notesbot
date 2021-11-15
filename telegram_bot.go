package notesbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type NotesBoot struct {
	store Store
	http.Handler
}

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

func (n *NotesServer) botHandler(res http.ResponseWriter, req *http.Request) {

	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	if strings.Contains(strings.ToLower(body.Message.Text), "get list") {
		if err := n.sendList(body.Message.Chat.ID); err != nil {
			fmt.Println("send list:", err)

			return
		}

		return
	}

	s := body.Message.Text
	requestdNote := strings.Replace(s, "/", "", -1)

	note, err := FindNoteByAttribute(n.store.AllNotes(), strconv.ParseInt(requestdNote, 10, 64))

	if nil != err {
		fmt.Println("error in sending reply:", err)

		return
	}

	if err := n.sendResponce(body.Message.Chat.ID, note); err != nil {
		fmt.Println("error in sending reply:", err)

		return
	}

	fmt.Println("reply sent")
}

func (n *NotesServer) ComandHelper(command string) string {

	return strings.Replace(command, "/", "", -1)
}

func (n *NotesServer) sendResponce(chatID int64, note Note) error {

	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   note.Body,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(n.urlPost(), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func (n *NotesServer) sendList(chatID int64) error {
	titlesList := n.ListMessage()

	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   titlesList,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(n.urlPost(), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return errors.New("response Body:" + string(body))
	}

	return nil
}

func (n *NotesServer) ListMessage() string {
	var titles []string
	for _, note := range n.store.AllNotes() {
		titles = append(titles, fmt.Sprintf("/%d %s", note.TelegramId, note.Title))
	}

	return strings.Join(titles[:], "\n")
}

func (n *NotesServer) OneMessage(id int64) string {
	var message string

	note, err := FindNoteByAttribute(n.store.AllNotes(), id)
	if nil != err {
		return "Message not found"
	}

	message = fmt.Sprintf(`%s \n %s`, note.Title, note.Body)

	return message

}

func (n *NotesServer) urlPost() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.telegramToken)
}
