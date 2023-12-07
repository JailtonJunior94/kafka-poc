import { Kafka } from "kafkajs";
import {
    SchemaRegistry,
    SchemaType,
} from "@kafkajs/confluent-schema-registry";
import { CourseMessage } from './generated/protos/v1/course'

const TOPIC = "poc-schemaregistry";

const kafka = new Kafka({
    brokers: ["localhost:9092"],
});

const registry = new SchemaRegistry({
    host: "http://localhost:8081",
});

const producer = kafka.producer();


const registerSchema = async () => {
    try {
        const schema = `
        syntax = "proto3";

        option go_package = "./pkg/v1/course";
        
        package course;
        
        message CourseMessage {
          string id = 1;
          string description = 2;
        }
      `


        const id = await registry.register({ type: SchemaType.PROTOBUF, schema: schema })
        return id



    } catch (e) {
        console.log(e);
    }
};

const produceToKafka = async (registryId: number, message: any) => {
    try {
        await producer.connect();
        const outgoingMessage = {
            key: message.id,
            value: await registry.encode(registryId, message),
        };

        await producer.send({
            topic: TOPIC,
            messages: [outgoingMessage],
        });

        await producer.disconnect();
    } catch (error) {
        console.log(`There was an error producing the message: ${error}`);
    }
};

const produce = async () => {
    const registryId = await registerSchema();


    const course: CourseMessage = {
        id: "id",
        description: "description"
    }
    const serializedMessage = CourseMessage.encode(course).finish()

    try {
        await produceToKafka(1, serializedMessage)

    } catch (error) {
        console.log(`There was an error producing the message: ${error}`);
    }
};

produce()