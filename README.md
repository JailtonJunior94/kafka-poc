# Escopo
Criação de um Kafka Connect para enviar notificações com HTTP Sink

[Documentação de Connectors](https://docs.confluent.io/cloud/current/connectors/index.html) <br>
[Documentação do HTTP Sink](https://docs.confluent.io/kafka-connectors/http/current/overview.html)

# Desenho da Solução
<p align="center">
  <img src="docs/kafka-poc.png" width="800" title="Main">
</p>

# Executando com Docker
- Para executar o projeto local com docker, devemos utilizar os comandos
  ```
  make start
  ```
- Para parar a execução do projeto

  ```
  make stop
  ```

# Configurando um novo connector 
- Podemos criar um arquivo `.properties` como no exemplo abaixo e fazer upload do mesmo no `control-center`
- Exemplo de arquivo `.properties`
  
  ```
  name=HttpSink
  connector.class=io.confluent.connect.http.HttpSinkConnector
  tasks.max=1
  value.converter=org.apache.kafka.connect.storage.StringConverter
  topics=http-messages
  http.api.url=http://api:3000/api/messages
  request.method=post
  confluent.topic.bootstrap.servers=kafka:9092
  confluent.topic.replication.factor=1
  reporter.bootstrap.servers=kafka:9092
  reporter.result.topic.name=success-responses
  reporter.result.topic.replication.factor=1
  reporter.error.topic.name=error-responses
  reporter.error.topic.replication.factor=1
  ```
- Podemos também usar a API Rest disponibilizada pelo connectors com os seguintes endpoints 
- Listando connectors instalados no Kafka
  ```
  curl --location 'http://localhost:8083/connector-plugins'
  ```
- Listando connectors criados 
  ```
  curl --location 'http://localhost:8083/connectors'
  ```
- Cadastrando um novo connector
  ```
    curl --location 'http://localhost:8083/connectors' \
  --header 'Content-Type: application/json' \
  --data '{
    "name": "HttpSink",
    "config": {
      "topics": "http-messages",
      "tasks.max": "1",
      "connector.class": "io.confluent.connect.http.HttpSinkConnector",
      "http.api.url": "http://api:3000/api/messages",
      "value.converter": "org.apache.kafka.connect.storage.StringConverter",
      "value.converter.schemas.enable": false,
      "confluent.topic.bootstrap.servers": "kafka:9092",
      "confluent.topic.replication.factor": "1",
      "reporter.bootstrap.servers": "kafka:9092",
      "reporter.result.topic.name": "success-responses",
      "reporter.result.topic.replication.factor": "1",
      "reporter.error.topic.name":"error-responses",
      "reporter.error.topic.replication.factor":"1"
    }
  }'
  ```
- Verificando status do connector cadastrado
  ```
  curl --location 'http://localhost:8083/connectors/HttpSink/status'
  ```
- Deletando um connector 
  ```
  curl --location --request DELETE 'http://localhost:8083/connectors/HttpSink'
  ```