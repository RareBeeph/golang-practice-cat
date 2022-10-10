/*
Copyright © 2022 Keiran Jensen <keiranjensen@gmail.com>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bphcat",
	Short: "A Go language clone of cat",
	Long: `This is a clone of the shell builtin "cat",
	which takes any number of files as arguments
	and prints them to Stdout. It additionally can read
	from Stdin if there are no arguments, or if one is "-".
	
	Example: ./bphcat main.go go.mod`,
	Run: func(cmd *cobra.Command, args []string) {
		var toRead *os.File
		var err error

		// Prints bytes as strings from files, or from Stdin when argument is empty or "-", to Stdout
		if len(args) == 0 {
			args = append(args, "-")
		}

		for _, fileArg := range args {
			if fileArg == "-" {
				toRead = os.Stdin
			} else {
				toRead, err = os.Open(fileArg)
				errHandle(err)
				defer toRead.Close()
			}
			printFileContents(toRead)
		}
	},
}

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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
