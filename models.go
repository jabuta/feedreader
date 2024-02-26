package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jabuta/feedreader/internal/database"
	"github.com/jabuta/feedreader/internal/fetchxml"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}

}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
}

func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(dbFeeds))
	for i := range dbFeeds {
		feeds[i] = databaseFeedToFeed(dbFeeds[i])
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		FeedID:    dbFeedFollow.FeedID,
		UserID:    dbFeedFollow.UserID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeeds []database.FeedFollow) []FeedFollow {
	feeds := make([]FeedFollow, len(dbFeeds))
	for i := range dbFeeds {
		feeds[i] = databaseFeedFollowToFeedFollow(dbFeeds[i])
	}
	return feeds
}

type Post struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func databasePostToPost(dbPost database.Post) Post {
	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		Description: &dbPost.Description.String,
		PublishedAt: &dbPost.PublishedAt.Time,
		FeedID:      dbPost.FeedID,
	}
}

func databasePostsToPosts(dbPosts []database.Post) []Post {
	posts := make([]Post, len(dbPosts))
	for i := range dbPosts {
		posts[i] = databasePostToPost(dbPosts[i])
	}
	return posts
}

func xmlPostToDatabaseCreatePostParams(post fetchxml.Item, feedid uuid.UUID) database.CreatePostParams {
	pubdate, _ := time.Parse(time.RFC1123, post.PubDate)
	return database.CreatePostParams{
		ID:    uuid.New(),
		Title: post.Title,
		Url:   post.Link,
		Description: sql.NullString{
			String: deRefOrEmpty(&post.Description),
			Valid:  isNotNil(&post.Description),
		},
		PublishedAt: sql.NullTime{
			Time:  deRefOrEmpty(&pubdate),
			Valid: isNotNil(&pubdate),
		},
		FeedID: feedid,
	}
}

// var to NullVar Converters
func deRefOrEmpty[T any](val *T) T {
	if val == nil {
		var empty T
		return empty
	}
	return *val
}
func isNotNil[T any](val *T) bool {
	return val != nil
}
