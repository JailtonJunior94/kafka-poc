defmodule ExampleProducer do
  def send_my_message({key, value}, topic) do
    Kaffe.Producer.produce_sync(topic, [{key, value}])
  end

  def send_my_message(key, value) do
    Kaffe.Producer.produce_sync(key, value)
  end

  def send_my_message(value) do
    Kaffe.Producer.produce_sync("sample_key", value)
  end

  def produce(topic \\ "poc-schemaregistry") do
    # Create a client connected to the schema registry
    # client = ConfluentSchemaRegistry.client(base_url: "http://localhost:8081/")

    msg = %Course.CourseMessage{id: "123", description: "teste"}
          |> Course.CourseMessage.json_encode!()

    message_object = %{
      key: "1",
      value: msg,
      headers: [
        {"Content-Type", "application/vnd.kafka.protobuf.v2+json"},
        {"Accept", "application/vnd.kafka.v2+json"}
      ]
    }

    Kaffe.Producer.produce(topic, [message_object])

    # case ConfluentSchemaRegistry.get_schema(client, "poc-schemaregistry-value", "latest") do
    #   {:ok, reg} ->
    #     IO.puts("Registro:")
    #     IO.inspect(reg)
    #     IO.inspect(reg["schema"])

    #   {:error, 404, %{"error_code" => 40401}} ->
    #     IO.puts("Not Found")
    #   {:error, 404, %{"error_code" => 40403}} ->
    #     IO.puts("Not Found")
    #   {:error, code, reason} ->
    #     IO.inspect("Other error:")
    #     IO.inspect(code)
    #     IO.inspect(reason)
    # end
  end
end
