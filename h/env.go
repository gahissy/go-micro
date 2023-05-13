package h

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func IsProduction(profile string) bool {
	return profile == "production" || profile == "staging"
}

func GetEnv(keys ...string) string {
	for _, key := range keys {
		if value := viper.GetString(key); value != "" {
			return value
		}
	}
	return ""
}

func RequireEnv(key string) string {
	value := viper.GetString(key)
	if value == "" {
		log.Fatalf("Missing env varaible: %s", key)
	}
	return value
}
