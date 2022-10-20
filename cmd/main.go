package main

import (
	"db-producer/common"
	"db-producer/logger"
	"db-producer/rdb"
	"db-producer/system"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func init() {
	logger.InitLog()
}
func main() {
	server := system.InitSystem()
	rdb.ProcessGetData(server.Pro, server.Schema, server.Rdb)
	for  {
		
	}
	//consumer()
}

func consumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "myGroup",
		"enable.auto.commit": false,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"foo"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)

			break
		}
	}

	c.Close()
}
func producer()  {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.LogError(common.SEND_MESSAGE_FAIL, fmt.Sprintf("%v \n%v", ev.TopicPartition, string(ev.Value)))
				} else {
					logger.LogInfo(common.SEND_MESSAGE, string(ev.Value))
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "foo"
	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries
	p.Flush(15 * 1000)
}