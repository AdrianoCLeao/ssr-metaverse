# Server Side Rendering on a Metaverse

A starting point for a metaverse-like server, heavily inspired by [Reticulum](https://github.com/Hubs-Foundation/reticulum).

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Makefile Commands](#commands)


## Features

- **Server Side Rendering (SSR):** Delivers pre-rendered content for better SEO and performance.
- **RESTful API:** Provides endpoints for user authentication, management, and additional features.
- **Swagger Documentation:** Auto-generated API documentation.
- **Dockerized Environment:** Easily build and run the application with Docker and Docker Compose.
- **Makefile Automation:** Simplifies tasks such as building, testing, and generating documentation.

## Prerequisites

- [Go](https://golang.org/) **1.22** or later
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Getting Started
1. **Clone the Repository:**
```bash
git clone https://github.com/AdrianoCLeao/ssr-metaverse.git
cd ssr-metaverse
```

2. **Set Environment Variables:**
Create a `.env` file in the root directory or set the following variables in your environment:
```env
JWT_SECRET=secret
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=metaverse
```

3. **Build and Run Locally:**
```bash
make build
make run
```

The server will be accessible at  [http://localhost:8080](http://localhost:8080).
## Commands

The available commands are:

 - `make build`: 
Build the Docker image using Docker Compose.

 - `make up`: 
Start the containers in detached mode using Docker Compose.

 - `make down`: 
Stop and remove Docker containers.

 - `make logs`: 
Display logs of the "app" container.

 - `make docs`: 
Generate Swagger documentation using the swag CLI.

 - `make clean`: 
Remove the locally generated binary.
