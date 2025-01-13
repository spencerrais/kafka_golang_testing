# Real-Time Data Pipeline with Kafka

This project implements a real-time data pipeline using **Kafka**, **Docker Compose**, and **Golang**. The pipeline processes `user-login` events in real-time, generates insights, and publishes results to an output topic called `processed-login`.

All pushes to the repository run through a simple CI pipeline which both runs a linter and ensures that the project compiles. This ensures that nobody (except admins) are able to break the project in production. The repository is also set up in such a manner that the commands being run are entirely separate from the app logic (`/cmd` vs `/internal`).

---

## Features

- **Kafka-Based Real-Time Processing**:
  - Consumes `user-login` events from Kafka.
  - Generates insights, including device usage percentages, app version percentages, and unique user counts.
  - Publishes processed insights to a Kafka output topic every minute for the time that has elapsed since the last published message.

- **Components**:
  - **Kafka & Zookeeper**: Message broker and coordination service.
  - **Python Producer**: Simulates event generation.
  - **Golang Consumer**: Processes events and publishes results.
  - **Docker**: Utilized to build the executable on a minimal platform (Alpine Linux) with minimal footprint.
  - **Docker Compose**: The platform used for local testing with a full Kafka deployment.

- **Scalable and Fault Tolerant**:
  - Kafka’s partitioning allows parallelism and scalability.
  - Docker ensures containerized deployments with fault recovery.
  - Golang is highly efficient, with built-in fault tolerance mechanisms and a sharable executable.
  - GitHub Actions and CI are utilized to ensure that production won't break.

---

## Golang Application Setup

The Golang application is modularly designed for scalability, maintainability, and reusability. Below is a breakdown of the key components:

1. **`cmd/consumer/main.go`**:
  - Entry point for the application.
  - Initializes configuration, logger, Kafka consumer(s), and Kafka producer(s).
  - Starts the consumer to process messages and producer to publish insights.

2. **`internal/config`**:
  - Manages configuration by reading environment variables for Kafka brokers, topics, and group IDs.
  - Provides sensible defaults for local development.
  - Centralizes all configuration logic to ensure consistency across the application.

3. **`internal/kafka`**:
  - **Consumer**: A reusable template for subscribing to topics (e.g., `user-login`) and processing messages with fault tolerance. Delegates the actual processing to a custom handler.
  - **Producer**: A generic implementation for publishing messages to topics (e.g., `processed-login`) with reliable delivery.

4. **`internal/logger`**:
  - Centralized logging outputs logs to both the console and JSON-formatted files.
  - Ensures traceability and structured debugging, making it easier to monitor and troubleshoot the application.

5. **`internal/processing`**:
  - Contains all business logic for processing messages.
  - Decouples domain-specific logic from Kafka-related operations to simplify testing and future extensions.

6. **`internal/models`**:
  - Defines the structure of the messages exchanged between Kafka topics.
  - Provides reusable data mappings and models, ensuring consistency and type safety throughout the application.
  - Example:
     ```go
     type UserLogin struct {
         UserID      string `json:"user_id"`
         AppVersion  string `json:"app_version"`
         DeviceType  string `json:"device_type"`
         IP          string `json:"ip"`
         Locale      string `json:"locale"`
         DeviceID    string `json:"device_id"`
         Timestamp   string `json:"timestamp"`
     }
     ```

---

## Why This Structure?

- **Separation of Concerns**:
  - `internal/kafka` focuses on Kafka interactions (consuming and producing).
  - `internal/processing` encapsulates all business logic, keeping it independent of Kafka.
  - `internal/models` standardizes message structures, reducing redundancy and improving type safety.

- **Scalability**:
  - The Kafka templates support horizontal scaling by adding more partitions and consumer instances.
  - Processing logic is modular, making it easy to add new insights or extend functionality.

- **Maintainability**:
  - Each package has a clear responsibility, making the codebase easier to understand, test, and extend.

This design ensures a clean and modular architecture, ready for real-world scaling and future feature additions.

---

## Setup Instructions

### Prerequisites
- Docker and Docker Compose installed.

### Steps to Run

1. Clone the repository:
  ```bash
  git clone https://github.com/spencerrais/golang_kafka_testing.git
  cd golang_kafka_testing
  ```
2. Build and run the Docker containers:
  ```bash
  docker-compose up --build -d
  ```
3. Verify Kafka topics: (`user-login` and `processed-login`)
  ```bash
  docker exec -it golang_kafka_testing-kafka-1 kafka-topics --bootstrap-server kafka:9092 --list
  ```
4. Watch the logging for the User Login Insights on a minute-by-minute basis (will not populate until the container has been running 1+ minutes):
  ```bash
  docker-compose logs -f my-golang-consumer
  ```
5. Kill the Docker containers:
  ```bash
  docker-compose down
  ```

---

## Next Steps

### How to deploy this application in production

1. Container Orchestration:
  - Utilize Kubernetes to deploy and orchestrate all containers.
  - Enable horizontal scaling.
2. Secret Management:
  - Manage all secrets and environment variables via an external system.
3. Monitoring:
  - Regularly push the log files being created in `logs/` to cloud storage.
  - Integrate a log monitoring tool to examine logs and application metrics.
4. Deployment:
  - Create a CD pipeline using GitHub Actions.
  - Create tests as needed.
5. Security:
  - Ensure that the principle of least access is followed, for everyone.

### Other components to add

1. Channel for Failed Messages:
  - Create a Dead Letter Queue for all messages which cannot be processed.
2. Utilize Schemas and Validation:
  - Enforce data quality throughout the system.
3. Anticipate High Traffic
  - Utilize both rate limiting and load balancing with mutiple producers and consumers.

### Scaling Options

1. Kafka Topics:
  - Split the data into more topics
2. Horizontally Scale:
  - Create more producers and consumers to handle the increased load along with Kafka Partitioning.
3. Process Messages More Efficiently:
  - Perform aggregation where possible
  - Drop Data that is not needed.
  - Batch Data together if SLAs allow.
