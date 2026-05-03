# go-solid-order-system

[![CI](https://github.com/zJunior99/Design_Principles_of_Software-JME/actions/workflows/ci.yml/badge.svg)](https://github.com/zJunior99/Design_Principles_of_Software-JME/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Repositorio de ejemplo para el artículo **"Principios de Diseño de Software en Go: Construyendo Código que Dura"** publicado en Medium.

Implementa los principios **SOLID** aplicados a un sistema de procesamiento de pedidos de e-commerce.

---

## 📐 Principios demostrados

| Principio | Dónde verlo |
|-----------|------------|
| **S** — Single Responsibility | `internal/service`, `internal/discount`, `internal/notifier` |
| **O** — Open/Closed | `internal/discount/discount.go` (Strategy pattern) |
| **L** — Liskov Substitution | `internal/notifier/notifier.go` (Email, SMS, Multi) |
| **I** — Interface Segregation | `internal/repository/repository.go` (Reader/Writer) |
| **D** — Dependency Inversion | `internal/service/order_service.go` (DI via interfaces) |

---


