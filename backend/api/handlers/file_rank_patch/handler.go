package file_rank_patch

import (
	"api/components/auth"
	"api/components/gallery"
	"api/handlers"
	"fibertools"
	"net/http"

	"github.com/DAtek/golidator"
	"github.com/gofiber/fiber/v2"
)

const path = "/api/file/:id/rank/:rank/"

func PluginPatchFileRank(store gallery.IFileStore, authCtx auth.IAuthorization) fibertools.Plugin {
	handler := func(ctx *fiber.Ctx) error {
		return patchFileRank(ctx, store)
	}
	return func(app *fiber.App) {
		app.Patch(path, handlers.RequireUsername(authCtx), handler)
	}
}

func patchFileRank(ctx *fiber.Ctx, store gallery.IFileStore) error {
	input, err := fibertools.BindAndValidateParams(&patchFileRankPathParams{}, ctx, store)
	if err := golidator.Validate(input, store); err != nil {
		ctx.Status(http.StatusNotFound)
		return nil
	}
	dbErr := store.UpdateFileRank(&gallery.UpdateFileRankInput{FileId: *input.FileId, Rank: *input.Rank})

	if err != nil {
		panic(dbErr)
	}

	return nil
}
