package producer

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"

	"github.com/ozoncp/ocp-remind-api/internal/configuration"
)

// Producer - an interface for send event messages
type Producer interface {
	Send(message Message) error
}

type producer struct {
	prod  sarama.SyncProducer
	topic string
}

// NewProducer - creates a new instance of Producer
func NewProducer() *producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer([]string{configuration.Instance().Kafka.URI()}, config)
	if err != nil {
		log.Err(err).Msg("failed to create Sarama new sync producer")
	}

	return &producer{
		prod:  syncProducer,
		topic: "reminds",
	}
}

func (p *producer) Send(message Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		log.Err(err).Msg("failed marshaling message to json:")
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(p.topic),
		Value:     sarama.StringEncoder(bytes),
		Partition: -1,
		Timestamp: time.Time{},
	}

	_, _, err = p.prod.SendMessage(msg)
	return err
}
