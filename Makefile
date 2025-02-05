.PHONY: help build up down logs docs

# Display available commands.
help:
	@echo "Available commands:"
	@echo "  make build 	    -> Build the Docker image using docker-compose"
	@echo "  make up      	    -> Start the containers in detached mode using docker-compose"
	@echo "  make down    	    -> Stop and remove containers using docker-compose"
	@echo "  make logs          -> Display logs of the 'app' container"
	@echo "  make swag          -> Generate Swagger documentation (swag init)"
	@echo "  make clean         -> Remove the locally generated binary"

# Build the Docker image using docker-compose.
build:
	docker-compose build

# Start the containers defined in docker-compose in detached mode.
up:
	docker-compose up -d

# Stop and remove the containers.
down:
	docker-compose down

# Display the logs of the "app" container.
logs:
	docker-compose logs -f app

# Generate the Swagger documentation using docs.
docs:
	swag init -g cmd/main.go --dir . --exclude vendor,assets,docs -o ./api/swagger
