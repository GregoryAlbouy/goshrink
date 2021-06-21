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

type multipartFormData struct {
	userID   string
	filename string
	file     io.Reader
}

type messageHandler struct {
	userService internal.UserService
}

// handle runs the full worker job process.
//
//
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

	avatarURL, err := h.postFileToStaticServer(data)
	if err != nil {
		return err
	}

	return h.updateAvatarURLInDatabase(userID, avatarURL)
}

// postFileToStaticServer makes a POST request to the static server.
// It uses a Multipart/FormData to send the userID and the image file.
func (h messageHandler) postFileToStaticServer(d multipartFormData) (string, error) {
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
	url := env["STATIC_SERVER_URL"] + "/static/avatar"
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

func (h messageHandler) updateAvatarURLInDatabase(userID, avatarURL string) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	return h.userService.SetAvatarURL(id, avatarURL)
}
