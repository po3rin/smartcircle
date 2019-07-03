package smartcircle_test

import (
	"flag"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/po3rin/smartcircle"
)

var genGoldenFiles = flag.Bool("gen_golden_files", false, "whether to generate the TestXxx golden files.")

func convertRGBA(raw image.Image) *image.RGBA {
	want, ok := raw.(*image.RGBA)
	if !ok {
		b := raw.Bounds()
		want = image.NewRGBA(b)
		draw.Draw(want, b, raw, b.Min, draw.Src)
	}
	return want
}

func TestCropCircle(t *testing.T) {
	tests := []struct {
		path           string
		goldenFilename string
	}{
		{
			path:           "testdata/gopher.png",
			goldenFilename: "testdata/golden/gopher.png",
		},
		{
			path:           "testdata/person.png",
			goldenFilename: "testdata/golden/person.png",
		},
		{
			path:           "testdata/pic.jpg",
			goldenFilename: "testdata/golden/pic.jpg",
		},
		{
			path:           "testdata/whale.jpeg",
			goldenFilename: "testdata/golden/whale.jpeg",
		},
		{
			path:           "testdata/women.jpeg",
			goldenFilename: "testdata/golden/women.jpeg",
		},
	}

	for _, tt := range tests {
		img, err := os.Open(tt.path)
		if err != nil {
			log.Fatal(err)
		}
		defer img.Close()
		src, _, err := image.Decode(img)
		if err != nil {
			log.Fatal(err)
		}

		cropper, err := smartcircle.NewCropper(smartcircle.Params{Src: src})
		if err != nil {
			t.Fatalf("not expected error: %v", err.Error())
		}

		got, err := cropper.CropCircle()
		if err != nil {
			t.Fatalf("not expected error: %v", err.Error())
		}

		if *genGoldenFiles {
			goldenFile, err := os.Create(tt.goldenFilename)
			if err != nil {
				t.Errorf("failed to create file\nerr: %v", err)
			}
			defer goldenFile.Close()
			err = png.Encode(goldenFile, got)
			if err != nil {
				t.Errorf("failed to encode file\nerr: %v", err)
			}
			continue
		}

		// want
		f, err := os.Open(tt.goldenFilename)
		if err != nil {
			t.Fatalf("failed to open file\nerr: %v", err)
		}
		defer f.Close()
		want, _, err := image.Decode(f)
		if err != nil {
			t.Fatalf("failed to decode file\nerr: %v", err)
		}

		// compare RGBA.
		if !reflect.DeepEqual(convertRGBA(got), convertRGBA(want)) {
			t.Errorf("actual image differs from golden image")
			return
		}
	}
}
