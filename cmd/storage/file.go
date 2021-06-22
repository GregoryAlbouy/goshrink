package main

import (
	"io"
	"os"
	"strings"
)

// saveUniqueUserImage saves an image file, making sure there is always
// only one for that user, overwritting the previous one.
// It uses a sub-directory named after the user ID that is systematically
// replaced.
//
// It returns the URL to access this image via GET requests.
func saveUniqueUserImage(userID, filename string, file io.Reader) (string, error) {
	dirpath := joinPath(env["STORAGE_FILE_PATH"], userID)
	return saveUniqueFileInDir(dirpath, filename, file)
}

// saveUniqueFileInDir saves a file in the target directory.
// It erases any other file in this directory before saving
// to make sure it is unique.
// It returns the path to the newly created file and the earliest
// error encountered in the process.
func saveUniqueFileInDir(dirpath, filename string, file io.Reader) (string, error) {
	// Remove directory used for storing the user image if it exists
	if err := os.RemoveAll(dirpath); err != nil {
		return "", err
	}

	// Create directory used for storing the user image
	if err := os.Mkdir(dirpath, os.ModeDir); err != nil {
		return "", err
	}

	// Create a destination on disk
	filepath := joinPath(dirpath, filename)
	dst, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy all bytes from the file to the destination on disk
	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	return filepath, nil
}

// joinPath returns a path string. The given parts are joined by a slash.
func joinPath(parts ...string) string {
	return strings.Join(parts, "/")
}
