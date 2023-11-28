package db

import (
	"db/models"
	"math/big"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Populator struct {
	transaction *gorm.DB
}

func NewPopulator(transaction *gorm.DB) *Populator {
	return &Populator{
		transaction: transaction,
	}
}

func (p *Populator) User(data ...map[string]any) *models.User {
	faker := newFaker()
	obj := &models.User{
		Username:     faker.Internet().User(),
		PasswordHash: faker.Internet().Password(),
	}

	return createWithCustomValues(p.transaction, obj, data...)
}

func (p *Populator) Gallery(data ...map[string]any) *models.Gallery {
	faker := newFaker()
	t := time.Now().UTC()
	d := datatypes.Date(t)

	obj := &models.Gallery{
		Title:       faker.Lorem().Word(),
		Slug:        faker.Lorem().Word(),
		Description: faker.Lorem().Sentence(5),
		Date:        &d,
		Directory:   faker.Bothify("??????????????"),
	}

	return createWithCustomValues(p.transaction, obj, data...)
}

func createRelationIfNeeded[T any](
	createRecord func(data ...map[string]any) *T,
	data ...map[string]any,
) (obj *T) {
	ok := false
	ref := reflect.ValueOf(createRecord).Pointer()
	fullName := runtime.FuncForPC(ref).Name()
	parts := strings.Split(fullName, ".")
	key := strings.Replace(parts[len(parts)-1], "-fm", "", -1)

	switch len(data) {
	case 0:
		obj, ok = nil, false
	case 1:
		obj, ok = data[0][key].(*T)
	default:
		panic("Wrong value for data")
	}

	if obj != nil && !ok {
		panic("Wrong type for " + key)
	}

	if obj == nil {
		obj = createRecord()
	}

	return obj
}

func createWithCustomValues[T any](transaction *gorm.DB, obj T, data ...map[string]any) T {
	var d map[string]any
	switch len(data) {
	case 0:
		d = nil
	case 1:
		d = data[0]
	default:
		panic("Wrong value for data")
	}

	for key, value := range d {
		obj := reflect.ValueOf(obj).Elem()
		field := obj.FieldByName(key)
		newVal := reflect.ValueOf(value)
		field.Set(newVal)
	}
	ResultOrPanic(transaction.Create(obj))
	return obj
}

func newFaker() faker.Faker {
	uid := uuid.New()
	uidInt := &big.Int{}
	uidInt.SetBytes(uid[:])
	seed := rand.NewSource(uidInt.Int64())
	return faker.NewWithSeed(seed)
}
