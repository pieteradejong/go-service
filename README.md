# Go Service: Message Bus

* Dependencies: Zookeeper, Kafka broker
* Module: Producer, Consumer, HTTP API


# Run
Zookeeper and Kafka (depending location of your installation)

```bash
$ zookeeper-server-start.sh config/zookeeper.properties
$ kafka-server-start.sh config/server.properties
```

Run Docker:
`docker-compose up`

`docker exec kafka-1 kafka-topics --create --topic topictest1 --partitions 4 --replication-factor 2 --if-not-exists --zookeeper zk1:22181,zk2:32181,zk3:42181`


# Ongoing work
* [DONE] Consolidated producer, consumer, API into one Go module.
* [DONE] Create `docker-compose` for Zookeeper and Kafka Broker, run production-relevant setup locally.
* [IN PROGRESS] Create single topic on Kafka Broker and produce + consume messages.
  * TODO: fix producer connect
* [OPTIONAL] Add custom `/config/kafka.cfg` for Kafka config.
* [OPTIONAL] Add custom `/config/zookeeper.cfg` for Zookeeper config.


Resources:
* Ideas:
    ** https://youtu.be/TAI4ZiKMcfY?si=gowvW6VeQLgH8NsY&t=589


