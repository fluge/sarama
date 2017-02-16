package sarama

import (
	"github.com/Shopify/sarama"
	"github.com/lancewf/concurrent"
	"os"
	"os/signal"
	"fmt"
	"strconv"
)

func StartAverageCalcConsumer(topicName string) {
	master := createConsumer()

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	consumer, err := master.ConsumePartition(topicName, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	replyChan := make(chan interface{}, 1)
	actor := NewActorConsumer()

	go func(){
		for msg := range consumer.Messages() {
			actor.Send <- concurrent.Request{*msg, replyChan}
		}
	}()

	go func(){
		for err := range consumer.Errors() {
			actor.Send <- concurrent.Request{*err, replyChan}
		}
	}()

	waitForInterrupt()

	// send request for average
	actor.Send <- concurrent.Request{GetAverageRequest{}, replyChan}

	// wait for reply from actor
	response := <- replyChan

	switch v := response.(type) {
	case AverageResponse:
		fmt.Println("Final average:", strconv.FormatFloat(v.Average, 'E', -1, 64))
	}
}

func waitForInterrupt(){
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	<-signals
	fmt.Println("Interrupt is detected")
}

func createConsumer() sarama.Consumer {
	brokers := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	return master
}
