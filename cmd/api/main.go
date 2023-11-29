package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	course "github.com/JailtonJunior94/kafka-poc/pkg/v1"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/protobuf/proto"
)

type messageRequest struct {
	Message string `json:"message"`
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/health"))

	router.Post("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		var message messageRequest
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		deliveryChan := make(chan kafka.Event)
		msg := &course.CourseMessage{Id: "Id", Description: "Novo curso"}

		producer, serializer := NewKafkaProducer()
		if err := Publish(serializer, msg, "http-messages", producer, nil, deliveryChan); err != nil {
			log.Println(err.Error())
		}

		go DeliveryReport(deliveryChan)
		producer.Flush(2000)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(message); err != nil {
			log.Fatalf(err.Error())
		}
	})

	server := http.Server{
		ReadTimeout:       time.Duration(10) * time.Second,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "4000"))
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}

func NewKafkaProducer() (*kafka.Producer, serde.Serializer) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "localhost:9092",
		"delivery.timeout.ms": "0",
		"acks":                "1",
		"enable.idempotence":  "false",
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	schemaRegistryClient, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://localhost:8081/"))
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	serializer, err := protobuf.NewSerializer(schemaRegistryClient, serde.ValueSerde, protobuf.NewSerializerConfig())
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	return producer, serializer
}

func Publish(serializer serde.Serializer, message proto.Message, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	payload, err := serializer.Serialize(topic, message)
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		Value: payload,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key: key,
	}

	if err := producer.Produce(msg, deliveryChan); err != nil {
		return err
	}
	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("erro ao enviar mensagem")
				return
			}
			fmt.Println("mensagem enviada: ", ev.TopicPartition)
			return
		}
	}
}
