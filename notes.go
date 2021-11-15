package notesbot

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func NewNotesFromFS(rootPath string, fileSystem fs.FS) ([]Note, error) {
	var notes []Note

	var index int
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			index++
			if err != nil {
				return err
			}

			note, _ := getNote(fileSystem, path)
			note.Title = filenameWithoutExtension(info.Name())
			note.TelegramId = int64(index)
			notes = append(notes, note)

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return notes, nil
}

func getNote(fileSystem fs.FS, fileName string) (Note, error) {
	postFile, err := os.Open(fileName)
	if err != nil {
		return Note{}, err
	}
	defer postFile.Close()

	return newNote(postFile)
}

func FindNoteByAttribute(notes []Note, needle int64) (Note, error) {
	n := Note{}
	for _, n := range notes {
		if needle == n.TelegramId {
			return n, nil
		}
	}

	return n, errors.New("Notes not found")
}

func AllNotesTitleList(notes []Note) []string {
	var index []string
	for _, n := range notes {
		index = append(index, n.Title)
	}

	return index
}

func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}
