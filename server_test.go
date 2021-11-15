package notesbot_test

import (
	"bytes"
	notesbot "deni/notesbot"
	restclient "deni/notesbot/utils"
	mocks "deni/notesbot/utils/mocks"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	firstBody = `Title: Post 1
Tags: tdd, go
---
Hello world Body`
	secondBody = `Title: Post 2
Tags: rust, borrow-checker
---
B
L
M`
)

type StubStore struct {
	notes []notesbot.Note
}

func (s *StubStore) StoreNotes(notes []notesbot.Note) {
	s.notes = notes
}

func (s *StubStore) AllNotes() []notesbot.Note {
	return s.notes
}

func init() {
	restclient.Client = &mocks.MockClient{}
}

func TestPOSTNotesReceiver(t *testing.T) {
	store := StubStore{}
	server := notesbot.NewServer(&store, "token")

	var notes []notesbot.Note
	notes = append(notes, notesbot.Note{
		Title: "hello world",
		Body:  firstBody,
	})
	notes = append(notes, notesbot.Note{
		Title: "hello-world2",
		Body:  firstBody,
	})
	jsonBytes, _ := json.Marshal(notes)

	t.Run("Sync notes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/sync-notes", bytes.NewReader(jsonBytes))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		gotNotesCount := len(store.notes)
		wantNotesCount := 2

		if gotNotesCount != wantNotesCount {
			t.Errorf("response body is wrong, got %d want %d", gotNotesCount, wantNotesCount)
		}
	})
}

func TestPOSTNotesTransmitter(t *testing.T) {
	var notes []notesbot.Note
	httpposturl := "https://reqres.in/api/users"
	json := `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}}`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 203,
			Body:       r,
		}, nil
	}

	request, _ := http.NewRequest("POST", httpposturl, r)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := restclient.Client
	response, error := restclient.SendNotes(client, httpposturl, notes)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	got := response.StatusCode
	want := 203

	if got != want {
		t.Errorf("response StatusCode is wrong, got %q want %q", got, want)
	}
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
