package notesbot

import (
	"encoding/json"
	"net/http"
)

const jsonContentType = "application/json"

type Store interface {
	StoreNotes(notes []Note)
	AllNotes() []Note
}

type NotesServer struct {
	store Store
	http.Handler
}

func NewServer(store Store) *NotesServer {
	s := new(NotesServer)
	s.store = store

	router := http.NewServeMux()
	router.Handle("/sync-notes", http.HandlerFunc(s.syncHandler))
	router.Handle("/bot", http.HandlerFunc(s.botHandler))

	s.Handler = router

	return s
}

func (n *NotesServer) SendNotes() {

}

func (n *NotesServer) getNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	w.Header().Set("content-type", jsonContentType)

	note, err := FindNoteByAttribute(n.store.AllNotes(), "")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func (s *NotesServer) syncHandler(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	if r.Method != http.MethodPost {
		return
	}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&notes)
	if err != nil {
		panic(err)
	}

	s.store.StoreNotes(notes)
}
