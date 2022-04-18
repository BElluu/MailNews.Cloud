package rss

import (
	"MailNews.Subscriber/models"
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"time"
)

func AmazonRSSParser() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext("https://aws.amazon.com/blogs/aws/feed/", ctx)
	fmt.Println(feed.Title)
	for _, item := range feed.Items {
		var FeedItem = models.FeedItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PublishDate: item.PublishedParsed,
		}
		if FeedItem.PublishDate.Before(time.Now()) { // before last rss check

			fmt.Println("Tytuł: " + FeedItem.Title)
			fmt.Println("Opis: " + FeedItem.Description)
			fmt.Println("Link: " + FeedItem.Link)
			fmt.Println("Data publikacji: " + FeedItem.PublishDate.String())
		}
	}
}
