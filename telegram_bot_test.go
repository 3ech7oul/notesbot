package notesbot_test

import (
	notesbot "deni/notesbot"
	"testing"
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
	var notes []notesbot.Note
	notes = append(notes, notesbot.Note{
		Title:      "hello world",
		Body:       firstBody,
		TelegramId: 50,
	})
	store := StubStore{notes: notes}
	bot := notesbot.NewServer(&store, "t")

	t.Run("List Message", func(t *testing.T) {
		got := bot.ListMessage()
		want := `/50 hello world`

		if want != got {
			t.Errorf("response body is wrong, got %s want %s", got, want)
		}

	})

}
