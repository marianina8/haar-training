package samples

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

/*
// createSamples creates positive samples from an image and applies distortions repeatedly
// createSamplesCmdOptions: "-bgcolor 0 -bgthresh 0 -maxxangle 1.1 -maxyangle 1.1 maxzangle 0.5 -maxidev 40 -w 20 -h 20"
// Opencv_createsamples -img watch5050.jpg -bg.txt -info info/info.lst -pngoutput info -maxxangle 0.5 -maxyangle -0.5 -maxzangle 0.5 -num 1950
func createSamples(posFile string, negFile string, vecOutputDir string, totalNum int, createSampleCmdOptions string) {
	// create vecOutputDir if it doesn't exist
	mode := int64(0777)
	if _, err := os.Stat(vecOutputDir); os.IsNotExist(err) {
		os.Mkdir(vecOutputDir, os.FileMode(mode))
	}
	// grab file paths from positives file
	fileHandle, _ := os.Open(posFile)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	posFiles := []string{}
	for fileScanner.Scan() {
		posFiles = append(posFiles, fileScanner.Text())
	}
	// grab file paths from negatives file
	fileHandle, _ = os.Open(negFile)
	defer fileHandle.Close()
	fileScanner = bufio.NewScanner(fileHandle)
	negFiles := []string{}
	for fileScanner.Scan() {
		negFiles = append(negFiles, fileScanner.Text())
	}
	// number of generated images from one image so that total will be totalNum
	num := int(round(float64(totalNum)/float64(len(posFiles)), .5, 0))
	numfloor := int(totalNum / len(posFiles))
	numremain := totalNum - numfloor*len(posFiles)

	runDatas := []runData{}
	k := 0
	for k < len(posFiles) {
		img := strings.TrimSpace(posFiles[k])
		// my $num = ($k < $numremain) ? $numfloor + 1 : $numfloor;
		if k < numremain {
			num = numfloor + 1
		} else {
			num = numfloor
		}
		// Pick up negative images randomly and write to tmpFile
		localNegatives := []string{}
		i := 0
		for i < num {
			ind := random(0, len(negFiles))
			localNegatives = append(localNegatives, negFiles[ind])
			i++
		}
		fmt.Println("len(localNegatives)=", len(localNegatives))
		// write randomly selected negative files to a temp file (tmpFile)
		tmpFile := fmt.Sprintf("tmp%d.txt", k)
		err := writeLines(localNegatives, tmpFile)
		if err != nil {
			fmt.Println("Cant write lines")
			return
		}
		vecFile := vecOutputDir + "/" + filepath.Base(img) + ".vec"
		// put together whole command
		cmdName, err := exec.LookPath("opencv_createsamples")
		if err != nil {
			fmt.Println("can't find opencv_createsamples")
			return
		}
		cmdArgs := []string{
			createSampleCmdOptions,
			"-img", img,
			"-bg", tmpFile,
			"-vec", vecFile,
			"-num", strconv.Itoa(num),
		}
		_, err = exec.Command(cmdName, cmdArgs...).CombinedOutput()
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		k++
	}
	return
}*/

// Opencv_createsamples -img watch5050.jpg -bg bg.txt -info info/info.lst -pngoutput info -maxxangle 0.5 -maxyangle -0.5 -maxzangle 0.5 -num 1950 -w 20 -h 20
func CreateSamples(posFile string, bgFile string, num int, createSampleCmdOptions string) {
	// if "info" doesn't exist create it
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
		"-img", posFile,
		"-bg", bgFile,
		//	"-info", "info/info.lst",
		"-w", "70",
		"-h", "70",
		createSampleCmdOptions,
		"-vec", "positives.vec",
		"-num", strconv.Itoa(num),
	}
	fmt.Println(cmdName, cmdArgs)
	_, err = exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
}
