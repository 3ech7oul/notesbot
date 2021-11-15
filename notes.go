package notesbot

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func NewNotesFromFS(rootPath string, fileSystem fs.FS) ([]Note, error) {
	var notes []Note

	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			note, _ := getNote(fileSystem, path)
			note.Title = filenameWithoutExtension(info.Name())
			notes = append(notes, note)

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	for _, n := range notes {
		n.TelegramId = n.GetTelegramId()
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

func FindNoteByAttribute(notes []Note, needle string) (Note, error) {
	n := Note{}
	for _, n := range notes {
		if needle == strconv.FormatInt(n.TelegramId, 16) {
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
