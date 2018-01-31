package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	// "math"
	"github.com/icza/mjpeg"
	"io/ioutil"
	"os"
	"strconv"
)

const magicNumber = 64

func main() {

	// Get two input images from files and make RGBA images from them
	existingImageFile, err := os.Open("A.jpg")
	if err != nil {

	}
	defer existingImageFile.Close()

	inputImgA, _, err := image.Decode(existingImageFile)
	if err != nil {

	}

	imgA := image.NewRGBA(image.Rect(0, 0, 250, 300))
	boundsA := inputImgA.Bounds()
	draw.Draw(imgA, imgA.Rect, inputImgA, boundsA.Min, draw.Over)

	existingImageFile, err = os.Open("B.jpg")
	if err != nil {

	}
	defer existingImageFile.Close()

	inputImgB, _, err := image.Decode(existingImageFile)
	if err != nil {

	}

	imgB := image.NewRGBA(image.Rect(0, 0, 250, 300))
	boundsB := inputImgB.Bounds()
	draw.Draw(imgB, imgB.Rect, inputImgB, boundsB.Min, draw.Over)

	// At this point, we have imgA and imgB, which are our starting points for creating the interpolated images.

	// Move into floating point. Create two arrays of floats that represent the images.
	i := 0
	var imgAf [300000]float32
	for i < imgA.Stride*imgA.Rect.Max.Y {
		imgAf[i] = float32(imgA.Pix[i])
		imgAf[i+1] = float32(imgA.Pix[i+1])
		imgAf[i+2] = float32(imgA.Pix[i+2])
		imgAf[i+3] = float32(imgA.Pix[i+3])
		i += 4
	}

	var imgBf [300000]float32
	i = 0
	for i < imgB.Stride*imgB.Rect.Max.Y {
		imgBf[i] = float32(imgB.Pix[i])
		imgBf[i+1] = float32(imgB.Pix[i+1])
		imgBf[i+2] = float32(imgB.Pix[i+2])
		imgBf[i+3] = float32(imgB.Pix[i+3])
		i += 4
	}

	makeImages(imgAf, imgBf)

	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// Video size: 200x100 pixels, FPS: 24
	aw, err := mjpeg.New("test.avi", 200, 100, 24)
	
	checkErr(err)

	// Create a movie from images: 1.jpg, 2.jpg, ..., 10.jpg
	for j := 0; j < 4; j++ {
		for i := 1; i < magicNumber + 1; i++ {
			data, err := ioutil.ReadFile(fmt.Sprintf("%d.jpg", i))
			checkErr(err)
			checkErr(aw.AddFrame(data))
		}
	}

	checkErr(aw.Close())

}

func makeImages(x [300000]float32, y [300000]float32) {

	var fn float32
	for fn = 1; fn < magicNumber + 1; fn++ {

		var xPart float32
		if fn == magicNumber {
			xPart = 0
		} else {
			xPart = (magicNumber - fn) / magicNumber
		}
		yPart := 1 - xPart

		// Make an image from the two arrays and write it to a file
		img := image.NewRGBA(image.Rect(0, 0, 250, 300))
		j := 0
		for j < (img.Stride * img.Rect.Max.Y) {
			img.Pix[j] = uint8(x[j]*xPart + y[j]*yPart)
			img.Pix[j+1] = uint8(x[j+1]*xPart + y[j+1]*yPart)
			img.Pix[j+2] = uint8(x[j+2]*xPart + y[j+2]*yPart)
			img.Pix[j+3] = uint8(x[j+3]*xPart + y[j+3]*yPart)
			j += 4
		}

		fname := strconv.Itoa(int(fn)) + ".jpg"
		f, _ := os.Create(fname)
		defer f.Close()
		jpeg.Encode(f, img, &jpeg.Options{Quality: 100})

	}

}

// func recursiveMakeImage(x [300000]float32, y [300000]float32, fn int, g int) {
// 	fmt.Println("Called recursiveMakeImage with fn =", fn, "and g =", g)
// 	makeImage(x, y, fn)
// 	g++
// 	// This won't work if the guard is an odd number...
// 	if g < 6 {
// 		fnBack := fn - int((16 / math.Pow(2, float64(g-1))))
// 		fnFwd := fn + int((16 / math.Pow(2, float64(g-1))))
// 		recursiveMakeImage(x, average(x, y), fnBack, g)
// 		recursiveMakeImage(average(x, y), y, fnFwd, g)
// 	}
// }

// func makeImage(x [300000]float32, y [300000]float32, fn int) {

// 	fmt.Println("Called makeImage with fn=", fn)

// 	// Make an image from the average of the two arrays and write it to a file
// 	img := image.NewRGBA(image.Rect(0, 0, 250, 300))
// 	j := 0
// 	for j < (img.Stride * img.Rect.Max.Y) {
// 		img.Pix[j] = uint8(x[j]/2 + y[j]/2)
// 		img.Pix[j+1] = uint8(x[j+1]/2 + y[j+1]/2)
// 		img.Pix[j+2] = uint8(x[j+2]/2 + y[j+2]/2)
// 		img.Pix[j+3] = uint8(x[j+3]/2 + y[j+3]/2)
// 		j += 4
// 	}

// 	fname := strconv.Itoa(fn) + ".jpg"
// 	f, _ := os.Create(fname)
// 	defer f.Close()
// 	jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
// }

// func average(x [300000]float32, y [300000]float32) [300000]float32 {
// 	var z [300000]float32
// 	i := 0
// 	for i < 300000 {
// 		z[i] = x[i]/2 + y[i]/2
// 		i++
// 	}
// 	return z
// }
