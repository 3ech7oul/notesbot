package notesbot_test

import (
	notesbot "deni/notesbot"
	"errors"
	"io/fs"
	"reflect"
	"testing"
)

type StubFailingFS struct{}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}

func assertNote(t *testing.T, got notesbot.Note, want notesbot.Note) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestFindNoteByAttribute(t *testing.T) {
	var notes []notesbot.Note
	notes = append(notes, notesbot.Note{
		Title:      "hello world",
		Body:       "Hello world Body",
		Tags:       []string{"tdd", "go"},
		TelegramId: 1,
	})
	notes = append(notes, notesbot.Note{
		Title:      "hello-world2",
		Body:       "secondBody",
		TelegramId: 2,
	})

	t.Run("Note found", func(t *testing.T) {
		note, _ := notesbot.FindNoteByAttribute(notes, 1)
		assertNote(t, note, notesbot.Note{
			Title:      "hello world",
			Body:       "Hello world Body",
			Tags:       []string{"tdd", "go"},
			TelegramId: 1,
		})
	})

	t.Run("Note not found", func(t *testing.T) {
		note, err := notesbot.FindNoteByAttribute(notes, 3)
		assertNote(t, note, notesbot.Note{})

		if err == nil {
			t.Fatal("expected an error")
		}
	})

	t.Run("Titles index", func(t *testing.T) {
		index := notesbot.AllNotesTitleList(notes)
		got := index[1]
		want := "hello-world2"
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}
