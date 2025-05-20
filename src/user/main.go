package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"user/kafka"
)

type UserInfo struct {
	UserId              string `json:"user_id"`
	NewReservationCount int    `json:"new_reservation_count"`
}

type UserResponse struct {
	UserInfo   []UserInfo `json:"user_info"`
	TotalCount int        `json:"total_count"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	http.HandleFunc("GET /user-info", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		kc := kafka.New()
		resp, err := kc.Poll()
		if err != nil {
			slog.Error("failed to poll reservation message", "error", err.Error())
			http.Error(w, "failed to poll reservation message", http.StatusInternalServerError)
			return
		}

		// userResp := &UserResponse{}
		// userInfoMap := make(map[string]int)
		// for _, v := range resp {
		// 	userInfoMap[v]
		// }

		w.WriteHeader(http.StatusOK)
		// TODO: json.Decodeしてstring(buf)した方が目に優しい
		// fmt.Fprintln(w, userResp)
		fmt.Fprintln(w, resp)
	})

	slog.Info("starting server...")
	http.ListenAndServe(":8081", nil)
}
