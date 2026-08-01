package user

import (
	"encoding/json"
	"net/http"

	"github.com/hardal7/chrono/internal/domains/topicevent"
	e "github.com/hardal7/chrono/internal/util/errors"
)

func GetTopUsers(w http.ResponseWriter, r *http.Request, gr GetTopUsersRequest) {
	users, err := Repo.GetTopUsers(r.Context(), gr.Cursor, gr.Limit)
	if err != nil {
		e.ErrNotFound.Handle(w, err, "users")
		return
	} else {
		response := GetTopUsersResponse{}
		for i, v := range users {
			response.Usernames[i] = v.Username
			response.TotalTimes[i] = v.TotalTime
			response.TodayTimes[i] = topicevent.Repo.GetToday()
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			e.ErrEncodeJSON.Handle(w, err)
		}
	}
}
