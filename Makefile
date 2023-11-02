build-api:
	@echo "Compiling API..."
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/api ./cmd/server/main.go

start: 
	docker compose -f deployment/docker-compose-fc.yml up -d --build

stop: 
	docker compose -f deployment/docker-compose-fc.yml down