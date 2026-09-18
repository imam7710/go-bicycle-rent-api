# 🚲 Bicycle Rental API

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)
![Echo](https://img.shields.io/badge/Echo-Framework-00ACD7)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?logo=postgresql&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-API_Documentation-85EA2D?logo=swagger&logoColor=black)

A RESTful backend API for managing bicycle rental operations, including
authentication, bicycle inventory, bookings, cancellations, wallet top-ups,
vouchers, dynamic rental pricing, and operational reports.

This individual project was developed in August 2026 over approximately one
week during Hacktiv8's Full-Time Golang Backend Development Bootcamp.

## Key Features

- User registration and login with bcrypt password hashing
- JWT-based authentication and protected endpoints
- Bicycle inventory and stock management
- Daily bicycle rental booking
- Booking cancellation
- User balance top-up
- Percentage-based vouchers with maximum discount limits
- Dynamic weekday, weekend, and public-holiday pricing
- Booking and revenue reports
- Swagger API documentation
- Service-layer tests for core business logic
- Third-party API integrations

## Business Rules

The API implements several real-world rental business rules:

- Weekday rental price: IDR 50,000 per day
- Weekend and public-holiday rental price: IDR 75,000 per day
- Users must have sufficient balance before completing a booking
- Vouchers use percentage-based discounts with a maximum discount cap
- Bicycle stock is updated based on booking activity
- Bookings can be cancelled through the booking workflow
- Reports provide booking totals, revenue, and booking status summaries

## External API Integrations

| API | Purpose |
|---|---|
| [Open-Meteo](https://open-meteo.com/) | Provides weather information |
| [Nager.Date](https://date.nager.at/) | Detects public holidays for rental pricing |
| [Frankfurter](https://frankfurter.dev/) | Provides currency conversion data |

## Tech Stack

| Technology | Usage |
|---|---|
| Go | Main programming language |
| Echo | HTTP web framework |
| PostgreSQL | Relational database |
| GORM | ORM and database operations |
| JWT | Authentication and authorization |
| bcrypt | Password hashing |
| Swagger | Interactive API documentation |
| Docker | Application containerization |
| Go Test | Service and business-logic testing |

## Architecture

The application follows a layered architecture to separate HTTP handling,
business logic, and database operations.

```text
Client
  ↓
Handler
  ↓
Use Case / Service
  ↓
Repository
  ↓
PostgreSQL
```

This separation makes the application easier to maintain, test, and extend.

## Project Structure

```text
.
├── config/         # Application and database configuration
├── docs/           # Generated Swagger documentation
├── entity/         # Database entities and request/response models
├── handler/        # HTTP request handlers
├── helper/         # Reusable helper functions
├── middleware/     # JWT authentication middleware
├── repository/     # Database access layer
├── service/        # External API and application services
├── usecase/        # Business logic layer
├── .env.example    # Environment variable example
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
```

## Getting Started

### Prerequisites

Make sure the following software is installed:

- Go
- PostgreSQL
- Git
- Docker (optional)

### Installation

Clone the repository:

```bash
git clone https://github.com/imam7710/go-bicycle-rental-api.git
cd go-bicycle-rental-api
```

Install the dependencies:

```bash
go mod download
```

Create the environment configuration:

```bash
cp .env.example .env
```

Create a PostgreSQL database and update the `.env` file with your local
configuration.

Example configuration:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_database_password
DB_NAME=bicycle_rental

JWT_SECRET=replace_with_your_jwt_secret
```

> Adjust the environment variable names if your configuration package uses
> different names.

Run the application:

```bash
go run main.go
```

The API will be available at:

```text
http://localhost:8080
```

## API Documentation

After starting the application, access the Swagger documentation at:

```text
http://localhost:8080/swagger/index.html
```

Swagger provides information about the available endpoints, request bodies,
authentication requirements, parameters, and response formats.

## Running Tests

Run all available tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Generate a test coverage report:

```bash
go test -cover ./...
```

## Running with Docker

Build the Docker image:

```bash
docker build -t bicycle-rental-api .
```

Run the container:

```bash
docker run --env-file .env -p 8080:8080 bicycle-rental-api
```

Make sure the PostgreSQL database is accessible from the Docker container.

## Project Information

- Development period: August 2026
- Development duration: Approximately one week
- Project type: Individual project
- Status: Completed

## Author

**Muhammad Imam Fadhilah**

Junior Backend Developer specializing in Go backend development.

- GitHub: [@imam7710](https://github.com/imam7710)
- LinkedIn: [Muhammad Imam Fadhilah](https://www.linkedin.com/in/muhammad-imam-fadhilah-63943641a/)