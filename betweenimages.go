package main

import (
	// "fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	// "math"
	"github.com/icza/mjpeg"
	// "io/ioutil"
	"bytes"
	"os"
	"strconv"
)

const magicNumber = 64

func main() {

	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// Video size: 300x300 pixels, FPS: 24
	aw, err := mjpeg.New("test.avi", 300, 300, 48)
	checkErr(err)

	for n := 1; n < 16; n++ {
		imgA, imgB := getImagesFromFiles(n)

		// Move into floating point. Create two arrays of floats that represent the images.
		i := 0
		var imgAf [360000]float32
		for i < imgA.Stride*imgA.Rect.Max.Y {
			imgAf[i] = float32(imgA.Pix[i])
			imgAf[i+1] = float32(imgA.Pix[i+1])
			imgAf[i+2] = float32(imgA.Pix[i+2])
			imgAf[i+3] = float32(imgA.Pix[i+3])
			i += 4
		}

		var imgBf [360000]float32
		i = 0
		for i < imgB.Stride*imgB.Rect.Max.Y {
			imgBf[i] = float32(imgB.Pix[i])
			imgBf[i+1] = float32(imgB.Pix[i+1])
			imgBf[i+2] = float32(imgB.Pix[i+2])
			imgBf[i+3] = float32(imgB.Pix[i+3])
			i += 4
		}

		addImages(aw, imgAf, imgBf)
	}

	checkErr(aw.Close())

}

func getImagesFromFiles(n int) (*image.RGBA, *image.RGBA) {

	var stA string
	if n < 10 {
		stA = "inputs/0"
	} else {
		stA = "inputs/"
	}
	fn := stA + strconv.Itoa(n) + ".png"
	existingImageFile, err := os.Open(fn)
	if err != nil {

	}
	defer existingImageFile.Close()

	inputImgA, _, err := image.Decode(existingImageFile)
	if err != nil {

	}

	imgA := image.NewRGBA(image.Rect(0, 0, 300, 300))
	boundsA := inputImgA.Bounds()
	draw.Draw(imgA, imgA.Rect, inputImgA, boundsA.Min, draw.Over)

	var stB string
	if n < 9 {
		stB = "inputs/0"
	} else {
		stB = "inputs/"
	}
	fn = stB + strconv.Itoa(n+1) + ".png"
	existingImageFile, err = os.Open(fn)
	if err != nil {

	}
	defer existingImageFile.Close()

	inputImgB, _, err := image.Decode(existingImageFile)
	if err != nil {

	}

	imgB := image.NewRGBA(image.Rect(0, 0, 300, 300))
	boundsB := inputImgB.Bounds()
	draw.Draw(imgB, imgB.Rect, inputImgB, boundsB.Min, draw.Over)

	return imgA, imgB
}

func addImages(aw mjpeg.AviWriter, x [360000]float32, y [360000]float32) {

	var fn float32
	for fn = 1; fn < magicNumber+1; fn++ {

		var xPart float32
		if fn == magicNumber {
			xPart = 0
		} else {
			xPart = (magicNumber - fn) / magicNumber
		}
		yPart := 1 - xPart

		// Make an image from the two arrays and write it to a file
		img := image.NewRGBA(image.Rect(0, 0, 300, 300))
		j := 0
		for j < (img.Stride * img.Rect.Max.Y) {
			img.Pix[j] = uint8(x[j]*xPart + y[j]*yPart)
			img.Pix[j+1] = uint8(x[j+1]*xPart + y[j+1]*yPart)
			img.Pix[j+2] = uint8(x[j+2]*xPart + y[j+2]*yPart)
			img.Pix[j+3] = uint8(x[j+3]*xPart + y[j+3]*yPart)
			j += 4
		}

		checkErr := func(err error) {
			if err != nil {
				panic(err)
			}
		}

		buf := &bytes.Buffer{}
		checkErr(jpeg.Encode(buf, img, nil))
		checkErr(aw.AddFrame(buf.Bytes()))

	}

}
