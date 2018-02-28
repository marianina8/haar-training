package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
)

/*
 *
 *  Name: store_images.go
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
 *      . height (int) - resize to height (if 0 preserve height)
 *      . width (int) - resize to width (if 0 preserve width)
 *
 */

func storeImages(links []string, folderName string, grayscale bool, start, limit, height, width int) {
	picNum := start
	for _, imageLink := range links {
		// create a request to image-net
		req, err := http.NewRequest(http.MethodGet, imageLink, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		// create folder if does not exist
		mode := int64(0666)
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.Mkdir(folderName, os.FileMode(mode))
		}
		// scan through each line of the response body
		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			if picNum > limit {
				return
			}
			// get the image associated with the link
			resp, e := http.Get(scanner.Text())
			if e != nil || resp.StatusCode != http.StatusOK {
				continue
			}
			defer resp.Body.Close()
			//open a file for writing
			filePath := filepath.Join(folderName, strconv.Itoa(picNum)+".jpg")
			file, err := os.Create(filePath)
			if err != nil {
				continue
			}
			// Use io.Copy to just dump the response body to the file. This supports huge files
			n, err := io.Copy(file, resp.Body)
			if err != nil || n < 3000 {
				_ = os.Remove(filePath)
				continue
			}
			file.Close()
			// open the file for image manipulation
			srcImg, err := imaging.Open(filePath)
			if srcImg == nil || err != nil {
				continue
			}
			// resize image to 100x100
			if height > 0 || width > 0 {
				srcImg = imaging.Resize(srcImg, width, height, imaging.Lanczos)
			}
			// convert image to grayscale
			if grayscale {
				srcImg = imaging.Grayscale(srcImg)
			}
			err = imaging.Save(srcImg, filePath)
			if err != nil {
				continue
			}
			picNum++
			// fmt.Println("Saved", filePath)
		}
	}

}
