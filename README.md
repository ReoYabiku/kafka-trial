# kafka-trial
kubernetesとkafkaで遊んでみる

## 起動手順

### Strimzi(Kafka)クラスターの立ち上げ

minikubeの設定を変更する
```shell
minikube start --momory=4096
```

Strimziのdeploymentを作成
```shell
kubectl create namespace kafka
kubectl create -f ./deployment/strimzi.yaml -n kafka
```

Kafkaのdeploymentを作成
```shell
kubectl create -f ./deployment/kafka-single-node.yaml -n kafka
kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n kafka
```

Kafka Bridgeのserviceを作成
```shell
kubectl create -f ./service/kafka-bridge.yaml -n kafka
```

### アプリケーションの立ち上げ

BFFのdeployment, serviceを作成
```shell
kubectl create -f ./deployment/bff.yaml -n kafka
kubectl expose deployment kafka-bff --type=LoadBalancer -n kafka
```

topicの作成。BFFのPodから実行する
```shell
curl -X POST -H "Content-Type: application/vnd.kafka.v2+json" -d '{"topic_name": "reservation"}' http://my-bridge-bridge-service:8080/admin/topics
```

userのdeployment, serviceを作成
```shell
kubectl create -f ./deployment/user.yaml -n kafka
kubectl expose deployment kafka-user --type=LoadBalancer -n kafka
```

consumerの作成。userのPodから実行する
```shell
curl -X POST -H "Content-Type: application/vnd.kafka.v2+json" -d '{"name": "consumer"}' \
    http://my-bridge-bridge-service:8080/consumers/my-group
curl -X POST -H "Content-Type: application/vnd.kafka.v2+json" -d '{"partitions": [{"topic": "reservation", "partition": 0}]}' \
    http://my-bridge-bridge-service:8080/consumers/my-group/instances/consumer/assignments
```

## 使い方
BFFにリクエストを送信する
```shell
curl -X POST -d '{"user_id": "user100", "event_id": "event100", "seat_count": 10}' http://localhost:8080/reservation
```

userで、メッセージを受信する
```shell
curl localhost:8081/user-info
```
