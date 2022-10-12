package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	exitCode := m.Run()

	os.Exit(exitCode)
}

var inputData = "Your test input data here\n"
var goModOutput = readFile("../go.mod")
var errorOutput = readFile("nonexistent") + "\n"
var flagtests = []struct {
	inputArgs []string
	output    string
}{
	{[]string{}, inputData},
	{[]string{"-"}, inputData},
	{[]string{"../go.mod"}, goModOutput},
	{[]string{"nonexistent"}, errorOutput},
	{[]string{"../go.mod", "-"}, goModOutput + inputData},
	{[]string{"-", "../go.mod"}, inputData + goModOutput},
	{[]string{"-", "nonexistent"}, inputData + errorOutput},
	{[]string{"nonexistent", "-"}, errorOutput + inputData},
	{[]string{"../go.mod", "nonexistent"}, goModOutput + errorOutput},
	{[]string{"nonexistent", "../go.mod"}, errorOutput + goModOutput},
}

func TestFlagArgs(t *testing.T) {
	var output []byte
	var bytes []byte

	stdin, errIn := os.CreateTemp("", "testStdin")
	handleErr(errIn)
	defer os.Remove(stdin.Name())
	stdout, errOut := os.CreateTemp("", "testStdout")
	handleErr(errOut)
	defer os.Remove(stdout.Name())

	stdin.Write([]byte(inputData))

	os.Stdin = stdin
	os.Stdout = stdout

	for i, flag := range flagtests {
		os.Stdin.Seek(0, 0)
		os.Stdout.Seek(0, 0)
		err := os.Truncate(os.Stdout.Name(), 0)
		errHandle(err)

		rootCmd.SetArgs(flag.inputArgs)
		rootCmd.Execute()
		bytesRead, err := io.ReadAll(os.Stdout)
		handleErr(err)
		output = trimWhitespace([]byte(flag.output))
		bytes = trimWhitespace(bytesRead)

		log.Printf("%v", flag.inputArgs)
		if string(bytes) != string(output) {
			log.Printf("%v: Match!", i)
		} else {
			log.Println(output)
			log.Println(bytes)
			log.Printf("%v: No match. Expected %s, found %s", i, output, bytes)
		}
	}
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func readFile(file string) string {
	toRead, err := os.Open(file)
	if err != nil {
		return err.Error()
	}
	defer toRead.Close()

	bytes, errRead := io.ReadAll(toRead)
	errHandle(errRead)

	return string(bytes)
}

func trimWhitespace(s []byte) []byte {
	return bytes.Trim(s, "\x00\n")
}
