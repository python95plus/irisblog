package library

import (
	"fmt"
	"image"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

func ThumbnailCrop(minWidth, minHeight uint, img image.Image) image.Image {
	originalBounds := img.Bounds()
	originalWidth := uint(originalBounds.Dx())
	originalHeight := uint(originalBounds.Dy())
	newWidth, newHeight := originalWidth, originalHeight

	if minWidth >= originalWidth && minHeight >= originalHeight {
		return img
	}

	if minWidth > originalWidth {
		minWidth = originalWidth
	}

	if minHeight > originalHeight {
		minHeight = originalHeight
	}

	if originalWidth > minWidth {
		newHeight = uint(originalHeight * minWidth / originalWidth)
		if newHeight < 1 {
			newHeight = 1
		}
	}

	if newHeight < minHeight {
		newWidth = uint(newWidth * minHeight / newHeight)
		if newWidth < 1 {
			newWidth = 1
		}
	}

	if originalWidth > originalHeight {
		newWidth = minWidth
		newHeight = 0
	} else {
		newWidth = 0
		newHeight = minHeight
	}

	thumbImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return thumbImg
}

func Resize(width, height uint, img image.Image, interp resize.InterpolationFunction) image.Image {
	return resize.Resize(width, height, img, interp)
}

func Thumbnail(width, height uint, img image.Image, interp resize.InterpolationFunction) image.Image {
	return resize.Thumbnail(width, height, img, interp)
}

func CropImg(srcImg image.Image, dstWidth, dstHeight int) image.Image {
	//origBounds := srcImg.Bounds()
	//origWidth := origBounds.Dx()
	//origHeight := origBounds.Dy()

	dstImg, err := cutter.Crop(srcImg, cutter.Config{
		Height: dstHeight,       // height in pixel or Y ratio(see Ratio Option below)
		Width:  dstWidth,        // width in pixel or X ratio
		Mode:   cutter.Centered, // Accepted Mode: TopLeft, Centered
		//Anchor: image.Point{
		//	origWidth / 12,
		//	origHeight / 8}, // Position of the top left point
		Options: 0, // Accepted Option: Ratio
	})
	fmt.Println()
	if err != nil {
		fmt.Println("Cannot crop image:" + err.Error())
		return srcImg
	}
	return dstImg
}
