package fetcher

import (
	"MailNews.Cloud/Backend/common"
	dbservice "MailNews.Cloud/Database/services"
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
	getLastFetchDate, _ := dbservice.GetConfigValue("LastFetchFeedsDate", client)
	lastFetchDateParsed, _ := time.Parse("02-01-2006 15:01:05", getLastFetchDate.Value)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(feedUrl, ctx)
	for _, item := range feed.Items {
		var News = dbservice.News{
			Title:       item.Title,
			Link:        item.Link,
			PublishDate: item.PublishedParsed,
		}
		if News.PublishDate.Before(lastFetchDateParsed) { //TODO change for prod to before
			dbservice.AddNews(News, client, provider)
		}
	}
}
