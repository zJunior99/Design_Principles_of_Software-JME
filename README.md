# Software Design Principles in Go

Demonstration project for the article **"Software Design Principles in Go: Building Code That Lasts"** published on Medium.

Implements the **SOLID** principles applied to an e-commerce order processing system.

## Requirements

- Go 1.23 or higher

## Run

```bash
go run ./cmd/...
```

## Test

```bash
go test ./... -v
```
cl
## Project Structure

```
├── cmd/                  # Entry point
├── internal/
│   ├── domain/           # Order entity
│   ├── service/          # Business logic
│   ├── discount/         # Discount strategies (OCP)
│   ├── notifier/         # Notification channels (LSP)
│   ├── repository/       # Data interfaces (ISP)
│   └── logger/           # Logging abstraction (DIP)
└── .github/workflows/    # CI/CD pipeline
```

## Article

[Software Design Principles in Go: Building Code That Lasts]
(https://medium.com/@jm2022075474/software-design-principles-in-go-building-code-that-lasts-2d86847d166b)

### Estudiante: Junior Mamani Estaña