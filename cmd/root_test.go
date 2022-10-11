package cmd

import (
	"io"
	"log"
	"os"
	"testing"
)

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func TestMain(m *testing.M) {

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestStdin(t *testing.T) {
	stdin, errIn := os.CreateTemp("", "testStdin")
	handleErr(errIn)
	defer os.Remove(stdin.Name())
	stdout, errOut := os.CreateTemp("", "testStdout")
	handleErr(errOut)
	defer os.Remove(stdout.Name())

	stdin.Write([]byte(inputData))
	stdin.Seek(0, 0)

	os.Stdin = stdin
	os.Stdout = stdout

	rootCmd.Execute()
	os.Stdout.Seek(0, 0)
	bytesRead, err := io.ReadAll(os.Stdout)
	handleErr(err)
	log.Println(string(bytesRead))
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

var inputData = "Your test input data here\n"
var goModOutput = readFile("../go.mod")
var errorOutput = readFile("nonexistent")
var flagtests = []struct {
	inputArgs []string
	output    string
}{
	{[]string{""}, inputData},
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
	stdin, errIn := os.CreateTemp("", "testStdin")
	handleErr(errIn)
	defer os.Remove(stdin.Name())
	stdout, errOut := os.CreateTemp("", "testStdout")
	handleErr(errOut)
	defer os.Remove(stdout.Name())

	stdin.Write([]byte(inputData))
	stdin.Seek(0, 0)

	os.Stdin = stdin
	os.Stdout = stdout

	for i, flag := range flagtests {
		rootCmd.Args(rootCmd, []string{})
		rootCmd.SetArgs(flag.inputArgs)
		rootCmd.Execute()
		os.Stdout.Seek(0, 0)
		bytesRead, err := io.ReadAll(os.Stdout)
		handleErr(err)
		log.Printf("%v", flag.inputArgs)
		if string(bytesRead) == flag.output {
			log.Printf("%v: Match!", i)
		} else {
			log.Printf("%v: No match. Expected %s, found %s", i, flag.output, string(bytesRead))
		}
	}
}

func TestFileIn(t *testing.T) {

	inputArgs := []string{"../go.mod", "-", "../go.mod", "error", "../go.mod"}

	stdin, errIn := os.CreateTemp("", "testStdin")
	handleErr(errIn)
	defer os.Remove(stdin.Name())
	stdout, errOut := os.CreateTemp("", "testStdout")
	handleErr(errOut)
	defer os.Remove(stdout.Name())

	stdin.Write([]byte(inputData))
	stdin.Seek(0, 0)

	os.Stdin = stdin
	os.Stdout = stdout

	rootCmd.SetArgs(inputArgs)

	rootCmd.Execute()
	os.Stdout.Seek(0, 0)
	bytesRead, err := io.ReadAll(os.Stdout)
	handleErr(err)
	log.Println(string(bytesRead))
}
