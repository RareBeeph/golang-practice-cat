package main

import (
	"fmt"
	"log"
	"os"
)

func errHandle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	for _, file := range os.Args[1:] {
		text, err := os.ReadFile(file)
		errHandle(err)
		fmt.Println(string(text))
	}
}
