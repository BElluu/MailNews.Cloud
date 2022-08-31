package services

import (
	"MailNews.Subscriber/common"
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/models"
	"context"
	"github.com/mmcdole/gofeed"
	"time"
)

func FetchFeeds() {
	rssParser("https://aws.amazon.com/blogs/aws/feed/", "Aws")
	rssParser("https://cloudblog.withgoogle.com/products/gcp/rss/", "Google")
	rssParser("https://azurecomcdn.azureedge.net/en-us/blog/feed/", "Azure")
}

func rssParser(feedUrl, provider string) {
	client := common.CreateLocalClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	getLastFetchDate, _ := database.GetConfigValue("LastFetchFeedsDate", client)
	lastFetchDateParsed, _ := time.Parse("02-01-2006 15:01:05", getLastFetchDate.Value)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(feedUrl, ctx)
	for _, item := range feed.Items {
		var FeedItem = models.FeedItem{
			Title:       item.Title,
			Link:        item.Link,
			PublishDate: item.PublishedParsed,
			//Provider:    provider,
		}
		if FeedItem.PublishDate.Before(lastFetchDateParsed) { // change for prod to before
			database.AddFeed(FeedItem, client, provider)
		}
	}
}
