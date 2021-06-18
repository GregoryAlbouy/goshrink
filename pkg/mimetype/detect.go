package mimetype

import (
	"io"
	"net/http"
)

// Detect returns the MIME type of file. The file must be provided as a stream.
//
// Using Detect moves the pointer reading bytes of the io.Reader. It the file
// needs to be read again, make sure to move the pointer back at the start with
// `file.Seek(0, io.SeekStart)`.
func Detect(file io.Reader) (string, error) {
	// Only the first 512 bytes are needed to get the content type.
	buf := make([]byte, 512)

	if _, err := file.Read(buf); err != nil {
		return "", err
	}
	return http.DetectContentType(buf), nil
}

// IsImage checks if the given file is an image type.
//
// Using IsImage moves the pointer reading bytes of the io.Reader. It the file
// needs to be read again, make sure to move the pointer back at the start with
// `file.Seek(0, io.SeekStart)`.
func IsImage(file io.Reader) bool {
	// An error while reading the file is interpreted as "not an image".
	kind, _ := Detect(file)
	return kind == JPEG || kind == PNG
}
