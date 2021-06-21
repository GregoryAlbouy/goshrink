package imaging

import (
	"bytes"
	"image"
	"io"

	"github.com/disintegration/imaging"
)

// Format is an image file format.
type Format imaging.Format

const (
	JPEGFormat = Format(imaging.JPEG)
	PNGFormat  = Format(imaging.PNG)
)

// Ext is an image file extension, including a leading dot.
type Ext string

const (
	JPEGExt Ext = ".jpeg"
	PNGExt  Ext = ".png"
)

// FormatExt maps a file extension to a Format.
var FormatExt = map[Format]Ext{
	JPEGFormat: JPEGExt,
	PNGFormat:  PNGExt,
}

// DecodeRaw tries to decode and return a series of bytes as an image.
func DecodeRaw(rawFile []byte) (image.Image, error) {
	r := bytes.NewReader(rawFile)
	return imaging.Decode(r)
}

// Encode encodes an *image.NRGBA as bytes. Returns any error occurring
// during the process.
func Encode(w io.Writer, img *image.NRGBA, fmt Format) error {
	return imaging.Encode(w, img, imaging.Format(fmt))
}

// Reader encodes an *image.NRGBA and returns it as a new io.Reader.
func Reader(img *image.NRGBA, fmt Format) (io.Reader, error) {
	buf := bytes.Buffer{}

	if err := Encode(&buf, img, fmt); err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}

// Rescale sets an image width to the given value, preserving its aspect ratio.
// Returns the modified image.
func Rescale(img image.Image, width int) *image.NRGBA {
	return imaging.Resize(img, width, 0, imaging.NearestNeighbor)
}
