package services

import (
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/models"
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/mmcdole/gofeed"
	"time"
)

func FetchFeeds() {
	rssParser("https://aws.amazon.com/blogs/aws/feed/", "Amazon")
	rssParser("https://cloudblog.withgoogle.com/products/gcp/rss/", "Google")
	rssParser("https://azurecomcdn.azureedge.net/en-us/blog/feed/", "Azure")
}

func rssParser(feedUrl, provider string) {
	client := database.CreateLocalClient()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	getLastFetchDate, _ := database.GetConfigValue(ctx, "LastFetchFeedsDate", client)
	lastFetchDate := attributevalue.Unmarshal(getLastFetchDate["Value"], &models.Config{})
	println(lastFetchDate)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(feedUrl, ctx)
	for _, item := range feed.Items {
		var FeedItem = models.FeedItem{
			Title:       item.Title,
			Link:        item.Link,
			PublishDate: item.PublishedParsed,
			Source:      provider,
		}
		if FeedItem.PublishDate.After(time.Now()) { // before last rss check
			database.AddFeed(ctx, FeedItem, client)
		}
	}
}
