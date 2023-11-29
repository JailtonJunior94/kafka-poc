package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
)

const (
	nullOffset = -1
)

type KafkaProducer interface {
	ProduceMessage(topic string, msg proto.Message) (int64, error)
	Close()
}

type kafkaProducer struct {
	producer   *kafka.Producer
	serializer serde.Serializer
}

func NewProducer(bootstrapServers, schemaRegistryURL string) (KafkaProducer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   bootstrapServers,
		"delivery.timeout.ms": "0",
		"acks":                "1",
		"enable.idempotence":  "false",
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, err
	}

	schemaRegistryClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(schemaRegistryURL))
	if err != nil {
		return nil, err
	}

	serializer, err := protobuf.NewSerializer(schemaRegistryClient, serde.ValueSerde, protobuf.NewSerializerConfig())
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{
		producer:   producer,
		serializer: serializer,
	}, nil
}

func (p *kafkaProducer) ProduceMessage(topic string, msg proto.Message) (int64, error) {
	kafkaChan := make(chan kafka.Event)
	defer close(kafkaChan)

	payload, err := p.serializer.Serialize(topic, msg)
	if err != nil {
		return nullOffset, err
	}

	if err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          payload,
	}, kafkaChan); err != nil {
		return nullOffset, err
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return int64(ev.TopicPartition.Offset), nil
	case kafka.Error:
		return nullOffset, err
	}

	return nullOffset, nil
}

func (p *kafkaProducer) Close() {
	p.serializer.Close()
	p.producer.Close()
}
