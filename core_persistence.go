package micro

import (
	"errors"
	"gorm.io/gorm"
)

type Entity interface {
}

type DB interface {
	First(model interface{}, conds ...interface{}) error

	FindAll(model interface{}, order interface{}, conds ...interface{}) error

	FindBy(model interface{}, conds ...interface{}) error

	FindById(model interface{}, id string) error

	DeleteById(model interface{}, id string) error

	DeleteBy(model interface{}, conds ...interface{}) error

	Create(model interface{}) error

	CountBy(model interface{}, query interface{}, conds ...interface{}) (int64, error)

	Save(model interface{}) error

	Patch(model interface{}, values map[string]interface{}) error

	Query(dest interface{}, query string, values ...interface{}) error

	UpdateRaw(query string, values ...interface{}) (int64, error)

	UpdateColumn(model interface{}, column string, value interface{}) error

	//Repo(entity interface{}) *Repo
}

type Repo[T interface{}] struct {
	DB     DB
	Entity T
}

func NewRepo[T Entity](ctx *Ctx, entity T) *Repo[T] {
	return &Repo[T]{
		DB:     ctx.Env.DB,
		Entity: entity,
	}
}

func (r *Repo[T]) FindById(id string) *T {
	var entity T
	err := r.DB.FindById(&entity, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return &entity
}

func (r *Repo[T]) FindBy(criteria string, params ...interface{}) (*T, error) {
	var entity T
	err := r.DB.FindBy(&entity, criteria, params)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *Repo[T]) CountBy(criteria string, params ...interface{}) (int64, error) {
	var entity T
	return r.DB.CountBy(entity, criteria, params)
}

func (r *Repo[T]) FindAll(orderBy string) ([]*T, error) {
	var entities []*T
	err := r.DB.FindAll(&entities, orderBy)
	return entities, err
}

func (r *Repo[T]) DeleteById(id string) error {
	return r.DB.DeleteById(r.Entity, id)
}

func (r *Repo[T]) Create(model *T) error {
	return r.DB.Create(model)
}

func (r *Repo[T]) Update(model *T) error {
	return r.DB.Save(model)
}

func (r *Repo[T]) Patch(model *T, changes map[string]interface{}) error {
	return r.DB.Patch(model, changes)
}
