package kafka

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	endpoint     = "http://my-bridge-bridge-service:8080"
	acceptHeader = "application/vnd.kafka.binary.v2+json"
	groupId      = "my-group"
	consumer     = "consumer"
)

type Reservation struct {
	UserID    string `json:"user_id"`
	EventID   string `json:"event_id"`
	SeatCount int    `json:"seat_count"`
}

type PollResponse struct {
	Topic     string `json:"topic"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Partition int    `json:"partition"`
	Offset    int    `json:"offset"`
}

type KafkaClient struct{}

func New() *KafkaClient {
	return &KafkaClient{}
}

func (kc *KafkaClient) Poll() ([]*Reservation, error) {
	client := &http.Client{}

	// TODO: ConfigSetからドメイン名を取得する
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.JoinPath("consumers", groupId, "instances", consumer, "records").String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respContentType := resp.Header.Get("Content-Type")
	slog.Debug("poll", "Content-Type", respContentType)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	slog.Debug("poll", "responseBody", respBody)

	var pollResp []*PollResponse
	err = json.Unmarshal(respBody, &pollResp)
	if err != nil {
		return nil, err
	}

	slog.Debug("poll", "json", pollResp)

	var reservations []*Reservation

	for _, v := range pollResp {
		slog.Debug("before decode", "value", v.Value)

		decoded, err := base64.StdEncoding.DecodeString(v.Value)
		if err != nil {
			return nil, err
		}

		slog.Debug("decode", "value", decoded)

		reservation := &Reservation{}
		if err := json.Unmarshal(decoded, reservation); err != nil {
			return nil, err
		}

		slog.Debug("unmarshal", "value", reservation)

		reservations = append(reservations, reservation)
	}

	return reservations, nil
}
