package imaging

import (
	"image"
	"log"
	"os"
	"testing"

	"github.com/disintegration/imaging"
)

const (
	imgPath   = "./fixtures/image.png"
	noimgPath = "./fixtures/not-an-image.png"
)

// Tests

func TestDecodeRaw(t *testing.T) {
	img := mustReadFile(imgPath)
	if _, err := DecodeRaw(img); err != nil {
		t.Fatalf("should decode a regular image, got err: %s", err)
	}

	noImg := mustReadFile(noimgPath)
	if _, err := DecodeRaw(noImg); err == nil { // if err EQUALS nil
		t.Fatal("should fail to decode a non-image fail")
	}
}

// Benchmarks

func BenchmarkFilters(b *testing.B) {
	b.Run("NearestNeighbor", benchmarkNearestNeighbor)
	b.Run("Lanczos", benchmarkLanczos)
	b.Run("Linear", benchmarkLinear)
	b.Run("Box", benchmarkBox)
	b.Run("CatmullRom", benchmarkCatmullRom)
	b.Run("MitchellNetravali", benchmarkMitchellNetravali)
}

func benchmarkNearestNeighbor(b *testing.B) {
	img := mustDecodeImage(imgPath)

	for i := 0; i < b.N; i++ {
		rescaleWithFilter(img, imaging.NearestNeighbor)
	}
}

func benchmarkLanczos(b *testing.B) {
	img := mustDecodeImage(imgPath)

	for i := 0; i < b.N; i++ {
		rescaleWithFilter(img, imaging.Lanczos)
	}
}

func benchmarkLinear(b *testing.B) {
	img := mustDecodeImage(imgPath)

	for i := 0; i < b.N; i++ {
		rescaleWithFilter(img, imaging.Linear)
	}
}

func benchmarkBox(b *testing.B) {
	img := mustDecodeImage(imgPath)

	for i := 0; i < b.N; i++ {
		rescaleWithFilter(img, imaging.Box)
	}
}

func benchmarkCatmullRom(b *testing.B) {
	img := mustDecodeImage(imgPath)

	for i := 0; i < b.N; i++ {
		rescaleWithFilter(img, imaging.CatmullRom)
	}
}

func benchmarkMitchellNetravali(b *testing.B) {
	img := mustDecodeImage(imgPath)

	for i := 0; i < b.N; i++ {
		rescaleWithFilter(img, imaging.MitchellNetravali)
	}
}

// Helpers

func mustReadFile(filepath string) []byte {
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func mustDecodeImage(filepath string) image.Image {
	file := mustReadFile(filepath)
	img, err := DecodeRaw(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func rescaleWithFilter(img image.Image, f imaging.ResampleFilter) {
	imaging.Resize(img, 200, 0, f)
}
