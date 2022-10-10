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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
