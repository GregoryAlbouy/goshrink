package main

import (
	"fmt"
	"io"

	"github.com/GregoryAlbouy/shrinker/pkg/imaging"
)

const imageFormat = imaging.PNGFormat // TODO: make it actionnable

// Rezise width used by the worker.
var resizeWidth int

// processImage runs the full image processing logic.
//
// It decodes the bytes input into an image, performs the rescale and
// returns a io.Reader to access the output file.
// Also returns the file extension.
func processImage(raw []byte) (io.Reader, imaging.Ext, error) {
	img, err := imaging.DecodeRaw(raw)
	if err != nil {
		return nil, "", fmt.Errorf("image decoding error: %w", err)
	}

	resized := imaging.Rescale(img, resizeWidth)

	// Encode and retrieve a file reader
	reader, err := imaging.Reader(resized, imageFormat)
	if err != nil {
		return nil, "", fmt.Errorf("image encoding error: %w", err)
	}

	// Manually get file extension
	ext := imaging.FormatExt[imageFormat]

	return reader, ext, nil
}
