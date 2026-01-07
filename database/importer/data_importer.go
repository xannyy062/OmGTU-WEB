package importer

import (
	"CarDealership/database/loader"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ImportData(ctx context.Context, conn *pgx.Conn, carsFile string, dealersFile string) error {
	// Загружаем дилеров
	dealers, err := loader.LoadDealersFromJSON(dealersFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки дилеров: %v", err)
	}

	// Загружаем машины
	cars, err := loader.LoadCarsFromJSON(carsFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки машин: %v", err)
	}

	// Вставляем дилеров
	dealerIDs := make(map[int]int) // индекс в слайсе -> id в базе
	for i, dealer := range dealers {
		var id int
		err := conn.QueryRow(ctx,
			`INSERT INTO dealers (name, city, address, area, rating) 
             VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			dealer.Name, dealer.City, dealer.Address, dealer.Area, dealer.Rating,
		).Scan(&id)

		if err != nil {
			return fmt.Errorf("ошибка вставки дилера %s: %v", dealer.Name, err)
		}
		dealerIDs[i] = id
	}

	fmt.Printf("✅ Добавлено %d дилеров\n", len(dealers))

	// Вставляем машины, распределяя их по дилерам
	for i, car := range cars {
		// Распределяем машины по дилерам циклически
		dealerIndex := i % len(dealers)
		dealerID := dealerIDs[dealerIndex]

		_, err := conn.Exec(ctx,
			`INSERT INTO cars (firm, model, year, power, color, price, dealer_id) 
             VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			car.Firm, car.Model, car.Year, car.Power, car.Color, car.Price, dealerID,
		)

		if err != nil {
			return fmt.Errorf("ошибка вставки машины %s %s: %v", car.Firm, car.Model, err)
		}
	}

	fmt.Printf("✅ Добавлено %d машин\n", len(cars))
	return nil
}
