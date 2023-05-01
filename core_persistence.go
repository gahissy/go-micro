package micro

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

	Updates(model interface{}, values map[string]interface{}) error

	Query(dest interface{}, query string, values ...interface{}) error

	UpdateRaw(query string, values ...interface{}) (int64, error)

	UpdateColumn(model interface{}, column string, value interface{}) error
}
