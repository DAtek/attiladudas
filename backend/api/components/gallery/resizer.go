package gallery

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

type ImageSize string

type IResizer interface {
	ResizeImage(newSize *Size, directory, filename string) ([]byte, error)
}

type Size struct {
	Width  uint
	Height uint
}

type resizer struct {
	baseDir string
}

func NewResizer(baseDir string) IResizer {
	return &resizer{baseDir: baseDir}
}

func (r *resizer) ResizeImage(newSize *Size, directory, filename string) ([]byte, error) {
	originalPath := filepath.Join(r.baseDir, directory, filename)
	imageContent, readErr := os.ReadFile(originalPath)

	if readErr != nil {
		return nil, readErr
	}

	subdir := fmt.Sprintf("%dx%d", newSize.Width, newSize.Height)
	directoryPath := filepath.Join(r.baseDir, directory, subdir)
	_, directoryErr := os.Stat(directoryPath)

	if os.IsNotExist(directoryErr) {
		os.Mkdir(directoryPath, 0750)
	} else if directoryErr != nil {
		return nil, directoryErr
	}

	img, _, decodingErr := image.Decode(bytes.NewBuffer(imageContent))

	if decodingErr != nil {
		return nil, decodingErr
	}

	thumbnail := newThumbnail(newSize, img)
	thumbnail.resize()
	thumbnail.crop()
	encodedBuf := &bytes.Buffer{}
	encodingErr := jpeg.Encode(encodedBuf, thumbnail.img, nil)

	if encodingErr != nil {
		return nil, encodingErr
	}

	writeErr := os.WriteFile(filepath.Join(directoryPath, filename), encodedBuf.Bytes(), 0750)
	return encodedBuf.Bytes(), writeErr
}

type thumbnail struct {
	originalSize        *Size
	originalAspectRatio float32
	newSize             *Size
	newAspectRatio      float32
	img                 image.Image
}

type subImager interface {
	SubImage(image.Rectangle) image.Image
}

func newThumbnail(newSize *Size, img image.Image) *thumbnail {
	return &thumbnail{
		originalSize:        &Size{uint(img.Bounds().Max.X), uint(img.Bounds().Max.Y)},
		originalAspectRatio: float32(img.Bounds().Max.X) / float32(img.Bounds().Max.Y),
		newSize:             newSize,
		newAspectRatio:      float32(newSize.Width) / float32(newSize.Height),
		img:                 img,
	}
}

func (t *thumbnail) resize() {
	newSize := t.getNewSizeForResizing()
	t.img = resize.Resize(newSize.Width, newSize.Height, t.img, resize.Lanczos3)
}

func (t *thumbnail) getNewSizeForResizing() *Size {
	if t.newAspectRatio > t.originalAspectRatio {
		height := float32(t.newSize.Width) / t.originalAspectRatio
		return &Size{t.newSize.Width, uint(height)}
	}

	if t.newAspectRatio < t.originalAspectRatio {
		width := t.originalAspectRatio * float32(t.newSize.Height)
		return &Size{uint(width), t.newSize.Height}
	}

	return &Size{t.newSize.Width, t.newSize.Height}
}

func (t *thumbnail) crop() error {
	if t.originalAspectRatio == t.newAspectRatio {
		return nil
	}

	if t.newAspectRatio < t.originalAspectRatio {
		return t.cropWidth()
	}

	return t.cropHeight()
}

func (t *thumbnail) cropWidth() error {
	img := t.img.(subImager)
	diff := t.img.Bounds().Max.X - int(t.newSize.Width)
	halfDiff := diff / 2
	t.img = img.SubImage(image.Rectangle{
		Min: image.Point{halfDiff, 0},
		Max: image.Point{
			int(t.newSize.Width) + halfDiff,
			int(t.newSize.Height),
		},
	})
	return nil

}

func (t *thumbnail) cropHeight() error {
	img := t.img.(subImager)
	diff := t.img.Bounds().Max.Y - int(t.newSize.Height)
	halfDiff := diff / 2
	t.img = img.SubImage(image.Rectangle{
		Min: image.Point{0, halfDiff},
		Max: image.Point{
			int(t.newSize.Width),
			int(t.newSize.Height) + halfDiff,
		},
	})
	return nil
}
