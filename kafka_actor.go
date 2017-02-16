package sarama

import (
	"github.com/lancewf/concurrent"
	"github.com/Shopify/sarama"
	"encoding/json"
	"strconv"
	"time"
	"fmt"
)

func NewActorConsumer() *concurrent.Actor {
	return concurrent.NewActor(&actorConsumer{0, 0.0})
}

type actorConsumer struct {
	msgCount int64
	total float64
}

type GetAverageRequest struct {
}
type AverageResponse struct {
	Average float64
}

func (a *actorConsumer) Receive(request interface{}, sender chan<- interface{}) {
	switch msg := request.(type) {
	case GetAverageRequest:
		sender <- AverageResponse{a.total / float64(a.msgCount)}
	case sarama.ConsumerMessage:
		a.msgCount++
		var m KafkaObsValue
		json.Unmarshal(msg.Value, &m)
		a.total += m.Value
		s := strconv.FormatFloat(m.Value, 'E', -1, 64)
		date := time.Unix(0, (1000000 * m.Date))
		fmt.Println("Received messages", string(msg.Key), s, date.String())
	case sarama.ConsumerError:
		fmt.Println(msg)
	}
}
