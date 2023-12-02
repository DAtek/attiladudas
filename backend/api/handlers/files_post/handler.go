package files_post

import (
	"api/components/auth"
	"api/components/gallery"
	"api/handlers"
	"api/handlers/shared"
	"bytes"
	"fibertools"
	"io"

	"github.com/DAtek/golidator"
	"github.com/DAtek/gotils"
	"github.com/gofiber/fiber/v2"
)

func PluginPostFiles(authCtx auth.IAuthorization, galleryStore gallery.IGalleryStore, fileStore gallery.IFileStore) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return postFiles(ctx, galleryStore, fileStore)
	}

	return func(app *fiber.App) {
		app.Post("/api/gallery/:id/files/", handlers.RequireUsername(authCtx), handler)
	}
}

func postFiles(ctx *fiber.Ctx, galleryStore gallery.IGalleryStore, fileStore gallery.IFileStore) error {
	pathParam, err := fibertools.BindAndValidateObj[shared.GalleryIdInPath](ctx.ParamsParser, galleryStore)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	form, formError := ctx.MultipartForm()
	if formError != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	files := form.File["files"]
	validableInput := &fileInput{
		Files: []*file{},
	}

	for _, file_ := range files {
		validableInput.Files = append(validableInput.Files, &file{Filename: file_.Filename})
	}

	if err := golidator.Validate(validableInput, pathParam.Gallery.Files); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fibertools.JsonErrorFromValidationError(err))
	}
	input := []*gallery.FileInput{}
	for _, file := range files {
		buf := &bytes.Buffer{}
		reader := gotils.ResultOrPanic(file.Open())
		io.Copy(buf, reader)
		input = append(input, &gallery.FileInput{Content: buf.Bytes(), Filename: file.Filename})
	}

	if err := fileStore.AddFiles(pathParam.Gallery, input); err != nil {
		panic(err)
	}

	ctx.Status(fiber.StatusCreated)
	return nil
}
