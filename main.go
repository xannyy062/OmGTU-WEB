package main

import (
	"CarDealership/database/connection"
	"CarDealership/database/importer"
	"CarDealership/database/simple_sql"
	"CarDealership/handlers"
	"CarDealership/messaging"
	"CarDealership/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()

	// Используем пул соединений
	pool, err := connection.CreateConnectionPool(ctx)
	if err != nil {
		log.Fatal("Ошибка создания пула соединений:", err)
	}
	defer pool.Close()

	fmt.Println("✅ База данных успешно подключена!")

	// Получаем одно соединение для создания таблиц
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatal("Ошибка получения соединения:", err)
	}

	if err := simple_sql.CreateTable(ctx, conn.Conn()); err != nil {
		conn.Release()
		panic(err)
	}
	conn.Release()

	fmt.Println("✅ Таблицы созданы/проверены!")

	// Автоматически проверяем и импортируем данные при запуске
	importDataIfNeeded(ctx, pool)

	rmq, err := messaging.NewRabbitMQ()
	if err != nil {
		log.Fatal("RabbitMQ connection error:", err)
	}
	defer rmq.Close()

	// Хендлеры для cars и для dealers
	carsHandler := handlers.NewCarsHandler(pool)
	carsHandler.Rabbit = rmq

	dealersHandler := handlers.NewDealersHandler(pool)

	// Роутер
	router.SetupRoutes(carsHandler, dealersHandler)

	// Оборачиваем все обработчики в CORS middleware
	handler := enableCORS(http.DefaultServeMux)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("🚀 Сервер успешно запущен на http://localhost%s\n", port)
	fmt.Println("🌐 CORS включен для всех доменов")
	fmt.Println("📋 Доступные эндпоинты:")
	fmt.Println("  GET    /api/cars          - Получить список всех машин")
	fmt.Println("  GET    /api/cars/{id}     - Получить автомобиль по его идентификатору")
	fmt.Println("  POST   /api/cars          - Создать новый автомобиль")
	fmt.Println("  PUT    /api/cars/{id}     - Обновить автомобиль по ID")
	fmt.Println("  DELETE /api/cars/{id}     - Удалить автомобиль по ID")
	fmt.Println("  GET    /api/dealers       - Получить всех дилеров")
	fmt.Println("  GET    /api/dealers/{id}  - Получить дилера по его идентификатору")
	fmt.Println("  POST   /api/dealers       - Создать нового дилера")
	fmt.Println("  PUT    /api/dealers/{id}  - Обновить дилера по ID")
	fmt.Println("  DELETE /api/dealers/{id}  - Удалить дилера по ID")

	log.Fatal(http.ListenAndServe(port, handler))
}

// importDataIfNeeded проверяет, есть ли данные в БД, и импортирует их если таблицы пустые
func importDataIfNeeded(ctx context.Context, pool *pgxpool.Pool) {
	// Получаем соединение из пула
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatal("Ошибка получения соединения:", err)
	}
	defer conn.Release()

	// Проверяем, есть ли данные в таблице dealers
	var dealerCount int
	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM dealers").Scan(&dealerCount)

	// Если произошла ошибка (например, таблица не существует) или таблица пустая
	if err != nil || dealerCount == 0 {
		fmt.Println("📥 Таблицы пустые, начинаю импорт данных...")

		carsFile := filepath.Join("cars.json")
		dealersFile := filepath.Join("dealers.json")

		// Проверяем существование JSON файлов
		if !checkJSONFilesExist(carsFile, dealersFile) {
			log.Fatal("❌ Ошибка: JSON файлы не найдены. Убедитесь, что cars.json и dealers.json существуют в корне проекта")
		}

		// Временное соединение для импорта
		importConn, err := pool.Acquire(ctx)
		if err != nil {
			log.Fatal("Ошибка получения соединения для импорта:", err)
		}
		defer importConn.Release()

		if err := importer.ImportData(ctx, importConn.Conn(), carsFile, dealersFile); err != nil {
			log.Fatal("❌ Не удалось импортировать данные:", err)
		}

		fmt.Println("✅ Данные успешно импортированы!")
	} else {
		fmt.Printf("📊 В базе уже есть %d дилеров, импорт не требуется\n", dealerCount)
	}
}

// checkJSONFilesExist проверяет существование JSON файлов
func checkJSONFilesExist(carsFile, dealersFile string) bool {
	if _, err := os.Stat(carsFile); os.IsNotExist(err) {
		fmt.Printf("❌ Файл %s не найден\n", carsFile)
		return false
	}

	if _, err := os.Stat(dealersFile); os.IsNotExist(err) {
		fmt.Printf("❌ Файл %s не найден\n", dealersFile)
		return false
	}

	return true
}
