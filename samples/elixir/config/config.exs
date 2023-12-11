import Config

config :kaffe,
  producer: [
    heroku_kafka_env: false,
    endpoints: [localhost: 9092],
    # endpoints references [hostname: port]. Kafka is configured to run on port 9092.
    # In this example, the hostname is localhost because we've started the Kafka server
    # straight from our machine. However, if the server is dockerized, the hostname will
    # be called whatever is specified by that container (usually "kafka")
    topics: ["our_topic", "another_topic"], # add a list of topics you plan to produce messages to
  ]
