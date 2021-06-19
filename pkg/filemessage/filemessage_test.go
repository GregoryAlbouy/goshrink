package message

import (
	"bytes"
	"log"
	"os"
	"testing"
)

// Tests

func TestConversion(t *testing.T) {
	in := newFileMessage()

	b, err := in.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	out, err := Read(b)
	if err != nil {
		t.Fatal(err)
	}

	if in.ID != out.ID {
		t.Error("ID are not equal")
	}

	if bytes.Compare(in.File, out.File) != 0 {
		t.Error("Files are not equal")
	}
}

// Benchmarks

func BenchmarkCompare(b *testing.B) {
	b.Run("Struct encoding", BenchmarkStruct)
	b.Run("JSON encoding", BenchmarkJSON)
}

func BenchmarkStruct(b *testing.B) {
	msg := newFileMessage()

	for i := 0; i < b.N; i++ {
		msg.Bytes()
	}
}

func BenchmarkJSON(b *testing.B) {
	msg := newFileMessage()

	for i := 0; i < b.N; i++ {
		msg.JSONBytes()
	}
}

// Helpers

func newFileMessage() FileMessage {
	img, err := os.ReadFile("image.png")
	if err != nil {
		log.Fatal(err)
	}
	return FileMessage{12, img}
}
