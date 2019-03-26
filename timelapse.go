package main

import (
	"flag"
	"fmt"
	"gocv.io/x/gocv"
	image2 "image"
	"io/ioutil"
	"runtime"
	"strings"
)

func main(){
	var (
		dir = flag.String("dir", "./images", "Directory to find the images in")
		fps = flag.Uint("fps", 10, "Framerate")
		outputFile = flag.String("outputFile", "timelapse.avi", "Path and name of output file.")
		limit = flag.Uint("limit", 3000, "Maximum number of images to stitch")
	)
	flag.Parse()
	runtime.GOMAXPROCS(2)

	stitchImages(*dir, *fps, *outputFile, *limit)
}

func stitchImages(dir string, fps uint, outputFile string, limit uint) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Errorf("Failed to load images in directory %s", dir)
	}

	fmt.Println("Found images: %s", len(files))

	if limit > uint(len(files)) {
		limit = uint(len(files))
	}


	newImage := gocv.NewMat()
	defer newImage.Close()

	firstFileName := files[0].Name()
	firstImage := dir + "/" + firstFileName

	newImage = gocv.IMRead(firstImage, gocv.IMReadAnyColor)

	imageWidth, imageHeight := newImage.Cols(), newImage.Rows()


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

		resizedImage := gocv.NewMat()

		var pointy image2.Point

		pointy.X = imageWidth
		pointy.Y = imageHeight


		gocv.Resize(newImage, &resizedImage, pointy, 0, 0, gocv.InterpolationLinear)


		err = writer.Write(resizedImage)

		if err != nil {
			fmt.Errorf("Failed to add image to video file ")
		}
	}
}
