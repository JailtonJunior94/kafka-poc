package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	consumerGroupID       = "test-consumer"
	defaultSessionTimeout = 6000
	noTimeout             = -1
)

type KafkaConsumer interface {
	Run(messageType protoreflect.MessageType, topic string) error
	Close()
}

type kafkaConsumer struct {
	consumer     *kafka.Consumer
	deserializer *protobuf.Deserializer
}

func NewConsumer(bootstrapServers, schemaRegistryURL string) (KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServers,
		"group.id":           consumerGroupID,
		"enable.auto.commit": false,
	})

	if err != nil {
		return nil, err
	}

	schemaRegistryClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(schemaRegistryURL))
	if err != nil {
		return nil, err
	}

	deserializer, err := protobuf.NewDeserializer(schemaRegistryClient, serde.ValueSerde, protobuf.NewDeserializerConfig())
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		consumer:     consumer,
		deserializer: deserializer,
	}, nil
}

func (c *kafkaConsumer) RegisterMessage(messageType protoreflect.MessageType) error {
	return nil
}

func (c *kafkaConsumer) Run(messageType protoreflect.MessageType, topic string) error {
	if err := c.consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		return err
	}

	if err := c.deserializer.ProtoRegistry.RegisterMessage(messageType); err != nil {
		return err
	}

	for {
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			return err
		}

		msg, err := c.deserializer.Deserialize(topic, kafkaMsg.Value)
		if err != nil {
			return err
		}

		c.handleMessage(msg, int64(kafkaMsg.TopicPartition.Offset))
		if _, err = c.consumer.CommitMessage(kafkaMsg); err != nil {
			return err
		}
	}
}

func (c *kafkaConsumer) handleMessage(message interface{}, offset int64) {
	fmt.Printf("message %v with offset %d\n", message, offset)
}

func (c *kafkaConsumer) Close() {
	if err := c.consumer.Close(); err != nil {
		log.Fatal(err)
	}
	c.deserializer.Close()
}
