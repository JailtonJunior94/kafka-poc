using Confluent.Kafka;
using Google.Protobuf;
using Producer.Common;
using Confluent.SchemaRegistry.Serdes;

namespace Producer;

public class ProtoProducer<T> : ProducerBase<T> where T : class, IMessage<T>, new()
{
    public ProtoProducer(string bootstrapServers, string schemaRegistryUrl, string topic)
       : base(bootstrapServers, schemaRegistryUrl, topic)
    {
    }

    public void Build()
    {
        AddSchemaRegistry();
        _producer = new ProducerBuilder<string, T>(_producerConfig).SetValueSerializer(new ProtobufSerializer<T>(_schemaRegistry)).Build();
    }
}
