# 🛍️ Go Product API

[![Go Version](https://img.shields.io/badge/Go-1.25.2-00ADD8?logo=go)](https://golang.org)
[![CI Status](https://github.com/Schrodinger71/go-product-api/actions/workflows/go-analysis.yml/badge.svg)](https://github.com/Schrodinger71/go-product-api/actions)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-336791?logo=postgresql)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?logo=redis)](https://redis.io)

RESTful API сервис для управления товарами с кэшированием, написанный на Go.

## 🚀 Технологии

- **Backend**: Go (чистый `net/http`)
- **База данных**: PostgreSQL + драйвер `lib/pq`
- **Кэширование**: Redis для ускорения запросов
- **Контейнеризация**: Docker & Docker Compose
- **CI/CD**: GitHub Actions

## 📦 Быстрый старт

### Локальная разработка

```bash
# Клонирование репозитория
git clone https://github.com/Schrodinger71/go-product-api.git
cd go-product-api

# Запуск сервисов
docker-compose up -d
