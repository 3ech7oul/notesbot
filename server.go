package notesbot

import (
	"encoding/json"
	"net/http"
)

const jsonContentType = "application/json"

var notes []Note

type NotesServer struct {
	http.Handler
}

type Store interface {
	Get(name string) Note
	Commmit(notes []Note)
}

func NewServer() *NotesServer {
	s := new(NotesServer)

	router := http.NewServeMux()
	router.Handle("/sync-notes", http.HandlerFunc(s.syncHandler))

	s.Handler = router

	return s
}

func (n *NotesServer) getNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	w.Header().Set("content-type", jsonContentType)

	note, err := FindNoteByAttribute(notes, "")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func (n *NotesServer) syncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&notes)
	if err != nil {
		panic(err)
	}
}
