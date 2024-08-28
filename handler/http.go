package handler

import (
	"encoding/json"
	"net/http"

	"github.com/khavazh7/matchmaker/internal/matchmaker"
)

type CreateUserRequest struct {
	Name    string  `json:"name"`
	Skill   float64 `json:"skill"`
	Latency float64 `json:"latency"`
}

func CreateUserHandler(matcher *matchmaker.Matcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		player := matchmaker.Player{
			Name:    req.Name,
			Skill:   req.Skill,
			Latency: req.Latency,
		}

		matcher.AddPlayer(player)
		w.WriteHeader(http.StatusCreated)
	}
}
