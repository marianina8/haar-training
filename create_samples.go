package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

/*
 *
 *  Name: create_samples.go
 *  Author: Marian Montagnino
 *  Description: storeImages is a function that given specified parameters will download negative images from www.image-net.org
 *  Parameters: takes in links to images from http://www.image-net.org/, folder name to store images and limit (int) the number of downloaded images.
 *      . links ([]string) - links contains a slice of strings representing a single link to image-net containing links to negative images
 *      example link: http://www.image-net.org/api/text/imagenet.synset.geturls?wnid=n12102133
 *      . folderName (string) - folder name to store all negative images (negative images are saved as numbers, ie. 1.jpg, 2.jpg, ...)
 *      . grayscale (bool) - if set to true, converts image to grayscale
 *      . start (int) - number to start saving file.  If function was run earlier and the last negative image saved is 20.jpg, pass in 21 for start so
 *      you do not overwrite any files saved in previous run
 *      . limit (int) - limit number of negative files saved
 *      . height (int) - resize sample images to height (if 0 preserve height)
 *      . width (int) - resize sample images to width (if 0 preserve width)
 *
 */

// createSamples creates positive samples from an image and applies distortions repeatedly
// createSamplesCmdOptions: "-bgcolor 0 -bgthresh 0 -maxxangle 1.1 -maxyangle 1.1 maxzangle 0.5 -maxidev 40 -w 20 -h 20"
func createSamples(posFile string, negFile string, vecOutputDir string, totalNum int, createSampleCmdOptions string) {
	// create vecOutputDir if it doesn't exist
	mode := int64(0666)
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
}
