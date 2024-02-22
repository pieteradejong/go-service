# Go Service: Message Bus

* Dependencies: Zookeeper, Kafka broker
* Module: Producer, Consumer, HTTP API


# Run
Zookeeper and Kafka

```bash
$ zookeeper-server-start.sh config/zookeeper.properties
$ kafka-server-start.sh config/server.properties
```

# Ongoing work
* [DONE] Consolidated producer, consumer, API into one Go module.
* [DONE] Locally install and easily run Zookeeper and Kafka Broker.
* [WIP] Create single topic on Kafka Broker and produce + consume messages.
* [TODO] Add custom `/config/kafka.cfg` for Kafka config.
* 


Resources:
* Ideas:
    ** https://youtu.be/TAI4ZiKMcfY?si=gowvW6VeQLgH8NsY&t=589


