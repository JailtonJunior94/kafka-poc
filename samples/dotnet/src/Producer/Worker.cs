using Course;

namespace Producer;

public class Worker : BackgroundService
{
    readonly string bootstrapServers = "localhost:9092";
    readonly string schemaRegistryUrl = "http://localhost:8081";
    readonly string topicName = "poc-schemaregistry";
    private readonly ILogger<Worker> _logger;

    public Worker(ILogger<Worker> logger)
    {
        _logger = logger;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        var course = new CourseMessage
        {
            Id = Guid.NewGuid().ToString(),
            Description = "Enviado via C#"
        };


        var producer = new ProtoProducer<CourseMessage>(bootstrapServers, schemaRegistryUrl, topicName);
        producer.Build();

        await producer.ProduceAsync(course);
    }
}
