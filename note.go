package notesbot

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
	tagsSeparator        = "Tags: "
)

type Note struct {
	Title string
	Tags  []string
	Body  string
}

func newNote(noteBody io.Reader) (Note, error) {
	scanner := bufio.NewScanner(noteBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	return Note{
		Title: readMetaLine(titleSeparator),
		Tags:  strings.Split(readMetaLine(tagsSeparator), ", "),
		Body:  readBody(scanner),
	}, nil
}

func readBody(scanner *bufio.Scanner) string {
	scanner.Scan() // ignore a line
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}

	return strings.TrimSuffix(buf.String(), "\n")
}
