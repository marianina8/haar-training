package training

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Opencv_traincascade -data data -vec positives.vec -bg bg.txt -numPos 1800 -numNeg 900 -numStages 10 -w 20 -h 20

func HaarCascade(dataFolder string, numPositive, numNegative, numStages int) {
	// if {dataFolder} doesn't exist create it
	mode := int64(0777)
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		os.MkdirAll(dataFolder, os.FileMode(mode))
	}
	cmdName, err := exec.LookPath("opencv_traincascade")
	if err != nil {
		fmt.Println("can't find opencv_traincascade")
		return
	}
	cmdArgs := []string{
		"-data", dataFolder,
		"-vec", "positives.vec",
		"-bg", "bg.txt",
		"-numPos", strconv.Itoa(numPositive),
		"-numNeg", strconv.Itoa(numNegative),
		"-numStages", strconv.Itoa(numStages),
		"-w", "70",
		"-h", "70",
	}
	fmt.Println(cmdName, cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe:", err)
		return
	}
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}
