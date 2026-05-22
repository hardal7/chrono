package user

import (
	"net/http"
	"time"

	"github.com/hardal7/chrono/internal/domains/topic_event"
	e "github.com/hardal7/chrono/internal/util/errors"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost    int = 10
	notRegistered int = 0
)

func Register(w http.ResponseWriter, r *http.Request, rr RegisterRequest) {
	logger.Info("Registering user with username: " + rr.Username)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcryptCost)
	if err != nil {
		ErrHashPassword.Handle(w, err)
		return
	}
	user, err := Repo.FindByUsername(r.Context(), rr.Username)
	if user.ID != notRegistered {
		ErrAlreadyRegistered.Handle(w, err)
		return
		// TODO: Check if this works as well
	} else if err != pgx.ErrNoRows {
		e.ErrCheckIfDuplicate.Handle(w, err, table)
		return
	} else {
		user = User{
			Email:     rr.Email,
			Username:  rr.Username,
			Password:  string(passwordHash),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := Repo.Create(r.Context(), user); err != nil {
			e.ErrCreate.Handle(w, err, table)
			return
		} else {
			logger.Info("Registered user: " + rr.Username)
			// TODO: Handle this error
			user, _ := Repo.FindByUsername(r.Context(), rr.Username)
			err := topicevent.Initialize(user.ID, r.Context())
			if err != nil {
				// TODO
				logger.Info("Failed to initialize topic events for user: " + rr.Username)
				logger.Debug(err.Error())
				logger.Warn("Continuing with errors")
			}
			w.WriteHeader(http.StatusCreated)
		}
	}
}
