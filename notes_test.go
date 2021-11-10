package notesbot_test

import (
	notesbot "deni/notesbot"
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

type StubFailingFS struct{}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}

func TestNewNotes(t *testing.T) {

	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	posts, err := notesbot.NewNotesFromFS("", fs)

	if err != nil {
		t.Fatal(err)
	}

	// rest of test code cut for brevity
	assertNote(t, posts[0], notesbot.Note{
		Title: "hello world",
		Tags:  []string{"tdd", "go"},
		Body:  `Hello world`,
	})
}

func assertNote(t *testing.T, got notesbot.Note, want notesbot.Note) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestFindNoteByAttribute(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	notes, err := notesbot.NewNotesFromFS("", fs)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("Note found", func(t *testing.T) {
		note, _ := notesbot.FindNoteByAttribute(notes, "hello world")
		assertNote(t, note, notesbot.Note{
			Title: "hello world",
			Tags:  []string{"tdd", "go"},
			Body:  `Hello world`,
		})
	})

	t.Run("Note not found", func(t *testing.T) {
		note, err := notesbot.FindNoteByAttribute(notes, "hello world asd")
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
