package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/koderkt/blog_aggregator/internal/database"
)

func (apiConfig *apiConfig) handlerFeedFollowCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Feed_Id uuid.UUID `json:"feed_id"`
	}
	reqBody := new(params)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	follow, err := apiConfig.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    reqBody.Feed_Id,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "could not follow feed")
	}

	respondWithJSON(w, http.StatusCreated, FeedFollow(follow))
}

func (apiConfig *apiConfig) handlerDeleteFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := uuid.Parse(r.PathValue("feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid feed follow id")
		return
	}

	err = apiConfig.DB.DeleteFeedFollow(r.Context(), id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "delete unsuccessful")
	}
	respondWithJSON(w, http.StatusCreated, struct{}{})
}

func (apiConfig *apiConfig) handlerGetAllFeedFollowsForAUser(w http.ResponseWriter, r *http.Request, user database.User) {

	follows, err := apiConfig.DB.GetAllFeedFeedFollowsForAUser(r.Context(), user.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "delete unsuccessful")
	}

	var getFollowFeedsResponse []FeedFollow

	for _, follow := range follows {
		var item FeedFollow
		item.CreatedAt = follow.CreatedAt
		item.UpdatedAt = follow.UpdatedAt
		item.FeedID = follow.FeedID
		item.ID = follow.ID
		item.UserID = follow.UserID
		getFollowFeedsResponse = append(getFollowFeedsResponse, item)
	}
	respondWithJSON(w, http.StatusAccepted, getFollowFeedsResponse)
}
