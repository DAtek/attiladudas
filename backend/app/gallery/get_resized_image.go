package app_gallery

import (
	"attiladudas/backend/components/gallery"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DAtek/golidator"
	"github.com/gin-gonic/gin"
)

func GetResizedImageHandler(resizer gallery.IResizer, mediaDirName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resizeImage(ctx, resizer, mediaDirName)
	}
}

func resizeImage(ctx *gin.Context, resizer gallery.IResizer, mediaDirName string) {
	params := &uriParams{}
	ctx.ShouldBindUri(params)

	if err := golidator.Validate(params); err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	if _, err := resizer.ResizeImage(params.size, params.Directory, params.Filename); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Redirect(http.StatusFound, fmt.Sprintf(
		"/%s/%s/%s/%s",
		mediaDirName,
		params.Directory,
		params.Size,
		params.Filename,
	))
}

type uriParams struct {
	Directory string `uri:"directory"`
	Size      string `uri:"size"`
	Filename  string `uri:"filename"`
	size      *gallery.Size
}

type sizes []string

func (s sizes) contains(item string) bool {
	for _, size := range s {
		if size == item {
			return true
		}
	}

	return false
}

var validSizes = sizes{
	"400x225",
	"600x337",
	"1366x768",
}

func (obj *uriParams) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	return golidator.ValidatorCollection{
		{Field: "size", Function: func() *golidator.ValueError {
			if !validSizes.contains(obj.Size) {
				return &golidator.ValueError{ErrorType: "INVALID_SIZE"}
			}

			parts := strings.Split(obj.Size, "x")
			values := []uint{}
			for _, part := range parts {
				v, _ := strconv.Atoi(part)
				values = append(values, uint(v))
			}
			obj.size = &gallery.Size{Width: values[0], Height: values[1]}

			return nil
		}},
	}
}
