package models

import "time"

type FeedItem struct {
	UUID        string
	Title       string
	Link        string
	PublishDate *time.Time
	Source      string
	Sent        bool
}
