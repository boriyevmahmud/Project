package events

import (
	"context"
	"fmt"

	conf "github.com/mahmud3253/Project/Post_Service/config"
	"github.com/mahmud3253/Project/Post_Service/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	kafkaReader *kafka.Reader
	log         logger.Logger
}
type KafkaConsumera interface {
	Consume(ctx context.Context, topic string)
}

func (p *KafkaConsumer) Consume(ctx context.Context, topic string) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	connString := fmt.Sprintf("%s:%d", conf.Load().KafkaHost, conf.Load().KafkaPort)
	r := kafka.NewReader(kafka.ReaderConfig{

		Brokers: []string{connString},
		Topic:   topic,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}
