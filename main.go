package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func errHandle(err error) {
	//End execution with a message on errors
	if err != nil {
		log.Fatalln(err)
	}
}

func handleFile(file os.File) {
	for {
		//Make a slice containing bytes from file
		slice := make([]byte, 5)
		int, err := file.Read(slice)

		//Harmlessly return to calling function instead of erroring on EOF
		if err == io.EOF {
			return
		}

		//Print the read bytes to Stdout
		errHandle(err)
		fmt.Print(string(slice[0:int]))
	}
}

func main() {
	//Get arguments from call; assume Stdin if none specified
	args := os.Args[1:]
	if len(args) == 0 {
		args = append(args, "-")
	}

	//Iterate through arguments
	for _, file := range args {
		var toRead *os.File

		//Distinguish between Stdin and File arguments; handle them
		if file == "-" {
			toRead = os.Stdin
		} else {
			var err error
			toRead, err = os.Open(file)
			errHandle(err)
		}
		handleFile(*toRead)
	}
}
