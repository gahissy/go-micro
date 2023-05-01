package micro

import (
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"os"
)

func initSentry(dsn string, release string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Release:     release,
		Environment: os.Getenv("ENV"),
		Debug:       os.Getenv("ENV") != "production",
	})
	if err != nil {
		log.Fatalf("sentry.New: %s", err)
	}
}
