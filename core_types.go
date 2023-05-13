package micro

import (
	"github.com/gahissy/go-micro/h"
	"github.com/gahissy/go-micro/ports"
	"github.com/go-co-op/gocron"
	"io/fs"
)

type Any = interface{}

type Err = error

type App struct {
	scheduler *gocron.Scheduler
	router    Router
	Env       *Env
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
	Profile             string
	DB                  DB
	App                 *AppInfo
	NotificationManager ports.NotificationManager
}

type Ctx struct {
	Env  *Env
	Auth *Auth
}

type Worker struct {
	Every  string
	Handle func(ctx *Ctx) error
}

type DbConfig struct {
	MigrationsFs       fs.FS
	MigrationsLocation string
	Seeders            []func(ctx *Ctx) error
}

func (e *Env) SendNotification(notification ports.Notification) error {
	return e.NotificationManager.Send(notification)
}

func (e *Env) IsProduction() bool {
	return h.IsProduction(e.Profile)
}
