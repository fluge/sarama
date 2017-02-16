package sarama

import (
	"testing"
)

func TestActorWorker(t *testing.T){

	kafkaTopic := "air_temperature_10"
	StartProducer(kafkaTopic, 308593)

	StartAverageCalcConsumer(kafkaTopic)
}