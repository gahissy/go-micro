package sentry

import (
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

func Configure(dsn string, appVersion string, profile string, debug bool) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Release:     appVersion,
		Environment: profile,
		Debug:       debug,
	})
	if err != nil {
		log.Fatalf("sentry.New: %s", err)
	}
}
