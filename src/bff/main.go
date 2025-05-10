package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Reservation struct {
	UserID    string `json:"user_id"`
	EventID   string `json:"event_id"`
	SeatCount int    `json:"seat_count"`
}

func main() {
	http.HandleFunc("POST /reservation", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("invalid request", "error", err.Error())
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		var reservation Reservation
		if err := json.Unmarshal(body, &reservation); err != nil {
			slog.Error("failed to unmarshal request body", "error", err.Error())
			http.Error(w, "failed to unmarshal request body", http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(&reservation)
		if err != nil {
			slog.Error("failed to unmarshal reservation", "error", err.Error())
			http.Error(w, "failed to unmarchal reservation", http.StatusInternalServerError)
			return
		}

		slog.Info("POST /reservation", "reservation", reservation)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(res))
	})

	slog.Info("starting server...")
	http.ListenAndServe(":8080", nil)
}
