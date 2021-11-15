package notesbot

import (
	"bufio"
	"bytes"
	"crypto/rand"
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
	Title      string
	Tags       []string
	Body       string
	TelegramId int64
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

func (n *Note) GetTelegramId() int64 {
	if 0 != n.TelegramId {
		return n.TelegramId
	}
	id, _ := rand.Prime(rand.Reader, 10)
	n.TelegramId = id.Int64()

	return n.TelegramId
}
