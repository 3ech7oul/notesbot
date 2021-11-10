package notesbot

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type FileSystemStore struct {
	database *json.Encoder
	Notes    []Note
}

func NewFileSystemStore(file *os.File) (*FileSystemStore, error) {

	err := initialiseDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	notes, err := LoadNotes(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemStore{
		database: json.NewEncoder(&tape{file}),
		Notes:    notes,
	}, nil
}

func initialiseDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

func LoadNotes(rdr io.Reader) ([]Note, error) {
	var notes []Note
	err := json.NewDecoder(rdr).Decode(&notes)

	if err != nil {
		err = fmt.Errorf("problem parsing notes, %v", err)
	}

	return notes, err
}
