package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/koderkt/blog_aggregator/internal/database"
)

func (cfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = specifiedLimit
	}

	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user")
		return
	}
	var postsRespose []Post

	for _, post := range posts {
		var item Post
		item.ID = post.ID
		item.CreatedAt = post.CreatedAt
		item.UpdatedAt = post.UpdatedAt
		item.Title = post.Title
		item.Url = post.Url
		item.Description = nullStringToStringPtr(post.Description)
		item.PublishedAt = nullTimeToTimePtr(post.PublishedAt)
		item.FeedID = post.FeedID
		postsRespose = append(postsRespose, item)
	}
	respondWithJSON(w, http.StatusOK, postsRespose)
}

func nullTimeToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
