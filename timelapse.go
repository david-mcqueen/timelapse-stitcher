package main

import (
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	"io/ioutil"
	"strings"
)

func main(){
	var (
		dir = flag.String("dir", "./images", "Directory to find the images in")
		fps = flag.Uint("fps", 10, "Framerate")
		outputFile = flag.String("outputFile", "timelapse.avi", "Path and name of output file.")
		limit = flag.Uint("limit", 30, "Maximum number of images to stitch")
	)
	stitchImages(*dir, *fps, *outputFile, *limit)
}

func stitchImages(dir string, fps uint, outputFile string, limit uint) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Errorf("Failed to load images in directory %s", dir)
	}

	fmt.Println("Found images: %s", len(files))

	newImage := gocv.NewMat()
	defer newImage.Close()

	firstFileName := files[0].Name()
	firstImage := dir + "/" + firstFileName

	newImage = gocv.IMRead(firstImage, gocv.IMReadAnyColor)

	writer, err := gocv.VideoWriterFile(outputFile, "MJPG", float64(fps), newImage.Cols(), newImage.Rows(), true)

	if err != nil {
		fmt.Errorf("Failed to creater writer")
	}

	defer writer.Close()

	fmt.Println("Video being generated with the name :", outputFile)


	for _, image := range files[:limit] {
		if image.IsDir() || !strings.HasSuffix(image.Name(), "jpg") {
			continue
		}

		newImage = gocv.IMRead(dir + "/" + image.Name(), gocv.IMReadAnyColor)
		err = writer.Write(newImage)

		if err != nil {
			fmt.Errorf("Failed to add image to video file ")
		}
	}
}
