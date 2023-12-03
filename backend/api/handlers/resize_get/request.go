package resize_get

import (
	"api"
	"api/components/gallery"
	"strconv"
	"strings"

	"github.com/DAtek/golidator"
	"github.com/DAtek/gotils"
)

type pathData struct {
	Directory string `params:"directory"`
	Size      string `params:"size"`
	Filename  string `params:"filename"`

	size *gallery.Size
}

func (obj *pathData) GetValidators(params ...any) golidator.ValidatorCollection {
	return golidator.ValidatorCollection{
		{Field: "size", Function: func() *golidator.ValueError {
			if !validSizes.contains(obj.Size) {
				return api.ErrorInvalid
			}

			parts := strings.Split(obj.Size, "x")
			values := []uint{}

			for _, part := range parts {
				v := gotils.ResultOrPanic(strconv.Atoi(part))
				values = append(values, uint(v))
			}

			obj.size = &gallery.Size{Width: values[0], Height: values[1]}
			return nil
		}},
	}
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
