package main

import (
	"CarDealership/database/connection"
	"CarDealership/database/importer"
	"CarDealership/database/simple_sql"
	"context"
	"fmt"
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

	carsFile := filepath.Join("cars.json")
	dealersFile := filepath.Join("dealers.json")

	if err := importer.ImportData(ctx, conn, carsFile, dealersFile); err != nil {
		panic(err)
	}

	fmt.Println("Данные успешно импортированы!")

}
