package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/koderkt/blog_aggregator/internal/database"
)

func (apiConfig *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	reqParams := new(params)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(reqParams)

	feed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqParams.Name,
		Url:       reqParams.Url,
		UserID:    user.ID,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "could not create feeed")
		return
	}

	feedResponse := Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}

	follow, err := apiConfig.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feedResponse.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil{
		log.Println(err)

	}

	type CreateFeedResposne struct{
		Feed Feed `json:"feed"`
		FeedFollow FeedFollow  `json:"feed_follow"`
	}

	respondWithJSON(w, http.StatusCreated, CreateFeedResposne{
		Feed: feedResponse,
		FeedFollow: FeedFollow(follow),
	})

}

func (apiConfig *apiConfig) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConfig.DB.GetAllFeeds(r.Context())

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "could not retrive feeds")
		return
	}
	var feedsResponse []Feed

	for _, feed := range feeds {
		var item Feed
		item.ID = feed.ID
		item.CreatedAt = feed.CreatedAt
		item.UpdatedAt = feed.UpdatedAt
		item.Name = feed.Name
		item.Url = feed.Url
		item.UserID = feed.UserID
		feedsResponse = append(feedsResponse, item)
	}
	respondWithJSON(w, http.StatusAccepted, feedsResponse)
}
