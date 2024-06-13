package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var runOsCommandMock func(commandName string, cmdFlags ...string) bool
var cmdRunMock func(c *exec.Cmd) error
var cmdCommandMock func(name string, arg ...string) *exec.Cmd

type execMock struct{}

func (u execMock) command(commandName string, cmdFlags ...string) *exec.Cmd {
	return new(exec.Cmd)
}

func (u execMock) run() error {
	return *new(error)
}

func TestRunOsCommand(t *testing.T) {
	res := runOsCommand("echo", "\"test\"")
	expect := true

	if res != expect {
		t.Errorf("Got %t, expected %t", res, expect)
	}
}

func TestCwdPath(t *testing.T) {
	res := addCwdToFilePath("test")

	currentWorkingDirectory, _ := os.Getwd()
	var expect strings.Builder
	expect.WriteString(currentWorkingDirectory)
	expect.WriteString("/test")

	if res != expect.String() {
		t.Errorf("Got %s, expected %s", res, expect.String())
	}
}

func TestGetFileStats(t *testing.T) {
	res := getFileStats("/bin/bash")

	if res.Name() != "bash" {
		t.Errorf("Unable to correctly get file properties, failed to get bash")
	}

}

func TestFileExists(t *testing.T) {
	res := fileExists("/bin/bash")

	if res != true {
		t.Errorf("Unable to correctly check for file existence, failed to find bash")
	}
}

func TestFileDoesntExist(t *testing.T) {
	res := fileExists("/bin/abcdefghz")

	if res != false {
		t.Errorf("Found file that shouldn't exist")
	}
}

func TestDDShred(t *testing.T) {

}
