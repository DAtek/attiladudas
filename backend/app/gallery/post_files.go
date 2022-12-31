package app_gallery

import (
	app_helpers "attiladudas/backend/app/helpers"
	"attiladudas/backend/components/gallery"
	"attiladudas/backend/helpers"
	"attiladudas/backend/models"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/DAtek/golidator"
	"github.com/gin-gonic/gin"
)

type file struct {
	Filename string
}

type fileInput struct {
	Files []*file
}

func PostFilesHandler(fileStore gallery.IFileStore, gallerStore gallery.IGalleryStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		postFiles(ctx, fileStore, gallerStore)
	}
}

func postFiles(ctx *gin.Context, fileStore gallery.IFileStore, gallerStore gallery.IGalleryStore) {
	id := &galleryIdUri{}
	if err := ctx.ShouldBindUri(id); err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	galleryObj, err := gallerStore.GetGallery(&gallery.GetGalleryInput{Id: &id.Id})

	if app_helpers.HandleError(ctx, err) {
		return
	}

	form, _ := ctx.MultipartForm()
	files := form.File["files"]
	validableInput := &fileInput{
		Files: []*file{},
	}
	for _, file_ := range files {
		validableInput.Files = append(validableInput.Files, &file{Filename: file_.Filename})
	}

	if err := golidator.Validate(validableInput, galleryObj.Files); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.JsonErrorFromValidationError(err))
	}
	input := []*gallery.FileInput{}
	for _, file := range files {
		buf := &bytes.Buffer{}
		reader, _ := file.Open()
		io.Copy(buf, reader)
		input = append(input, &gallery.FileInput{Content: buf.Bytes(), Filename: file.Filename})
	}
	if err := fileStore.AddFiles(galleryObj, input); err != nil {
		fmt.Printf("Error happened when saving the files: %v\n", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (obj *fileInput) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	return golidator.GetValidatorsForList("files", obj.Files, params...)
}

func (obj *file) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	existingFiles, ok := params[0].([]*models.File)
	if !ok {
		panic("Existing files were not provided")
	}

	return []*golidator.Validator{
		{Field: "filename", Function: func() *golidator.ValueError {
			for _, file := range existingFiles {
				if obj.Filename == file.Filename {
					return &golidator.ValueError{ErrorType: helpers.ErrorAlreadyExists}
				}
			}
			return nil
		}},
	}
}
