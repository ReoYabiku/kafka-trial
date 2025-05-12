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

BFFのdeployment, serviceを作成
```shell
kubectl create -f ./deployment/bff.yaml -n kafka
kubectl expose deployment kafka-bff --type=LoadBalancer -n kafka
```
