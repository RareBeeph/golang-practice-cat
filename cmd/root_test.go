package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Create temporary files to read/write Stdin and Stdout
	stdin, _ := os.CreateTemp("", "testStdin")
	defer os.Remove(stdin.Name())
	stdout, _ := os.CreateTemp("", "testStdout")
	defer os.Remove(stdout.Name())

	os.Stdin = stdin
	os.Stdout = stdout
	stdin.Write([]byte(inputData))

	log.SetFlags(0)

	exitCode := m.Run()

	os.Exit(exitCode)
}

// Specify arguments and the expected outputs of rootCmd on them
func fileToString(file string) string {
	toRead, err := os.Open(file)
	if err != nil {
		return err.Error()
	}
	defer toRead.Close()

	bytes, _ := io.ReadAll(toRead)

	return string(bytes)
}

const inputData = "Your test input data here\n"

var errorOutput = fileToString("nonexistent") + "\n"
var goModOutput = fileToString("../go.mod")
var flagtests = []struct {
	inputArgs []string
	output    string
}{
	{[]string{}, inputData},                                           // Tests no args -> Stdin
	{[]string{"-"}, inputData},                                        // Tests "-" -> Stdin
	{[]string{"../go.mod"}, goModOutput},                              // Tests reading files
	{[]string{"nonexistent"}, errorOutput},                            // Tests reporting file open errors without exiting
	{[]string{"../go.mod", "-"}, goModOutput + inputData},             // Tests file to Stdin transition
	{[]string{"-", "../go.mod"}, inputData + goModOutput},             // Tests Stdin to file transition
	{[]string{"-", "nonexistent"}, inputData + errorOutput},           // Tests Stdin to error transition
	{[]string{"nonexistent", "-"}, errorOutput + inputData},           // Tests error to Stdin transition
	{[]string{"../go.mod", "nonexistent"}, goModOutput + errorOutput}, // Tests file to error transition
	{[]string{"nonexistent", "../go.mod"}, errorOutput + goModOutput}, // Tests error to file transition
}

func trimWhitespace(s []byte) []byte {
	return bytes.Trim(s, "\x00\n")
}

func TestFlagArgs(t *testing.T) {
	var output []byte
	var bytes []byte

	for i, flag := range flagtests {
		// Reset initial conditions between each test
		log.SetOutput(os.Stdout)

		os.Truncate(os.Stdout.Name(), 0)
		os.Stdin.Seek(0, 0)

		rootCmd.SetArgs(flag.inputArgs)
		rootCmd.Execute()
		os.Stdout.Seek(0, 0)

		bytesRead, _ := io.ReadAll(os.Stdout)
		output = trimWhitespace([]byte(flag.output))
		bytes = trimWhitespace(bytesRead)

		log.SetOutput(os.Stderr)
		log.Printf("%v", flag.inputArgs)
		if string(bytes) == string(output) {
			log.Printf("%v: Match!", i)
		} else {
			log.Println(output)
			log.Println(bytes)
			log.Printf("%v: No match. Expected %s, found %s", i, output, bytes)
		}
	}
}
