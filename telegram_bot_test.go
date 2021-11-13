package notesbot_test

import (
	"bytes"
	notesbot "deni/notesbot"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func TestTelegramBot(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	notes, _ := notesbot.NewNotesFromFS("", fs)
	store := StubStore{notes: notes}
	bot := notesbot.NewServer(&store, "t")

	t.Run("List Message", func(t *testing.T) {
		got := bot.ListMessage()
		want := `/hello world \n /hello-world2`

		if want != got {
			t.Errorf("response body is wrong, got %s want %s", got, want)
		}

	})

	t.Run("One Message", func(t *testing.T) {
		got := bot.OneMessage("hello world")
		want := `hello world \n Hello world Body`

		if want != got {
			t.Errorf("response body is wrong, got %s want %s", got, want)
		}
	})

	t.Run("Send Message", func(t *testing.T) {
		j := &webhookReqBody{}
		j.Message.Text = "get list"
		j.Message.Chat.ID = 111

		jsonBytes, _ := json.Marshal(j)

		request, _ := http.NewRequest(http.MethodPost, "/bot", bytes.NewReader(jsonBytes))
		response := httptest.NewRecorder()

		bot.ServeHTTP(response, request)

	})
}
