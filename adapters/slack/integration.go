package slack

import (
	"github.com/gahissy/go-micro/ports"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type Accessory struct {
	Type     string `json:"type"`
	ImageUrl string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

type TextType struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Block struct {
	Type      string     `json:"type"`
	Text      TextType   `json:"text"`
	Accessory *Accessory `json:"accessory,omitempty"`
}

type Message struct {
	Blocks []Block `json:"blocks"`
}

type NotificationManager struct {
	ports.NotificationManager
	webHookUrl string
	client     *resty.Client
}

func NewSlackNotificationManager(webhook string) *NotificationManager {
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

	var accessory *Accessory

	if message.Image != "" {
		accessory = &Accessory{
			Type:     "image",
			ImageUrl: message.Image,
			AltText:  message.ImageAlt,
		}
	}

	_, err := s.client.R().
		SetBody(Message{
			Blocks: []Block{
				{
					Type:      "section",
					Text:      TextType{Type: "mrkdwn", Text: message.Text},
					Accessory: accessory,
				},
			},
		}).
		SetContentLength(true). // Dropbox expects this value
		Post(s.webHookUrl)      // for upload Dropbox supports PUT too

	if err != nil {
		log.Error("failed to send slack message -- %v", err)
	}

	return err
}
