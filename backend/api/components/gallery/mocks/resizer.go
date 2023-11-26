package gallery_mocks

import "attiladudas/backend/components/gallery"

type MockResizer struct {
	ResizeImage_ func(newSize *gallery.Size, directory, filename string) ([]byte, error)
}

func (m *MockResizer) ResizeImage(newSize *gallery.Size, directory, filename string) ([]byte, error) {
	return m.ResizeImage_(newSize, directory, filename)
}
