package simple_sql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTable(ctx context.Context, conn *pgx.Conn) error {

	sqlQueryDealers := `
	CREATE TABLE IF NOT EXISTS dealers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		city VARCHAR(100) NOT NULL,
		address VARCHAR(100) NOT NULL,
		area VARCHAR(100) NOT NULL,
		rating DECIMAL(3,1) NOT NULL CHECK (rating >= 0 AND rating <= 5)
	);
	`

	sqlQueryCars := `
	CREATE TABLE IF NOT EXISTS cars (
		id SERIAL PRIMARY KEY,
		firm VARCHAR(100) NOT NULL,
		model VARCHAR(100) NOT NULL,
		year INTEGER NOT NULL,
		power INTEGER NOT NULL,
		color VARCHAR(100) NOT NULL,
		price INTEGER NOT NULL,
		dealer_id INTEGER REFERENCES dealers(id) ON DELETE CASCADE
	);
	`

	if _, err := conn.Exec(ctx, sqlQueryDealers); err != nil {
		return err
	}

	if _, err := conn.Exec(ctx, sqlQueryCars); err != nil {
		return err
	}

	return nil
}
