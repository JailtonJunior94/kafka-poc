package main

import (
	"fmt"
	"log"

	"github.com/JailtonJunior94/kafka-poc/pkg/kafka"
	pix "github.com/JailtonJunior94/kafka-poc/pkg/v3"
)

const (
	topic = "poc-schemaregistry"
)

func main() {
	producer, err := kafka.NewProducer("localhost:9092", "http://localhost:8081/")
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	msg := &pix.PixMessage{Type: "pix_out", Amount: 100}
	offset, err := producer.ProduceMessage(topic, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(offset)
}
