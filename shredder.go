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
var globalFileObjStat (os.FileInfo) = nil
var globalFileObjPtr (*os.File) = nil
var globalFilePath string = ""
var globalFileSize int64 = 0

// func execute(file string, iterations int) {
// 	currentWorkingDirectory, err := os.Getwd()
// 	if err != nil {
// 		return
// 	}
// 	fileName := currentWorkingDirectory + "/" + file
// 	commandOptions := "-vun" + strconv.Itoa(iterations) + " " + fileName
// 	fmt.Println("Running command: " + commandName + " " + commandOptions)
// 	// cmd := exec.Command("bash", "-c", commandName, commandOptions)

// }

func runOsCommand(commandName string, cmdFlags ...string) bool {
	var commandOptions strings.Builder
	for _, lstr := range cmdFlags {
		commandOptions.WriteString(lstr)
	}
	cmd := exec.Command(commandName, commandOptions.String())
	err := *new(error)
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
	if err != nil {
		return false
	} else {
		return true
	}
}

func addCwdToFilePath(file string) string {
	fileName := ""
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fileName = currentWorkingDirectory + "/" + file
	return fileName
}

func getFileStats(file string) (fp os.FileInfo) {
	fileObjStat, err := os.Stat(file)
	if err != nil {
		panic(err)
	}
	return fileObjStat
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

func getFileObject(file string) (fp *os.File) {
	fileObjOpen, err := os.OpenFile(file, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
	return fileObjOpen
}

func getFileSize() int64 {
	return globalFileSize
}

func execute_basic_shred(iterations int) bool {
	array := make([]byte, getFileSize())
	for iter := 0; iter < iterations; iter++ {
		rand.Read(array)
		_, err := globalFileObjPtr.WriteAt(array, 0)
		if err != nil {
			panic(err)
		}
		globalFileObjPtr.Sync()
	}
	return true
}

func execute_dd_shred(iterations int) bool {
	for iter := 0; iter < iterations; iter++ {
		if !runOsCommand("dd", "status=progress", "bs=", strconv.FormatInt(globalFileSize, 10), "count=1", "if=/dev/urandom", "of=", globalFilePath) {
			return false
		}
	}
	return true
}

func execute_wipe_shred(iterations int) bool {
	return runOsCommand("wipe", "-f", "-Q", strconv.FormatInt(int64(iterations), 10), globalFilePath)
}

func validateFileName(file string) bool {
	if strings.Contains(file, "/") == false {
		globalFilePath = addCwdToFilePath(file)

	} else {
		globalFilePath = file
	}

	if !fileExists(globalFilePath) {
		return false
	}
	return true
}

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("Can't Execute this on a windows machine")
	} else {

		fmt.Println("DISCLAIMER: this program is a GoLang programming exercise.")
		fmt.Println("            Data may still be available when using the right instruments")
		fmt.Println("")
		fmt.Println("")
		fileStr := flag.String("file", "", "filename of the file to be shredded")
		iterPtr := flag.Int("iter", 3, "number of iterations of shredding")
		fileRmv := flag.Bool("rem", false, "remove file after shredding")
		flagAlgo := flag.Int("algo", 0, `shredding algorithms
0 - basic writing of random bytes\n
1 - shredding using dd and random bytes\n
2 - shredding using wipe\n`)
		flag.Usage = func() {
			fmt.Println("Usage of ", os.Args[0], ":")
			fmt.Println("   ", os.Args[0], " -iter 3 -file myfile -algo 2")
			flag.PrintDefaults()
		}
		flag.Parse()

		if flag.NFlag() == 0 {
			flag.Usage()
			os.Exit(0)
		}
		fmt.Println("filename:", *fileStr)
		fmt.Println("iterations:", *iterPtr)

		if len(*fileStr) == 0 {
			fmt.Println("Empty file name, terminating")
			os.Exit(0)
		}

		if !validateFileName(*fileStr) {
			fmt.Println("Can't find file")
			os.Exit(1)
		}

		globalFileObjStat = getFileStats(globalFilePath)
		globalFileSize = globalFileObjStat.Size()
		globalFileObjPtr = getFileObject(globalFilePath)

		defer func() {
			if err := globalFileObjPtr.Close(); err != nil {
				panic(err)
			}
		}()

		switch *flagAlgo {
		case 0:
			execute_basic_shred(*iterPtr)
		case 1:
			execute_dd_shred(*iterPtr)
		case 2:
			execute_wipe_shred(*iterPtr)
		}

		if fileRmv != nil {
			os.Remove(globalFilePath)
		}
	}
}
