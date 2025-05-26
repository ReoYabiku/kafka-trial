## トピックの作成
curl -X POST -H "Content-Type: application/vnd.kafka.v2+json" \
  -d '{"topic_name": "reservation"}' \
  http://my-bridge-bridge-service:8080/admin/topics

## consumerの作成
curl -X POST -X "Content-Type: application/vnd.kafka.v2+json" \
  -d '{"name": "consumer1", "format": "binary", "auto.offset.reset": "earliest", "enable.auto.commit": false, "fetch.min.bytes": 512, "consumer.request.timeout.ms": 30000, "isolation.lebel": "read_commtted"}' \
  http://my-bridge-bridge-service:8080/consumers/my-group

## トピックのアサイン
curl -X POST -X "Content-Type: application/vnd.kafka.v2+json" \
  -d '{"partitions": [{"topic": "reservation", "partition": 0}]}' \
  http://my-bridge-bridge-service:8080/consumers/my-group/instances/consumuer1/assignments
