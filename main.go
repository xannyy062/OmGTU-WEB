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

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	conn, err := connection.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}

	rmq, err := messaging.NewRabbitMQ()
	if err != nil {
		log.Fatal("RabbitMQ connection error:", err)
	}
	defer rmq.Close()

	defer conn.Close(ctx)

	fmt.Println("База данных успешно подключена!")

	if err := simple_sql.CreateTable(ctx, conn); err != nil {
		panic(err)
	}

	fmt.Println("Таблицы созданы/проверены!")

	// Автоматически проверяем и импортируем данные при запуске
	importDataIfNeeded(ctx, conn)

	// Хендлеры для cars и для dealers
	carsHandler := handlers.NewCarsHandler(conn)
	carsHandler.Rabbit = rmq

	dealersHandler := handlers.NewDealersHandler(conn)

	// Роутер
	router.SetupRoutes(carsHandler, dealersHandler)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Сервер успешно запущен на http://localhost:8080\n")
	fmt.Println("Доступные эндпоинты:")
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

	log.Fatal(http.ListenAndServe(port, nil))
}

// importDataIfNeeded проверяет, есть ли данные в БД, и импортирует их если таблицы пустые
func importDataIfNeeded(ctx context.Context, conn *pgx.Conn) {
	// Проверяем, есть ли данные в таблице dealers
	var dealerCount int
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM dealers").Scan(&dealerCount)

	// Если произошла ошибка (например, таблица не существует) или таблица пустая
	if err != nil || dealerCount == 0 {
		fmt.Println("Таблицы пустые, начинаю импорт данных...")

		carsFile := filepath.Join("cars.json")
		dealersFile := filepath.Join("dealers.json")

		// Проверяем существование JSON файлов
		if !checkJSONFilesExist(carsFile, dealersFile) {
			log.Fatal("Ошибка: JSON файлы не найдены. Убедитесь, что cars.json и dealers.json существуют в корне проекта")
		}

		if err := importer.ImportData(ctx, conn, carsFile, dealersFile); err != nil {
			log.Fatal("Не удалось импортировать данные:", err)
		}

		fmt.Println("Данные успешно импортированы!")
	} else {
		fmt.Printf("В базе уже есть %d дилеров, импорт не требуется\n", dealerCount)
	}
}

// checkJSONFilesExist проверяет существование JSON файлов
func checkJSONFilesExist(carsFile, dealersFile string) bool {
	if _, err := os.Stat(carsFile); os.IsNotExist(err) {
		fmt.Printf("Файл %s не найден\n", carsFile)
		return false
	}

	if _, err := os.Stat(dealersFile); os.IsNotExist(err) {
		fmt.Printf("Файл %s не найден\n", dealersFile)
		return false
	}

	return true
}
