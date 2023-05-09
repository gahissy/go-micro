package micro

import (
	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func New(name string, version string) *App {
	_ = godotenv.Load()

	env := &Env{
		App: &AppInfo{
			Name:    name,
			Version: version,
		},
	}

	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		log.Infof("SENTRY_DSN detected, configuring")
		initSentry(dsn, env.App)
	}

	cron := gocron.NewScheduler(time.UTC)
	router := NewEchoRouter(env)

	app := &App{
		scheduler: cron,
		router:    router,
		env:       env,
	}

	return app
}

func (a *App) Handler() http.Handler {
	return a.router.Handler()
}

func (a *App) Start(port ...string) {
	a.scheduler.StartAsync()
	a.router.Start(port...)
}

func (a *App) CleanUp() {
	sentry.Flush(2 * time.Second)
}

func (a *App) InitDB(config DatabaseConfig) {
	a.env.DB = useGorm(config)
}

func (a *App) AddWorkers(workers []*Worker) {
	for _, w := range workers {
		_, err := a.scheduler.Every(w.Every).Do(w.Handle, a.env)
		if err != nil {
			log.Fatalf("failed to schedule scheduler: %s", err)
		}
	}
}
func (a *App) Seed(seeders ...func(env *Env) error) {
	for _, seed := range seeders {
		if err := seed(a.env); err != nil {
			log.Fatalf("failed to schedule scheduler: %s", err)
		}
	}
}

func (a *App) WithRouter(cb func(r Router)) {
	cb(a.router)
}

func (a *App) AddRoutes(routes []func(r Router)) {
	for _, route := range routes {
		route(a.router)
	}
}

func (a *App) AddPublicRoute(path ...string) {
	a.router.AddPublicRoute(path...)
}
