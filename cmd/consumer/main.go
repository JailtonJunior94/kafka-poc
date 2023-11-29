package main

import (
	"log"

	"github.com/JailtonJunior94/kafka-poc/pkg/kafka"
	course "github.com/JailtonJunior94/kafka-poc/pkg/v1"
)

const (
	topic = "poc-schemaregistry"
)

func main() {
	consumer, err := kafka.NewConsumer("localhost:9092", "http://localhost:8081/")
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	messageType := (&course.CourseMessage{}).ProtoReflect().Type()
	if err := consumer.Run(messageType, topic); err != nil {
		log.Fatal(err)
	}
}
