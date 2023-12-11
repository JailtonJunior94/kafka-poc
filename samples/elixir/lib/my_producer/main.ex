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

  def test() do
    client = ConfluentSchemaRegistry.client(base_url: "http://localhost:8081/")
    case ConfluentSchemaRegistry.get_schema(client, "poc-schemaregistry-value", "1") do
      {:ok, reg} ->
        IO.puts("Registro:")
        IO.inspect(reg)
      {:error, 404, %{"error_code" => 40401}} ->
        IO.puts("Not Found")
      {:error, 404, %{"error_code" => 40403}} ->
        IO.puts("Not Found")
      {:error, code, reason} ->
        IO.puts("Other error: #{code} #{reason}")
    end
  end
end
