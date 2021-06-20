package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// TODO: rename
type fileMessage struct {
	userID   string
	filename string
	file     io.Reader
}

type messageHandler struct {
	userService internal.UserService
}

func (h messageHandler) handle(d amqp.Delivery) error {
	log.Println("Got avatar from " + d.MessageId)

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

	// Build the data to be sent to database and static server
	filename := uuid.NewString() + string(ext)
	fm := fileMessage{
		userID:   d.MessageId,
		filename: filename,
		file:     imageReader,
	}

	// Post image file to static server
	avatarURL, err := h.postFileToStaticServer(fm)
	if err != nil {
		return err
	}

	// Insert avatarURL into the database
	return h.updateAvatarURLInDatabase(userID, avatarURL)
}

// postFileToStaticServer makes a POST request to the static server.
// It uses a Multipart/FormData to send the userID and the image file.
func (h messageHandler) postFileToStaticServer(fm fileMessage) (string, error) {
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	defer writer.Close()

	// Create and write userID part
	idPart, err := writer.CreateFormField("userId")
	if err != nil {
		return "", err
	}
	idPart.Write([]byte(fm.userID))

	// Create and write file part
	filePart, err := writer.CreateFormFile("image", fm.filename)
	if err != nil {
		return "", err
	}
	if n, err := io.Copy(filePart, fm.file); err != nil {
		return "", errors.New(fmt.Sprintf("%s\n%d bytes written\n", err, n))
	}

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

	if err := validateStatusCode(resp.StatusCode); err != nil {
		return "", err
	}

	// Handle static server response
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	avatarURL := string(content)
	log.Println(content)
	return avatarURL, nil
}

func (h messageHandler) updateAvatarURLInDatabase(userID, avatarURL string) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	return h.userService.SetAvatarURL(id, avatarURL)
}

func validateStatusCode(code int) error {
	if code < 200 || code > 299 {
		return errors.New("bad")
	}
	return nil
}
