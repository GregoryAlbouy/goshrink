package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// multipartFormData represents the necessary data for creating
// a new request with content type "multipart/form-data".
type multipartFormData struct {
	userID   string
	filename string
	file     io.Reader
}

// messageHandler represents a service to handle
// incoming messages from the queue.
type messageHandler struct {
	userService internal.UserService
}

// handle runs the full worker job process.
//
// It parses messages from the queue and invokes the resizing function on
// the file held inside. It then posts to the storage the modified image
// will all the data needed using "multipart/form-data" content type.
// Upon receiving confirmation of the creation of the file storage side,
// it writes the newly created file URL to the database.
//
// Any error during the process will is caught and returned to the caller.
func (h messageHandler) handle(d amqp.Delivery) error {
	log.Printf("Start handling message %d...", d.DeliveryTag)

	// Retrieve message values
	userID := d.MessageId
	rawFile := d.Body

	// UserID is required
	if userID == "" {
		return errors.New("no userID provided")
	}

	// Read and rescale image
	imageReader, ext, err := processImage(rawFile)
	if err != nil {
		return err
	}

	data := multipartFormData{
		userID:   d.MessageId,
		filename: uuid.NewString() + string(ext),
		file:     imageReader,
	}

	avatarURL, err := h.postImageToStorage(data)
	if err != nil {
		return err
	}

	return h.writeURLToDatabase(userID, avatarURL)
}

// postImageToStorage makes a POST request to the storage.
// It sends a "multipart/form-data" request with the userID and the image file.
func (h messageHandler) postImageToStorage(d multipartFormData) (string, error) {
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	defer writer.Close()

	// Create and write userID part
	idPart, err := writer.CreateFormField("userId")
	if err != nil {
		return "", err
	}
	idPart.Write([]byte(d.userID))

	// Create and write file part
	filePart, err := writer.CreateFormFile("image", d.filename)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(filePart, d.file); err != nil {
		return "", err
	}

	writer.Close()

	// Build HTTP request
	url := env["STATIC_SERVER_URL"] + "/storage/avatar"
	request, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", "Bearer "+env["STATIC_SERVER_KEY"])

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("static server sent: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Handle static server response
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(content), nil
}

func (h messageHandler) writeURLToDatabase(userID, avatarURL string) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	return h.userService.SetAvatarURL(id, avatarURL)
}
