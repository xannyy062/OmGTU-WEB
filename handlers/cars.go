package handlers

import (
	"CarDealership/database/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type CarsHandler struct {
	DB *pgx.Conn
}

func NewCarsHandler(db *pgx.Conn) *CarsHandler {
	return &CarsHandler{DB: db}
}

// GetAllCars возвращает все автомобили
func (h *CarsHandler) GetAllCars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.DB.Query(ctx,
		"SELECT id, firm, model, year, power, color, price, dealer_id FROM cars")
	if err != nil {
		http.Error(w, "Не удалось извлечь автомобили: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		if err := rows.Scan(&car.ID, &car.Firm, &car.Model, &car.Year,
			&car.Power, &car.Color, &car.Price, &car.DealerID); err != nil {
			http.Error(w, "Не удалось найти автомобиль: "+err.Error(), http.StatusInternalServerError)
			return
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка итерации автомобилей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

// GetCarByID возвращает автомобиль по ID
func (h *CarsHandler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ИСПРАВЛЕНИЕ: Корректно извлекаем ID из пути
	// Путь будет вида /api/cars/1
	path := r.URL.Path
	pathParts := strings.Split(path, "/")

	// Проверяем что у нас минимум 4 части: ["", "api", "cars", "id"]
	if len(pathParts) < 4 {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[3] // Это будет ID
	if idStr == "" {
		http.Error(w, "Необходим ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var car models.Car
	err = h.DB.QueryRow(ctx,
		"SELECT id, firm, model, year, power, color, price, dealer_id FROM cars WHERE id = $1", id).
		Scan(&car.ID, &car.Firm, &car.Model, &car.Year,
			&car.Power, &car.Color, &car.Price, &car.DealerID)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Автомобиль не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}
