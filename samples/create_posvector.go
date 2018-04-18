package samples

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Opencv_createsamples -info info/info.lst -num 1950 -w 20 -h 20 -vec positives.vec
func CreatePositiveVectorFile(num int) {
	mode := int64(0777)
	if _, err := os.Stat("info"); os.IsNotExist(err) {
		os.MkdirAll("info", os.FileMode(mode))
	}
	cmdName, err := exec.LookPath("opencv_createsamples")
	if err != nil {
		fmt.Println("can't find opencv_createsamples")
		return
	}
	cmdArgs := []string{
		"-info", "info/info.lst",
		"-num", strconv.Itoa(num),
		"-w", "70",
		"-h", "70",
		"-vec", "positives.vec",
	}
	fmt.Println(cmdName, cmdArgs)
	_, err = exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
}
