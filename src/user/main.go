package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	kc := kafka.New()
	for {
		resp, _ := kc.Poll()
		jsonResp, _ := json.MarshalIndent(resp, "", "  ")
		if string(jsonResp) != "null" {
			fmt.Println(string(jsonResp))
		}
	}
}
