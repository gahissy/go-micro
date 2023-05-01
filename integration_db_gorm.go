package micro

import (
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
)

type GormDBAdapter struct {
	DB
	db *gorm.DB
}

type QueryOpts struct {
	Query string
}

func (r *GormDBAdapter) First(model interface{}, conds ...interface{}) error {
	res := r.db.First(model, conds...)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *GormDBAdapter) FindBy(model interface{}, conds ...interface{}) error {
	res := r.db.Find(model, conds...)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *GormDBAdapter) FindAll(model interface{}, order interface{}, conds ...interface{}) error {
	res := r.db.Order(order).Find(model, conds...)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *GormDBAdapter) FindById(model interface{}, id string) error {
	res := r.db.First(model, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *GormDBAdapter) DeleteById(model interface{}, id string) error {
	res := r.db.Delete(model, "id = ?", id)
	return res.Error
}

func (r *GormDBAdapter) DeleteBy(model interface{}, conds ...interface{}) error {
	res := r.db.Delete(model, conds...)
	return res.Error
}

func (r *GormDBAdapter) Create(model interface{}) error {
	res := r.db.Create(model)
	return res.Error
}

func (r *GormDBAdapter) CountBy(model interface{}, query interface{}, conds ...interface{}) (int64, error) {
	var value int64
	res := r.db.Model(model).Where(query, conds...).Count(&value)
	return value, res.Error
}

func (r *GormDBAdapter) Save(model interface{}) error {
	res := r.db.Save(model)
	return res.Error
}

func (r *GormDBAdapter) Updates(model interface{}, values map[string]interface{}) error {
	res := r.db.Model(model).Updates(values)
	return res.Error
}

func (r *GormDBAdapter) Query(dest interface{}, query string, values ...interface{}) error {
	res := r.db.Raw(query, values...).Scan(dest)
	return res.Error
}

func (r *GormDBAdapter) UpdateRaw(query string, values ...interface{}) (int64, error) {
	var count int64
	res := r.db.Raw(query, values...).Scan(&count)
	return count, res.Error
}
func (r *GormDBAdapter) UpdateColumn(model interface{}, column string, value interface{}) error {
	res := r.db.Model(model).Update(column, value)
	return res.Error
}

func useGorm(config DatabaseConfig) DB {

	databaseUrl := os.Getenv("DATABASE_URL")
	var dialector gorm.Dialector
	if strings.HasPrefix(databaseUrl, "postgres") {
		dialector = postgres.Open(databaseUrl)
		// }else if strings.HasSuffix(databaseUrl, ".db") {
		//	dialector = sqlite.Open(databaseUrl)
	} else {
		panic("unknown database type: " + databaseUrl)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if config.MigrationsFs != nil {
		applyMigrations(config, db, dialector)
	}

	return &GormDBAdapter{db: db}
}

func applyMigrations(config DatabaseConfig, db *gorm.DB, dialector gorm.Dialector) {
	goose.SetBaseFS(config.MigrationsFs)
	if err := goose.SetDialect(dialector.Name()); err != nil {
		panic(err)
	}
	unwrapped, _ := db.DB()
	if err := goose.Up(unwrapped, config.MigrationsLocation); err != nil {
		panic(err)
	}
}
