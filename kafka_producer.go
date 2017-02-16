package sarama

import (
	"github.com/Shopify/sarama"
	"log"
	"github.com/lancewf/sarama/sensorcollection"
	"encoding/json"
	"strconv"
)

type KafkaObsValue struct {
	Value float64
	Date  int64
}

func StartProducer(topicName string, deviceId int) {
	producer := createProducer()
	defer producer.Close()

	sensorCacheDevice := sensorcollection.GetSensorCacheDevice(deviceId)

	if len(sensorCacheDevice.Obs) > 0 {
		sov := sensorCacheDevice.Obs[0]
		kafkaObsValueCollection := mapToKafkaObsValueCollection(sov)

		for _, kafkaObsValue := range kafkaObsValueCollection {
			message := createMessage(kafkaObsValue, topicName)

			sendMessage(message, producer)
		}
	}
}

func sendMessage(message *sarama.ProducerMessage, producer sarama.SyncProducer){
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatal(err)
	}

	println("offset", offset)

	println("partition", partition)
}

func mapToKafkaObsValueCollection(sov sensorcollection.Sov) []KafkaObsValue{
	kafkaObsValueCollection := []KafkaObsValue{}
	numberOfObservations := len(sov.Dates)
	for index := 0; index < numberOfObservations; index++ {
		date := sov.Dates[index]
		value := sov.Values[index]
		kafkaObsValue := KafkaObsValue{value, date}
		kafkaObsValueCollection = append(kafkaObsValueCollection, kafkaObsValue)

	}
	return kafkaObsValueCollection
}

func createMessage(kafkaObsValue KafkaObsValue, topicName string) *sarama.ProducerMessage{

	value, err := json.Marshal(kafkaObsValue)
	if err != nil {
		log.Fatal(err)
	}

	message := &sarama.ProducerMessage{
		Topic:     topicName,
		Partition: 0,
		Key:       sarama.StringEncoder(strconv.FormatInt(kafkaObsValue.Date, 10)),
		Value:     sarama.ByteEncoder(value),
	}

	return message
}

func createProducer() sarama.SyncProducer {
	brokers := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
	//setup relevant config info
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}

	return producer
}
