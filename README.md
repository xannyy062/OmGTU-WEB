# Car Dealership Management System

Простая система управления автосалоном с бэкендом на Go и фронтендом на React.

## Требования

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- RabbitMQ

## Быстрый запуск

### 1. Клонирование проекта
```bash
git clone <repository-url>

2. Настройка базы данных

Убедитесь, что PostgreSQL запущен и доступен:

# Проверьте статус PostgreSQL
sudo systemctl status postgresql

# Если не запущен
sudo systemctl start postgresql

3. Запуск rabbitmq
docker-compose up -d

4. Запуск бэкенда (Go сервер)
go run main.go

5. Запуск фронтенда
cd /.frontend/
npm install
npm start

Доступ к приложению
- Фронтенд: http://localhost:3000
- Бэкенд API: http://localhost:8080/api

# Получить все автомобили
curl http://localhost:8080/api/cars

# Получить автомобиль с ID 1
curl http://localhost:8080/api/cars/1
