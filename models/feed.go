package models

import "time"

type FeedItem struct {
	Title       string
	Description string
	Link        string
	PublishDate *time.Time
}
