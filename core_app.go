package micro

import (
	"github.com/gahissy/go-micro/adapters/discord"
	"github.com/gahissy/go-micro/adapters/slack"
	"github.com/gahissy/go-micro/h"
	"github.com/gahissy/go-micro/ports"
	"github.com/gahissy/go-micro/sentry"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

var AppName string

var AppVersion string

func New(profile string, name string, version string) *App {
	AppName = name
	AppVersion = version

	match, err := h.FindFileInParents(".Env")
	if err == nil {
		_ = godotenv.Load(match)
	}
	if profile != "" && profile != "default" {
		filename := ".Env." + strings.ToLower(profile)
		match, err = h.FindFileInParents(filename)
		if err == nil {
			_ = godotenv.Load(match)
		}
	}

	viper.AutomaticEnv()
	_ = viper.ReadInConfig() // Find and read the config file

	env := &Env{
		Profile: profile,
		App: &AppInfo{
			Name:    name,
			Version: version,
		},
	}

	if env.Profile == "" || env.Profile == "default" {
		env.Profile = h.GetEnv("PROFILE")
	}

	env.configureErrorReporter()
	env.configureNotificationsManager()

	cron := gocron.NewScheduler(time.UTC)

	app := &App{
		scheduler: cron,
		//router:    router,
		Env: env,
	}

	return app
}

func (e *Env) configureErrorReporter() {
	dsn := h.GetEnv("SENTRY_DSN")
	if dsn != "" {
		log.Infof("SENTRY_DSN detected, configuring")
		sentry.Configure(dsn, e.App.Version, e.Profile, !e.IsProduction())
	}
}

func (e *Env) configureNotificationsManager() {
	if e.detectDiscord() {
		return
	}
	if e.detectSlack() {
		return
	}
	e.NotificationManager = &ports.NoNotificationsManager{}
}

func (e *Env) detectSlack() bool {
	wh := h.GetEnv("SLACK_WEBHOOK_URL", "SLACK_WEBHOOK")
	if wh == "" {
		log.Info("no slack webhook found, skipping slack integration")
		return false
	}
	e.NotificationManager = slack.NewSlackNotificationManager(wh)
	return true
}

func (e *Env) detectDiscord() bool {
	wh := h.GetEnv("DISCORD_WEBHOOK_URL", "DISCORD_WEBHOOK")
	if wh == "" {
		log.Info("no discord webhook found, skipping discord integration")
		return false
	}
	e.NotificationManager = discord.NewDiscordNotificationManager(wh)
	return true
}

func (a *App) Handler() http.Handler {
	return a.router.Handler()
}

func (a *App) Start(port ...string) {

	a.scheduler.StartAsync()
	a.router.Start(port...)
}

func (a *App) CleanUp() {
	// sentry.Flush(2 * time.Second)
}

func (a *App) WithDatabase(config DbConfig) {
	ctx := &Ctx{
		Env:  a.Env,
		Auth: nil,
	}
	a.Env.DB = useGorm(config)
	for _, seeder := range config.Seeders {
		if err := seeder(ctx); err != nil {
			log.Fatalf("failed to seed database: %s", err)
		}
	}
}

func (a *App) WithWorkers(workers []*Worker) {
	ctx := &Ctx{
		Env:  a.Env,
		Auth: nil,
	}
	for _, w := range workers {
		_, err := a.scheduler.Every(w.Every).Do(w.Handle, ctx)
		if err != nil {
			log.Fatalf("failed to schedule scheduler: %s", err)
		}
	}
}

func (a *App) WithRouter(config RouterConfig) {
	config.Routes = append(config.Routes, ActuatorRoutes())
	a.router = NewEchoRouter(a.Env, config)

	//a.r.AddRoutes([]func(r Router){ActuatorRoutes()})
	//cb(a.router)
}

/*func (a *App) AddRoutes(routes []func(r Router)) {
	for _, route := range routes {
		route(a.router)
	}
}*/

/*
func (a *App) AddPublicRoute(path ...string) {
	a.router.AddPublicRoute(path...)
}
*/
