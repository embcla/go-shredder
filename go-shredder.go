package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var commandName = "shred"

func execute(file string, iterations int) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return
	}
	fileName := currentWorkingDirectory + "/" + file
	commandOptions := "-vun" + strconv.Itoa(iterations) + " " + fileName
	fmt.Println("Running command: " + commandName + " " + commandOptions)
	// cmd := exec.Command("bash", "-c", commandName, commandOptions)
	cmd := exec.Command(commandName, commandOptions)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	// out, err := exec.Command("ls", "-lrt").Output()
	if err != nil {
		fmt.Println("Result: " + out.String())
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	} else {
		fmt.Println("Command Successfully Executed")
		fmt.Println("Result: " + out.String())
	}

}

func execute_shred(file string, iterations int) {
	fileName := ""
	res := strings.Contains(file, "/")
	if res == false {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			return
		}
		fileName = currentWorkingDirectory + "/" + file
	} else {
		fileName = file
	}

	fileObjStat, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	fileObjOpen, err := os.OpenFile(fileName, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fileObjOpen.Close(); err != nil {
			panic(err)
		}
	}()

	array := make([]byte, fileObjStat.Size())

	for iter := 0; iter < iterations; iter++ {
		rand.Read(array)
		_, err := fileObjOpen.WriteAt(array, 0)
		if err != nil {
			panic(err)
		}
		fileObjOpen.Sync()
	}

	os.Remove(fileName)

}

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a windows machine")
	} else {

		fmt.Println("DISCLAIMER: this program is a GoLang programming exercise.")
		fmt.Println("            It is not a true attempt at a data shredder.")
		fmt.Println("            Do not expect your data to be truly unavailable after usage.")
		filePtr := flag.String("file", "", "filename")
		iterPtr := flag.Int("iter", 3, "an int")

		flag.Parse()

		fmt.Println("filename:", *filePtr)
		fmt.Println("iterations:", *iterPtr)

		if len(*filePtr) == 0 {
			fmt.Println("Empty file name, terminating")
			os.Exit(0)
		}

		execute_shred(*filePtr, *iterPtr)
	}
}
