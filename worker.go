package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jabuta/feedreader/internal/database"
	"github.com/jabuta/feedreader/internal/fetchxml"
)

func (cfg apiConfig) fetchFeeds(ctx context.Context) {
	ticker := time.NewTicker(cfg.fetchSleep)
	wg := &sync.WaitGroup{}

	for {
		feeds2fetch, err := cfg.DB.GetNextFeedsToFetch(ctx, cfg.numFeeds2fetch)
		if err != nil {
			panic(err)
		}

		for _, feed2fetch := range feeds2fetch {
			wg.Add(1)
			go func(feed database.Feed, wg *sync.WaitGroup) {
				defer wg.Done()
				rss, err := fetchxml.FetchXmlFeed(feed.Url)
				if err != nil {
					log.Println(err)
					return
				}
				cfg.DB.MarkFeedFetched(ctx, feed.ID)
				for i, item := range rss.Channel.Items {
					log.Printf(`Feed: %s - Item #: %v, Title: %s`, rss.Channel.Title, i, item.Title)
				}
			}(feed2fetch, wg)
		}
		wg.Wait()
		log.Println("---------------------------------Finished fetch-------------------------------")
		<-ticker.C
	}
}
