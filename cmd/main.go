package main

import (
	"bufio"
	notesbot "deni/notesbot"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	notes, err := notesbot.NewNotesFromFS("../knowledge", os.DirFS("../knowledge"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Posts indexed: %d\n", len(notes))

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	input = strings.TrimSpace(input)
	noteIndex, err := strconv.ParseInt(input, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(notes[noteIndex])
}
