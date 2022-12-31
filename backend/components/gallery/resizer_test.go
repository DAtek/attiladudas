package gallery

import (
	"attiladudas/backend/components"
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResizeImage(t *testing.T) {
	resizer := getResizer()
	fruitGalleryDir := "fruit_gallery"
	filename := "apple.jpg"

	sizes := []Size{
		{200, 115},
		{115, 200},
		{474, 358},
	}
	for _, size := range sizes {
		testName := fmt.Sprintf("Resizes image properly %dx%d", size.Width, size.Height)
		t.Run(testName, func(t *testing.T) {
			newSizeString := fmt.Sprintf("%dx%d", size.Width, size.Height)
			defer func() {
				os.RemoveAll(
					filepath.Join(components.EnvTestFilesDir.Load(),
						fmt.Sprintf("%s/%dx%d", fruitGalleryDir, size.Width, size.Height)),
				)
			}()

			result, err := resizer.ResizeImage(&size, fruitGalleryDir, filename)
			assert.Nil(t, err)

			imgPath := filepath.Join(
				components.EnvTestFilesDir.Load(),
				fmt.Sprintf("%s/%s/%s", fruitGalleryDir, newSizeString, filename),
			)
			_, fileErr := os.Stat(imgPath)
			assert.Nil(t, fileErr)

			imgContent, _ := os.ReadFile(imgPath)
			assert.Equal(t, imgContent, result)
			img, _, _ := image.Decode(bytes.NewBuffer(imgContent))
			assert.Equal(t, image.Point{int(size.Width), int(size.Height)}, img.Bounds().Max)
		})
	}
}
