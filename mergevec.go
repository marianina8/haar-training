package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lunixbochs/struc"
)

/*
 *
 *  Name: mergevec.go
 *  Author: Marian Montagnino
 *  Description: mergeVecFiles is a function that merges .vec files.
 *  I made it as a replacement for mergevec.cpp (created by Naotoshi Seo) in order to avoid recompiling openCV with mergevec.cpp.
 *  Parameters:
 *      . vecDir (string) - single directory containing all .vec files you need to merge
 *      . outputFilename (string) - the name of the output file
 *
 */

// imgData file header
type imgData struct {
	TotalImages  int32
	SizeOfImages int32
	Min          int16
	Max          int16
}

// mergeVecFiles takes a directory of .vec files and merges them and combines them to a single .vec file
// 	(1) Iterates through files getting a count of the total images in the .vec files
// 	(2) checks that the image sizes in all files are the same
// 	The format of a .vec file is:
// 	4 bytes denoting number of total images (int)
// 	4 bytes denoting size of images (int)
// 	2 bytes denoting min value (short)
// 	2 bytes denoting max value (short)
// 	ex: 	6400 0000 4605 0000 0000 0000
// 		hex		6400 0000  	4605 0000 		0000 		0000
// 			   	# images  	size of h * w		min		max
// 		dec	    	100     	1350			0 		0
func mergeVecFiles(vecDir, outputFilename string) {
	// Check that the vecDir does not end in '/' and if it does, remove it.
	if vecDir[len(vecDir)-1:] == "/" {
		vecDir = vecDir[:len(vecDir)-1]
	}
	// Get .vec files
	files, err := filepath.Glob(filepath.Join(vecDir) + "/*.vec")
	if err != nil {
		fmt.Printf("Encountered an error (%s) globbing"+filepath.Join(vecDir)+"/*.vec\n", err.Error())
		return
	}
	// Check to make sure there are .vec files in the directory
	if len(files) <= 0 {
		fmt.Printf("Vec files to be merged could not be found from directory: %s\n", vecDir)
		return
	}
	// Check to make sure there are more than one .vec files
	if len(files) == 1 {
		fmt.Printf("Only 1 vec file was found in directory: %s. Cannot merge a single file.\n", vecDir)
		return
	}
	// Get the value for the first image size
	fp, err := os.Open(files[0])
	if err != nil {
		fmt.Println("err opening file:", err)
		return
	}
	defer fp.Close()
	// Get first image's file header
	n := 12
	data := make([]byte, n)
	_, err = fp.Read(data)
	if err != nil {
		fmt.Println("err reading 12 bytes")
		return
	}
	imgHeaderData := &imgData{}
	err = binary.Read(bytes.NewBuffer(data), binary.LittleEndian, imgHeaderData)
	if err != nil {
		fmt.Println("err unpacking binary data", err)
		return
	}
	// Get image size
	prevImageSize := imgHeaderData.SizeOfImages
	// Get the total number of images
	totalNumImages := int32(0)
	// Loop through the vector files
	for _, f := range files {
		fp, _ := os.Open(f)
		defer fp.Close()
		n := 12
		data := make([]byte, n)
		_, err = fp.Read(data)
		if err != nil {
			fmt.Println("err reading header bytes", err)
			return
		}
		imgHeaderData := &imgData{}
		err = binary.Read(bytes.NewBuffer(data), binary.LittleEndian, imgHeaderData)
		if err != nil {
			fmt.Println("err unpacking binary data", err)
			return
		}
		fmt.Println("imgHeaderData.SizeOfImages", imgHeaderData.SizeOfImages)
		if imgHeaderData.SizeOfImages != prevImageSize {
			fmt.Printf("the image sizes in the .vec files differ. These values must be the same. \n The image size of file %s: %d\n The image size of previous files: %d", f, imgHeaderData.SizeOfImages, prevImageSize)
			return
		}
		// Calculate the total number of images
		fmt.Println("imgHeaderData.TotalImages:", imgHeaderData.TotalImages)
		totalNumImages += imgHeaderData.TotalImages
	}
	fmt.Println("totalNumImages:", totalNumImages)
	//  Iterate through the .vec files, writing their data (not the header) to the output file
	var buf bytes.Buffer
	imgHeaderData = &imgData{totalNumImages, prevImageSize, 0, 0}
	err = struc.Pack(&buf, imgHeaderData)
	if err != nil {
		fmt.Println("error packing header, err")
	}
	// create the output file
	outFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println("error creating output file", err)
	}
	defer outFile.Close()
	// write the image header data
	_, err = outFile.Write(buf.Bytes())
	if err != nil {
		fmt.Println("error writing header", err)
	}
	// write all the data of the vector files (omitting their header file)
	for _, fp := range files {
		data, err := ioutil.ReadFile(fp)
		if err != nil {
			fmt.Println("error reading file", fp, "err:", err)
			return
		}
		// omit header when writing to the output file
		_, err = outFile.Write(data[12:])
		if err != nil {
			fmt.Println("error writing header, err")
			return
		}
	}
	return
}
