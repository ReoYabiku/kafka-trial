package kafka

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	endpoint = "http://my-bridge-bridge-service:8080"
	// endpoint    = "http://localhost:3000"
	contentType = "application/vnd.kafka.json.v2+json"
	topic       = "reservation"
)

type KafkaRecord struct {
	Key       string `json:"key,omitempty"`
	Value     string `json:"value"`
	Partition int    `json:"partition,omitempty"`
}

type KafkaSendRequest struct {
	Records []KafkaRecord `json:"records"`
}

type KafkaOffset struct {
	Partition int `json:"partition"`
	Offset    int `json:"offset"`
}

type KafkaSendResponse struct {
	Offsets []KafkaOffset `json:"offsets"`
}

type KafkaClient struct{}

func New() *KafkaClient {
	return &KafkaClient{}
}

func (kc *KafkaClient) Send(msg []byte) (*KafkaSendResponse, error) {
	body := KafkaSendRequest{
		Records: []KafkaRecord{{
			Key:   "key",
			Value: string(msg),
		}},
	}

	buf, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	slog.Debug("requestBody", "json", string(buf))

	// TODO: ConfigSetからドメイン名を取得する
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(u.JoinPath("topics", topic).String(), contentType, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	slog.Debug("topic", "body", respBody)

	var kafkaOffsets KafkaSendResponse
	err = json.Unmarshal(respBody, &kafkaOffsets)
	if err != nil {
		return nil, err
	}

	slog.Debug("kafka send", "response", kafkaOffsets)

	return &kafkaOffsets, nil
}
