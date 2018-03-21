package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	localhost := os.Getenv("KAFKA_CLUSTER")
	if localhost == "" {
		localhost = "localhost"
	}

	broker := flag.String("b", localhost, "Kafka Brokers")
	group := flag.String("g", "test", "Group")
	topicName := flag.String("t", "test", "Topic")
	timeout := flag.Int("s", 60000, "Session Timeout")
	help := flag.Bool("h", false, "Help")
	flag.Parse()

	topic := strings.Fields(*topicName)

	if *help {
		fmt.Println(`
         consumer -b <broker> -g <group> -t <topic> -s<timeout> -h
	     -b : Broker, default localhost
	     -g : Consumer group, default test
	     -t : Topic, default test
	     -s : Timeout, default 60s
	 `)
		os.Exit(0)
	}

	fmt.Printf(" Broker: %s \n Group: %s\n Topic: %s\n\n", *broker, *group, *topicName)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    *broker,
		"group.id":             *group,
		"session.timeout.ms":   *timeout,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(c)
	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topic, nil)

	for {
		done := false
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			done = true
		default:
			ev := c.Poll(500)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				done = true
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
		if done {
			break
		}
	}

	fmt.Printf("Closing Kafka %s consumer\n", *broker)
	c.Close()
}
