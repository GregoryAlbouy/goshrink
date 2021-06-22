package message

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// FileMessage represents a message carrying a file.
// The file is represented as a stream of bytes.
// An ID can be set to help identification.
type FileMessage struct {
	ID   int
	File []byte
}

// Bytes encodes a FileMessage struct into a slice of bytes.
func (msg FileMessage) Bytes() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(&msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// BytesJSON encodes a FileMessage struct into a slice of bytes
// that can be decoded into JSON. It is about 4 times slower
// than Bytes(), but the output is supported by any language.
func (msg FileMessage) BytesJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)

	if err := enc.Encode(msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Read reads a slice of bytes and returns a FileMessage
// and a non-nil error in case it went wrong.
func Read(in []byte) (FileMessage, error) {
	buf := bytes.NewBuffer(in)
	dec := gob.NewDecoder(buf)
	msg := FileMessage{}

	if err := dec.Decode(&msg); err != nil {
		return FileMessage{}, err
	}
	return msg, nil
}
