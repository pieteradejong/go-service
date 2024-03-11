# Go Service: Message Bus

* Dependencies: Zookeeper, Kafka broker
* Module: Producer, Consumer, HTTP API


# Basic Run Local
Start Zookeeper and Kafka (depending your specific setup)

```bash
$ zookeeper-server-start.sh config/zookeeper.properties
$ kafka-server-start.sh config/server.properties
```


Basic Kafka commands (from `$KAFKA_HOME/bin`):

List topics:
`$ kafka-topics --list --bootstrap-server localhost:9092`

Produce from one terminal tab:
`$ kafka-console-producer --topic testtopic1 --bootstrap-server localhost:9092`

Consume from another tab:
`$ kafka-console-consumer --topic testtopic1 --bootstrap-server localhost:9092`

Run `producer/producer.go` to send a message:
`$ go run producer.go`

Test `sign-service`:
`$ curl -X GET localhost:8080/health`

Send message:
`$ curl -X POST http://localhost:8080/sign -H "Content-Type: application/json" -d '{"message": "hello to sign service"}'`

## Run Docker:
`docker-compose up`

`docker exec kafka-1 kafka-topics --create --topic topictest1 --partitions 4 --replication-factor 2 --if-not-exists --zookeeper zk1:22181,zk2:32181,zk3:42181`


# Ongoing work
* [DONE] Consolidated producer, consumer, API into one Go module.
* [DONE] Create `docker-compose` for Zookeeper and Kafka Broker, run production-relevant setup locally.
* [DONE] Create single topic on Kafka Broker and produce + consume messages.
* [TODO] `POST /sign {message}` is sent to producer
* [TODO] Consumer service reads from Kafka and logs message
* [TODO]: Fix Docker-compose setup - hostnames and ports for ZooKeeper and Kafka
* [OPTIONAL] Add custom `/config/kafka.cfg` for Kafka config.
* [OPTIONAL] Add custom `/config/zookeeper.cfg` for Zookeeper config.

Use case idea:
1) `POST /sign` - message body -> Message Service
2) message service saves message in DB
3) Message svc puts message on Kafka topic
4) Signing service reads message from Kafka topic, signs it, and sends to separate Kafka topic
5) Message service reads signed message from Kafka, and saves message with signature to DB
6) Message service notifies user that message has been signed and saved


Resources:
* Ideas:
    ** https://youtu.be/TAI4ZiKMcfY?si=gowvW6VeQLgH8NsY&t=589


