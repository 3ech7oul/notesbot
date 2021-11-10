package notesbot

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func NewNotesFromFS(rootPath string, ileSystem fs.FS) ([]Note, error) {
	var posts []Note

	posts, _ = readFiles(rootPath, ileSystem, posts)

	return posts, nil
}

func readFiles(rootPath string, fileSystem fs.FS, posts []Note) ([]Note, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}

	for _, f := range dir {
		post, err := getNote(fileSystem, f.Name())

		if f.IsDir() {
			var notesInDir []Note
			str := fmt.Sprintf("./%s/%s", rootPath, f.Name())
			notesInDir, _ = readFiles(rootPath, os.DirFS(str), posts)

			for _, n := range notesInDir {
				posts = append(posts, n)
			}
		}

		if err != nil {
			return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
		}

		post.Title = strings.ReplaceAll(f.Name(), ".md", "")
		posts = append(posts, post)
	}

	return posts, nil
}

func getNote(fileSystem fs.FS, fileName string) (Note, error) {
	postFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Note{}, err
	}
	defer postFile.Close()

	return newNote(postFile)
}

func FindNoteByAttribute(notes []Note, needle string) (Note, error) {
	n := Note{}
	for _, n := range notes {
		if needle == n.Title {
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
