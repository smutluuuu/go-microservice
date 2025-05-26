# Go Microservices Project

A comprehensive microservices application built with Go, Docker, and Kubernetes.

## Technologies Used

- **Go**: Backend services
- **Docker**: Containerization
- **Kubernetes**: Container orchestration
- **RabbitMQ**: Message broker
- **PostgreSQL**: SQL database
- **MongoDB**: NoSQL database
- **gRPC**: Internal service communication
- **RESTful API**: External communication

## Architecture

This project consists of the following microservices:

1. **Broker Service**: API Gateway that routes requests to appropriate services
2. **Authentication Service**: Handles user authentication
3. **Logger Service**: Centralized logging with MongoDB
4. **Mail Service**: Handles email sending
5. **Listener Service**: Consumes messages from RabbitMQ

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Kubernetes (minikube or Docker Desktop)
- Make

### Running with Docker Compose

```bash
make ./project/up_build
```

### Running with Kubernetes

```bash
# Apply Kubernetes manifests
kubectl apply -f ./project/k8s

# Start minikube
minikube start

# Expose Broker Service
kubectl expose deployment broker-service --type=LoadBalancer --port=8080 --target-port=8080

# To make accessible LoadBalancer
minikube tunnel

# To see minikube dashboard
 minikube dashboard
#
```

## Project Structure

- `/authentication-service`: User authentication
- `/broker-service`: API Gateway
- `/front-end`: Web interface
- `/logger-service`: Logging service
- `/mail-service`: Email service
- `/listener-service`: RabbitMQ consumer
- `/k8s`: Kubernetes manifests

## Features

- User authentication
- Message broker integration
- Centralized logging
- Email sending
- gRPC and REST API communication

## License

MIT
