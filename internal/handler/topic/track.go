package topic

import (
	"net/http"
	"time"

	logger "github.com/hardal7/chrono/internal/util"

	"github.com/hardal7/chrono/internal/model"
	"github.com/hardal7/chrono/internal/repository"
)

func Track(w http.ResponseWriter, r *http.Request, tr model.TrackTopicRequest) {
	logger.Info("Tracking time for topic with name: " + tr.Topic)

	topicUser, err := repository.GetTopicUserByName(r.Context(), tr.Topic, r.Context().Value("userID").(int))
	topicUser.UpdatedAt = time.Now()
	if err != nil {
		// TODO: Create if not exists in DB
		logger.Warn(err.Error())
	} else {
		topicUser.TimeTracked = topicUser.TimeTracked.Add(time.Duration(tr.Time.Unix()))
		repository.Update(r.Context(), topicUser, "topic_users")
		logger.Info("Tracked time")
		w.WriteHeader(http.StatusOK)
	}
}
