package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	config, err := GetConfigFromFile()
	if err != nil {
		fmt.Println("err load config file " + err.Error())
		panic(err)
	}

	var (
		topic      string   = "my-topic"
		brokersUrl []string = []string{config.Kafka.Hostname + ":" + config.Kafka.Port}
	)
	// connect
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Consumer.Return.Errors = true

	go func() {
		// consumers
		conn, err := sarama.NewConsumer(brokersUrl, saramaConfig)
		if err != nil {
			panic(err)
		}

		consumer, err := conn.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		fmt.Println("Consumer started")
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
				consumer.Close()
			case msg := <-consumer.Messages():
				fmt.Printf("Received message: topic: %v: value: %v \n", string(msg.Topic), string(msg.Value))
			}
		}
	}()

	conn, err := sarama.NewSyncProducer(brokersUrl, saramaConfig)
	if err != nil {
		panic(err)
	}

	for i := 1; i < 10; i++ {
		// send message
		partition, offset, err := conn.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("hello world %d", i)),
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
