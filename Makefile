build-api:
	@echo "Compiling API..."
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/api ./cmd/server/main.go

start: 
	docker compose -f deployment/docker-compose.yml up -d --build

stop: 
	docker compose -f deployment/docker-compose.yml down

start-poc-avro: 
	docker compose -f deployment/docker-compose-avro.yml up -d --build

stop-poc-avro: 
	docker compose -f deployment/docker-compose-avro.yml down

proto-gen-v1:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/v1/*.proto
	mv protos/v1/*.go pkg/v1/

proto-gen-v2:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/v2/*.proto
	mv protos/v2/*.go pkg/v2/