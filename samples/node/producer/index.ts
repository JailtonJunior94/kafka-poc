const { Kafka } = require('kafkajs')
const protobuf = require('protobufjs')

const kafka = new Kafka({
    clientId: 'my-app',
    brokers: ['10.8.0.1:9092']
})

const producer = kafka.producer()

const run = async () => {
    await producer.connect()

    protobuf.load('data.proto', async (err, root) => {
        console.log("TESTING")
        console.log(err)

        let searchRequest = root.lookupType('awesomepackage.SearchRequest')
        let payload = { query: "test", page: 2 }

        let errMsg = searchRequest.verify(payload);
        console.log(errMsg)

        let msg = searchRequest.create(payload)
        let buffer = searchRequest.encode(msg).finish();
        console.log(buffer)

        await producer.send({
            topic: 'test-topic',
            messages: [
                { key: 'key1', value: buffer }
            ]
        })
    })
}

run()