package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	// load the config file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// checking arguments
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <group> <topics..>\n",
			os.Args[0])
		os.Exit(1)
	}
	// getting topic
	topic := os.Getenv("TOPIC")
	topics := []string{topic}
	// geting group
	group := os.Getenv("GROUPLABEL")
	// getting broker address
	broker := os.Args[1]
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("%v\n,%v\n ", topics, broker)
	// creating consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"broker.address.family": "v4",
		"group.id":              group,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(10)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				commitOffset(c, topic, int(e.TopicPartition.Partition))
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}

func commitOffset(c *kafka.Consumer, topic string, partition int) {
	committedOffsets, err := c.Committed([]kafka.TopicPartition{{
		Topic:     &topic,
		Partition: int32(partition),
	}}, 5000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch offset: %s\n", err)
		os.Exit(1)
	}

	committedOffset := committedOffsets[0]

	fmt.Printf("Committed partition %d offset: %d", committedOffset.Partition, committedOffset.Offset)
}
