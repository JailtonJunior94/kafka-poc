using Confluent.Kafka;
using Confluent.SchemaRegistry;

namespace Producer.Common;

public abstract class ProducerBase<T> where T : class
{
    protected readonly string _topic;
    protected readonly ProducerConfig _producerConfig;
    protected readonly SchemaRegistryConfig _schemaRegistryConfig;
    protected CachedSchemaRegistryClient _schemaRegistry;
    protected IProducer<string, T> _producer;

    private ProducerBase() { }

    public ProducerBase(string bootstrapServers, string schemaRegistryUrl, string topic)
    {
        _topic = topic;
        
        _producerConfig = new ProducerConfig
        {
            BootstrapServers = bootstrapServers
        };

        _schemaRegistryConfig = new SchemaRegistryConfig
        {
            Url = schemaRegistryUrl
        };
    }

    protected void AddSchemaRegistry()
    {
        _schemaRegistry = new CachedSchemaRegistryClient(_schemaRegistryConfig);
    }

    public async Task ProduceAsync(T message)
    {
        await _producer.ProduceAsync(_topic, new Message<string, T> { Value = message });
    }

    public void Close()
    {
        _producer.Dispose();
        _schemaRegistry.Dispose();
    }
}
