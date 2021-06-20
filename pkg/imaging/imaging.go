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

// FormatExt associates a file extension to a Format
var FormatExt = map[Format]Ext{
	JPEGFormat: JPEGExt,
	PNGFormat:  PNGExt,
}

// DecodeRaw decodes a series of bytes and returns an image.Image
// or a non-nil error.
func DecodeRaw(rawFile []byte) (image.Image, error) {
	r := bytes.NewReader(rawFile)
	return imaging.Decode(r)
}

// Encode encodes an *image.NRGBA to a []byte and returns an error.
func Encode(w io.Writer, img *image.NRGBA, fmt Format) error {
	return imaging.Encode(w, img, imaging.Format(fmt))
}

// Reader encodes an *image.NRGBA and returns a new file reader
// or a non-nil error.
func Reader(img *image.NRGBA, fmt Format) (io.Reader, error) {
	buf := bytes.Buffer{}

	if err := Encode(&buf, img, fmt); err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}

// Rescale sets an image width to the given value, preserving the aspect ratio.
func Rescale(img image.Image, width int) *image.NRGBA {
	return imaging.Resize(img, width, 0, imaging.NearestNeighbor)
}
