package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var start int64

func main() {

	localhost := os.Getenv("KAFKA_CLUSTER")
	if localhost == "" {
		localhost = "localhost"
	}
	topic := "test-topic"
	brokers := flag.String("b", localhost, "Kafka Brokers")
	group := flag.String("g", "test", "Group")
	help := flag.Bool("h", false, "Help")
	flag.Parse()

	config := &kafka.ConfigMap{
		"metadata.broker.list": *brokers,
		"group.id": *group,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
		"debug":                           "generic,broker,security",
	}

	if *help {
		fmt.Println(`
         producer -b <broker> -g <group> -t <topic>
	     -b : Broker, default localhost:9092 or localhost:9092,localhost:9093
	 `)
		os.Exit(0)
	}

	p, err := kafka.NewProducer(config)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	deliveryChan := make(chan bool)

	go func() {
		defer close(deliveryChan)
		var ix = 0
		for e := range p.Events() {
			ix++
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					//fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					//		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
					if ix == 50000 {
						var duration = time.Now().UnixNano()/int64(time.Millisecond) - start
						var rate = int64(ix) / duration * 1000
						fmt.Printf("message sent %d, duration %d ms, rate=%d mess/s\n", ix, duration, rate)
					}
				}
				//return

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	value := "Hello Go!"
	start = time.Now().UnixNano() / int64(time.Millisecond)
	ix := 0
	for ix = 0; ix <= 50002; ix++ {
		p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}
	}
	_ = <-deliveryChan
	p.Close()

}
