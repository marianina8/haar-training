package main

import (
	"math/rand"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	links := []string{"http://www.image-net.org/api/text/imagenet.synset.geturls?wnid=n12102133"}
	storeImages(links, "negatives", true, 1, 7000, 24, 24)
	// generate descriptor files here named positives.txt and negatives.txt
	// -w -h should be the width and height of your positive images
	createSampleCmdOptions := "-bgcolor 0 -bgthresh 0 -maxxangle 1.1 -maxyangle 1.1 maxzangle 0.5 -maxidev 40 -w 24 -h 24"
	createSamples("positives.txt", "negatives.txt", "samples", 7000, createSampleCmdOptions)
	// merge vec files
	mergeVecFiles("samples", "samples.vec")
}
