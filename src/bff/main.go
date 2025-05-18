package main

import (
	"bff/kafka"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type Reservation struct {
	UserID    string `json:"user_id"`
	EventID   string `json:"event_id"`
	SeatCount int    `json:"seat_count"`
}

type ReservationResponse struct {
	UserID    string `json:"user_id"`
	EventID   string `json:"event_id"`
	SeatCount int    `json:"seat_count"`
	Partition int    `json:"partition"`
	Offset    int    `json:"offset"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

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

		slog.Debug("POST /reservation", "reservation", reservation)

		kc := kafka.New()
		resp, err := kc.Send(res)
		if err != nil {
			slog.Error("failed to produce reservation message", "error", err.Error())
			http.Error(w, "failed to produce reservation message", http.StatusInternalServerError)
			return
		}
		if len(resp.Offsets) == 0 {
			slog.Error("no response from kafka")
			http.Error(w, "no response from kafka", http.StatusInternalServerError)
			return
		}

		resResp := &ReservationResponse{
			UserID:    reservation.UserID,
			EventID:   reservation.EventID,
			SeatCount: reservation.SeatCount,
			Partition: resp.Offsets[0].Partition,
			Offset:    resp.Offsets[0].Offset,
		}

		w.WriteHeader(http.StatusOK)
		// TODO: json.Decodeしてstring(buf)した方が目に優しい
		fmt.Fprintln(w, resResp)
	})

	slog.Info("starting server...")
	http.ListenAndServe(":8080", nil)
}
