package main

import (
	"fmt"
	"io"

	"github.com/GregoryAlbouy/shrinker/pkg/imaging"
)

const (
	imageFormat = imaging.PNGFormat // TODO: make it actionnable
)

var resizeWidth int

// processImage executes the whole image processing logic.
// It takes a slice of bytes in input, decodes it to an image,
// performs the rescale, and finally return a reader to access
// the resulting file along with the image file extension.
// It returns a non-nil error if anything went wrong.
//
// FIXME: imaging.DecodeRaw produces an error
func processImage(raw []byte) (io.Reader, imaging.Ext, error) {
	// // temp debug
	// log.Println("Processing image")

	// Decode raw file
	img, err := imaging.DecodeRaw(raw)
	if err != nil {
		return nil, "", fmt.Errorf("image decoding error: %w", err)
	}

	// Rescale
	resized := imaging.Rescale(img, resizeWidth)

	// Encode and retrieve a file reader
	reader, err := imaging.Reader(resized, imageFormat)
	if err != nil {
		return nil, "", fmt.Errorf("image encoding error: %w", err)
	}

	// Get file extension
	ext := imaging.FormatExt[imageFormat]

	return reader, ext, nil
}
