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
	endpoint    = "http://my-bridge-bridge-service:8080"
	contentType = "application/vnd.kafka.v2+json"
	topic       = "reservation"
)

type KafkaRecord struct {
	Key       string `json:"key,omitempty"`
	Value     string `json:"value"`
	Partition int    `json:"partition,omitempty"`
}

type KafkaSendRequest struct {
	Record []KafkaRecord `json:"record"`
}

type KafkaOffset struct {
	Partition string `json:"partition"`
	Offset    int    `json:"offset"`
}

type KafkaSendResponse struct {
	Offsets []KafkaOffset `json:"offsets"`
}

type KafkaClient struct{}

func New() *KafkaClient {
	return &KafkaClient{}
}

func (kc *KafkaClient) Send(msg []byte) error {
	body := KafkaSendRequest{
		Record: []KafkaRecord{{
			Key:   "key",
			Value: string(msg),
		}},
	}

	buf, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	// TODO: ConfigSetからドメイン名を取得する
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	resp, err := http.Post(u.JoinPath("topic", topic).String(), contentType, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var kafkaOffsets KafkaSendResponse
	err = json.Unmarshal(respBody, &kafkaOffsets)
	if err != nil {
		return err
	}

	slog.Info("kafka send", "response", kafkaOffsets)

	return nil
}
