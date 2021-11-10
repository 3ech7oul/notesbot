package notesbot_test

import (
	"bytes"
	notesbot "deni/notesbot"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

const (
	firstBody = `Title: Post 1
Tags: tdd, go
---
Hello world`
	secondBody = `Title: Post 2
Tags: rust, borrow-checker
---
B
L
M`
)

func TestPOSTNotes(t *testing.T) {
	server := notesbot.NewServer()

	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	notes, _ := notesbot.NewNotesFromFS("", fs)

	jsonBytes, _ := json.Marshal(notes)

	t.Run("Import notes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/sync-notes", bytes.NewReader(jsonBytes))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		fmt.Println(response.Body.String())
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
