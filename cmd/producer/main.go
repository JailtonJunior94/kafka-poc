package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/health"))

	router.Post("/api/send", func(w http.ResponseWriter, r *http.Request) {
		deliveryChan := make(chan kafka.Event)

		producer := NewKafkaProducer()
		if err := Publish("Mensagem via GO", "http-messages", producer, nil, deliveryChan); err != nil {
			log.Println(err.Error())
		}

		go DeliveryReport(deliveryChan)
		producer.Flush(2000)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode("data"); err != nil {
			log.Fatalf(err.Error())
		}
	})

	server := http.Server{
		ReadTimeout:       time.Duration(10) * time.Second,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "6000"))
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "localhost:9094",
		"delivery.timeout.ms": "0",
		"acks":                "1",
		"enable.idempotence":  "false",
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	return producer
}

func Publish(message, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	msg := &kafka.Message{
		Value: []byte(message),
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
