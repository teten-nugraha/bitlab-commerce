E-Commerce Microservices Architecture

https://example.com/microservices-diagram.png

This repository contains a collection of microservices that power our e-commerce platform. Each service is designed to be independently deployable and scalable.
Services Overview
Service	Description	Technology Stack
User Service	Manajemen user dan auth (JWT/OAuth2)	Go, MongoDB, gRPC, JWT
Product Service	Katalog produk dan inventory	Go, MongoDB, gRPC, Redis
Order Service	Proses order dan pembayaran	Go, MongoDB, gRPC, Kafka
Payment Service	Integrasi gateway pembayaran	Go, PostgreSQL, gRPC
Shipping Service	Logistik dan pengiriman	Go, MongoDB, Kafka
Notification Service	Email/SMS/Web notification	Node.js, RabbitMQ
Recommendation Service	Rekomendasi produk (ML)	Python, TensorFlow, gRPC
Service Details
User Service

Responsibilities:

    User registration and authentication

    Profile management

    JWT token generation and validation

    Role-based access control

API Endpoints:

    POST /register - Register new user

    POST /login - User login (returns JWT)

    GET /profile - Get user profile

    PUT /profile - Update user profile

Environment Variables:
env

PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB=user_service
JWT_SECRET=your_jwt_secret_key
KAFKA_BROKERS=localhost:9092

Product Service

Responsibilities:

    Product catalog management

    Inventory tracking

    Stock validation

    Product recommendations base

gRPC Methods:

    ValidateProducts - Validate product stock

    GetProductDetails - Get product information

    UpdateStock - Update product inventory

Environment Variables:
env

GRPC_PORT=50051
MONGO_URI=mongodb://localhost:27017
MONGO_DB=product_service
REDIS_URL=redis://localhost:6379

Order Service

Responsibilities:

    Order creation and processing

    Payment processing coordination

    Order status tracking

    Event publishing for other services

gRPC Methods:

    CreateOrder - Create new order

    ProcessPayment - Initiate payment process

    GetOrderStatus - Check order status

Environment Variables:
env

GRPC_PORT=50052
MONGO_URI=mongodb://localhost:27017
MONGO_DB=order_service
KAFKA_BROKERS=localhost:9092
PRODUCT_SERVICE_ADDR=product-service:50051
PAYMENT_SERVICE_ADDR=payment-service:50053

Payment Service

Responsibilities:

    Payment gateway integration

    Payment processing

    Transaction history

    Refund processing

gRPC Methods:

    CreatePayment - Process payment

    GetPaymentStatus - Check payment status

    ProcessRefund - Initiate refund

Environment Variables:
env

GRPC_PORT=50053
POSTGRES_URL=postgres://user:pass@localhost:5432/payment_service
STRIPE_API_KEY=your_stripe_key
PAYPAL_API_KEY=your_paypal_key

Shipping Service

Responsibilities:

    Order fulfillment

    Shipping logistics

    Delivery tracking

    Carrier integration

Event Subscriptions:

    order.created - Process new orders

    payment.processed - Handle paid orders

Environment Variables:
env

PORT=8082
MONGO_URI=mongodb://localhost:27017
MONGO_DB=shipping_service
KAFKA_BROKERS=localhost:9092
SHIPPING_API_KEY=your_shipping_api_key

Notification Service

Responsibilities:

    Email notifications

    SMS alerts

    Web push notifications

    Notification templates

Event Subscriptions:

    order.created - Send order confirmation

    payment.processed - Send payment receipt

    shipping.updated - Send shipping updates

Environment Variables:
env

PORT=8083
RABBITMQ_URL=amqp://localhost:5672
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=user
SMTP_PASS=pass
TWILIO_SID=your_twilio_sid
TWILIO_TOKEN=your_twilio_token

Recommendation Service

Responsibilities:

    Product recommendations

    Personalized suggestions

    Machine learning models

    Trend analysis

gRPC Methods:

    GetRecommendations - Get personalized recommendations

    TrainModel - Retrain ML models

    UpdateUserPreferences - Update user preference data

Environment Variables:
env

GRPC_PORT=50054
MONGO_URI=mongodb://localhost:27017
MONGO_DB=recommendation_service
ML_MODEL_PATH=/models/product_recommendation.h5

Getting Started
Prerequisites

    Docker and Docker Compose

    Go 1.19+

    Python 3.8+ (for Recommendation Service)

    Node.js 16+ (for Notification Service)

Running Locally

    Clone the repository:
    bash

git clone https://github.com/yourorg/ecommerce-microservices.git
cd ecommerce-microservices

Start the infrastructure:
bash

docker-compose -f docker-compose.infrastructure.yml up -d

Run individual services (example for User Service):
bash

    cd user-service
    go run cmd/server/main.go

Deployment

Each service can be deployed independently using the provided Dockerfiles:
bash

docker build -t user-service -f user-service/Dockerfile .
docker run -d --name user-service -p 8080:8080 user-service

Architecture Principles

    Loose Coupling: Services communicate via well-defined APIs and events

    Event-Driven: Kafka/RabbitMQ for asynchronous communication

    Resilience: Circuit breakers and retry mechanisms

    Observability: Centralized logging and metrics

    CI/CD: Independent deployment pipelines per service

Monitoring and Observability

    Metrics: Prometheus + Grafana

    Logging: ELK Stack

    Tracing: Jaeger

    Health Checks: /health endpoint on each service

API Documentation

Each service provides:

    gRPC service definitions in /proto directory

    Swagger/OpenAPI docs for HTTP endpoints

    Postman collection for testing

Contributing

    Fork the repository

    Create feature branches (git checkout -b feature/AmazingFeature)

    Commit changes (git commit -m 'Add some AmazingFeature')

    Push to the branch (git push origin feature/AmazingFeature)

    Open a Pull Request

License

Distributed under the MIT License. See LICENSE for more information.