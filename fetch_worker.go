package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jabuta/feedreader/internal/database"
	"github.com/jabuta/feedreader/internal/fetchxml"
	"github.com/lib/pq"
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
				for _, item := range rss.Channel.Items {
					createdPost, err := cfg.DB.CreatePost(ctx, xmlPostToDatabaseCreatePostParams(item, feed.ID))
					if err != nil {
						if pqErr, ok := err.(*pq.Error); !(ok && pqErr.Code == "23505") {
							log.Println(err)
						} else {
							continue
						}
					}
					log.Println(createdPost)
				}
			}(feed2fetch, wg)
		}
		wg.Wait()
		log.Println("---------------------------------Finished fetch-------------------------------")
		<-ticker.C
	}
}
