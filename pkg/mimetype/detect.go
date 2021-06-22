package mimetype

import (
	"io"
	"net/http"
)

const contentTypeMaxBytes = 512

// Detect returns the MIME type of file. The file must be provided as a an
// interface that implements the basic Read and Seek methods.
func Detect(file io.ReadSeeker) (string, error) {
	buf := make([]byte, contentTypeMaxBytes)

	if _, err := file.Read(buf); err != nil {
		return "", err
	}
	// Sets the pointer to the beginning of the file so it can be read again.
	defer file.Seek(0, io.SeekStart)
	return http.DetectContentType(buf), nil
}

// IsImage checks if the given file is an image type. The file must be provided as a an
// interface that implements the basic Read and Seek methods.
func IsImage(file io.ReadSeeker) bool {
	// An error while reading the file is interpreted as "not an image".
	kind, _ := Detect(file)
	return kind == JPEG || kind == PNG
}
