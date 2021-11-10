package notesbot_test

import (
	notesbot "deni/notesbot"
	"io/ioutil"
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func TestFileSystemStore(t *testing.T) {

	t.Run("Read from file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[{"Title":"hello world", "Body":"Hello world"}]`)
		defer cleanDatabase()

		store, err := notesbot.NewFileSystemStore(database)

		assertNoError(t, err)

		got, _ := notesbot.FindNoteByAttribute(store.Notes, "hello world")

		want := notesbot.Note{
			Title: "hello world",
			Body:  `Hello world`,
		}

		assertNote(t, got, want)
	})

}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
