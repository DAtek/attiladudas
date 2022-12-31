package gallery

import (
	"attiladudas/backend/components"
	"attiladudas/backend/helpers"
	"attiladudas/backend/models"

	"github.com/jaswdr/faker"
	"gorm.io/gorm"
)

// useful singletons used by tests

var db *gorm.DB = nil

func getDb() *gorm.DB {
	if db == nil {
		db = createOrPanic(components.NewDbFromEnv)
	}

	return db.Begin()
}

var emptyDb *gorm.DB = nil

func getEmptyDb() *gorm.DB {
	if emptyDb == nil {
		emptyDb = createOrPanic(newEmptyDb)
	}

	return emptyDb
}

func newEmptyDb() (*gorm.DB, error) {
	return components.NewDb(
		components.EnvDbHost.Load(),
		components.EnvDbUser.Load(),
		components.EnvDbPassword.Load(),
		components.EnvDbName.Load(),
		components.EnvEmptyPostgresPort.Load(),
	)
}

var res IResizer = nil

func getResizer() IResizer {
	if res == nil {
		res = NewResizer(components.EnvTestFilesDir.Load())
	}

	return res
}

func createOrPanic[T interface{}](f func() (T, error)) T {
	result, err := f()
	if err != nil {
		panic(err)
	}
	return result
}

func panicErrorResult(result *gorm.DB) {
	if result.Error != nil {
		panic(result.Error)
	}
}

var fake = faker.New()

func NewDummyGallery() *models.Gallery {
	return &models.Gallery{
		Title:     fake.Company().Name(),
		Slug:      fake.Bothify("?????####"),
		Date:      helpers.DateFromISO8601Panic("2022-01-01"),
		Directory: fake.UUID().V4(),
		Active:    fake.Bool(),
	}
}
