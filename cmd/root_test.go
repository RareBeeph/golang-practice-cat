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

// Specify arguments and the expected outputs of rootCmd on them
var errorOutput = readFile("nonexistent") + "\n"
var goModOutput = readFile("../go.mod")
var inputData = "Your test input data here\n"
var flagtests = []struct {
	inputArgs []string
	output    string
}{
	{[]string{}, inputData},                                           //Tests no args -> Stdin
	{[]string{"-"}, inputData},                                        //Tests "-" -> Stdin
	{[]string{"../go.mod"}, goModOutput},                              //Tests reading files
	{[]string{"nonexistent"}, errorOutput},                            //Tests reporting file open errors without exiting
	{[]string{"../go.mod", "-"}, goModOutput + inputData},             //Tests file to Stdin transition
	{[]string{"-", "../go.mod"}, inputData + goModOutput},             //Tests Stdin to file transition
	{[]string{"-", "nonexistent"}, inputData + errorOutput},           //Tests Stdin to error transition
	{[]string{"nonexistent", "-"}, errorOutput + inputData},           //Tests error to Stdin transition
	{[]string{"../go.mod", "nonexistent"}, goModOutput + errorOutput}, //Tests file to error transition
	{[]string{"nonexistent", "../go.mod"}, errorOutput + goModOutput}, //Tests error to file transition
}

func TestFlagArgs(t *testing.T) {
	var output []byte
	var bytes []byte

	// Create temporary files to read/write Stdin and Stdout
	stdin, errIn := os.CreateTemp("", "testStdin")
	handleErr(errIn)
	defer os.Remove(stdin.Name())
	stdout, errOut := os.CreateTemp("", "testStdout")
	handleErr(errOut)
	defer os.Remove(stdout.Name())

	os.Stdin = stdin
	os.Stdout = stdout
	stdin.Write([]byte(inputData))

	for i, flag := range flagtests {
		// Reset initial conditions between each test
		err := os.Truncate(os.Stdout.Name(), 0)
		errHandle(err)
		os.Stdin.Seek(0, 0)

		rootCmd.SetArgs(flag.inputArgs)
		rootCmd.Execute()
		os.Stdout.Seek(0, 0)

		bytesRead, err := io.ReadAll(os.Stdout)
		handleErr(err)
		output = trimWhitespace([]byte(flag.output))
		bytes = trimWhitespace(bytesRead)

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
