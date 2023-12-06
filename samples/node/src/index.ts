import { Kafka } from "kafkajs";
import {
    SchemaRegistry,
    readAVSCAsync,
} from "@kafkajs/confluent-schema-registry";

const TOPIC = "poc-schemaregistry";

// configure Kafka broker
const kafka = new Kafka({
    brokers: ["localhost:9092"],
});

const registry = new SchemaRegistry({
    host: "http://localhost:8081",
});

const producer = kafka.producer();

declare type MyMessage = {
    id: string;
    value: number;
};

const registerSchema = async () => {
    try {
        const schema = await readAVSCAsync("./avro/schema.avsc");
        const { id } = await registry.register(schema);
        return id;
    } catch (e) {
        console.log(e);
    }
};

const produceToKafka = async (registryId: number, message: MyMessage) => {
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
};

// create the kafka topic where we are going to produce the data
const createTopic = async () => {
    try {
        const topicExists = (await kafka.admin().listTopics()).includes(TOPIC);
        if (!topicExists) {
            await kafka.admin().createTopics({
                topics: [
                    {
                        topic: TOPIC,
                        numPartitions: 1,
                        replicationFactor: 1,
                    },
                ],
            });
        }
    } catch (error) {
        console.log(error);
    }
};

const produce = async () => {
    await createTopic();
    try {
        const registryId = await registerSchema();
        if (registryId) {
            const message: MyMessage = { id: "1", value: 1 };
            registryId && (await produceToKafka(registryId, message));
            console.log(`Produced message to Kafka: ${JSON.stringify(message)}`);
        }
    } catch (error) {
        console.log(`There was an error producing the message: ${error}`);
    }
};

produce()