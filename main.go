package main

import (
	"CarDealership/database/connection"
	"CarDealership/database/importer"
	"CarDealership/database/simple_sql"
	"CarDealership/handlers"
	"CarDealership/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	ctx := context.Background()
	conn, err := connection.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}

	defer conn.Close(ctx)

	fmt.Println("База данных успешно подключена!")

	if err := simple_sql.CreateTable(ctx, conn); err != nil {
		panic(err)
	}

	fmt.Println("Создание таблицы прошло успешно!")

	if len(os.Args) > 1 && os.Args[1] == "--import" {
		carsFile := filepath.Join("cars.json")
		dealersFile := filepath.Join("dealers.json")

		if err := importer.ImportData(ctx, conn, carsFile, dealersFile); err != nil {
			log.Fatal("Не удалось импортировать данные:", err)
		}
		fmt.Println("Данные успешно импортированы !")
	}

	// Хендлеры для cars и для dealers
	carsHandler := handlers.NewCarsHandler(conn)
	dealersHandler := handlers.NewDealersHandler(conn)

	// Роутер
	router.SetupRoutes(carsHandler, dealersHandler)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Сервер успешно запущен на http://localhost:8080")
	fmt.Println("Доступные эндпоинты:")
	fmt.Println("GET /api/cars          - Получить список всех машин")
	fmt.Println("GET /api/cars/{id}     - Получить автомобиль по его идентификатору")
	fmt.Println("GET /api/dealers       - Получить всех диллеров")
	fmt.Println("GET /api/dealers/{id}  - Получить диллера по его идентификатору")

	log.Fatal(http.ListenAndServe(port, nil))
}
