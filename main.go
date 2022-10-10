package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func errHandle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func printFileContents(file *os.File) {
	bytesRead, err := io.ReadAll(file)
	errHandle(err)

	fmt.Print(string(bytesRead))
}

func main() {
	var toRead *os.File
	var err error

	// Prints bytes as strings from files, or from Stdin when argument is empty or "-", to Stdout
	args := os.Args[1:]
	if len(args) == 0 {
		args = append(args, "-")
	}

	for _, file := range args {
		if file == "-" {
			toRead = os.Stdin
		} else {
			toRead, err = os.Open(file)
			errHandle(err)
			defer toRead.Close()
		}

		printFileContents(toRead)
	}
}
