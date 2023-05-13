package ports

import (
	log "github.com/sirupsen/logrus"
)

type Notification struct {
	Text     string
	User     string
	Image    string
	ImageAlt string
}

type NotificationManager interface {
	Send(message Notification) error
}

type NoNotificationsManager struct {
	NotificationManager
}

func (n *NoNotificationsManager) Send(message Notification) error {
	log.Warn("no notifications manager found, skipping notification")
	return nil
}
