package discord

import (
	"github.com/gahissy/go-micro/ports"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type NotificationManager struct {
	ports.NotificationManager
	webHookUrl string
	client     *resty.Client
}

type Message struct {
	Content string `json:"content"`
}

func NewDiscordNotificationManager(webhook string) *NotificationManager {
	return &NotificationManager{
		webHookUrl: webhook,
	}
}

func (s *NotificationManager) Send(message ports.Notification) error {
	if s.webHookUrl == "" {
		log.Warn("no notifications manager found, skipping notification")
		return nil
	}
	if s.client == nil {
		s.client = resty.New()
	}

	_, err := s.client.R().
		SetBody(Message{
			Content: message.Text,
		}).
		SetContentLength(true). // Dropbox expects this value
		Post(s.webHookUrl)      // for upload Dropbox supports PUT too

	if err != nil {
		log.Error("failed to send slack message -- %v", err)
	}

	return err
}
