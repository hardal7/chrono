package user

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/domains/topic"
	"github.com/hardal7/chrono/internal/util/logger"

	"github.com/hardal7/chrono/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost   int    = 10
	defaultTopic string = "General"
)

func Register(w http.ResponseWriter, r *http.Request, rr RegisterRequest) {
	logger.Info("Registering user with username: " + rr.Username)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcryptCost)
	if err != nil {
		logger.Info("Failed to create user: Could not hash password")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// TODO: Not only query by username but also the other fields
	user, err := repository.Find[User](r.Context(), "users", "username", rr.Username)
	if user.ID != 0 {
		logger.Info("User " + rr.Username + " is already registered")
		http.Error(w, "User is already registered", http.StatusConflict)
		return
	} else if err != nil {
		logger.Info("Failed to check if user" + user.Username + " is duplicate")
		logger.Debug(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		user = User{
			Email:     rr.Email,
			Username:  rr.Username,
			Password:  string(passwordHash),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := repository.Create(r.Context(), user, "users"); err != nil {
			logger.Info("Failed to create user: " + rr.Username)
			logger.Debug(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			logger.Info("Registered user: " + rr.Username)
			user, _ = repository.Find[User](r.Context(), "users", "username", rr.Username)
			firstTopic, _ := repository.Find[topic.Topic](r.Context(), "topics", "name", defaultTopic)
			topicEvent := topic.TopicEvent{
				UserID:      user.ID,
				TopicID:     firstTopic.ID,
				TimeTracked: 0,
				Date:        int(time.Now().Unix()),
			}
			repository.Create(r.Context(), topicEvent, "topic_events")
			w.WriteHeader(http.StatusCreated)
		}
	}
}
