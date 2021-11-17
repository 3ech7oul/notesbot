package notesbot_test

import (
	notesbot "deni/notesbot"
	"testing"
)

func TestTelegramBot(t *testing.T) {
	var notes []notesbot.Note
	notes = append(notes, notesbot.Note{
		Title:      "hello world",
		Body:       firstBody,
		TelegramId: 345,
	})
	store := StubStore{notes: notes}
	bot := notesbot.NewServer(&store, "t", "s")

	t.Run("List Message", func(t *testing.T) {
		got := bot.ListMessage()
		want := `/345 hello world`

		if want != got {
			t.Errorf("response body is wrong, got %s want %s", got, want)
		}

	})

}

func TestCommand(t *testing.T) {

}
