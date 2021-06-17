package mimetype

import (
	"io"
	"net/http"
)

// Detect returns the MIME type of file. The file must be provided as a stream.
func Detect(file io.Reader) (string, error) {
	buf := make([]byte, 512)

	if _, err := file.Read(buf); err != nil {
		return "", err
	}
	return http.DetectContentType(buf), nil
}

// IsImage checks if the given file is an image type.
func IsImage(file io.Reader) bool {
	// An error while reading the file is interpreted as "not an image".
	kind, _ := Detect(file)
	return kind == jpeg || kind == png
}
