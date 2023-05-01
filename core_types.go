package micro

import (
	"github.com/go-co-op/gocron"
	"io/fs"
)

type Any = interface{}

type Err = error

type App struct {
	scheduler *gocron.Scheduler
	router    Router
	env       *Env
}

type AppInfo struct {
	Name    string
	Version string
}

type Auth struct {
	Id            string
	Role          string
	Authenticated bool
}

type Service interface {
	Workers() []func(r Worker)
	Routes() []func(r Router)
}

type Env struct {
	DB  DB
	App *AppInfo
}

type Ctx struct {
	Env  *Env
	Auth *Auth
}

type Worker struct {
	Every  string
	Handle func(ctx *Env) error
}

type DatabaseConfig struct {
	MigrationsFs       fs.FS
	MigrationsLocation string
}
