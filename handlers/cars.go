package handlers

import (
	"CarDealership/database/models"
	"CarDealership/messaging"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarsHandler struct {
	DB     *pgxpool.Pool
	Rabbit *messaging.RabbitMQ
}

func NewCarsHandler(db *pgxpool.Pool) *CarsHandler {
	return &CarsHandler{DB: db}
}

// GetAllCars возвращает все автомобили
func (h *CarsHandler) GetAllCars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем соединение из пула
	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		http.Error(w, "Не удалось получить соединение с БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		"SELECT id, firm, model, year, power, color, price, dealer_id FROM cars ORDER BY id")
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
			http.Error(w, "Ошибка чтения данных автомобиля: "+err.Error(), http.StatusInternalServerError)
			return
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка обработки результатов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

// GetCarByID возвращает автомобиль по ID
func (h *CarsHandler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	path := r.URL.Path
	pathParts := strings.Split(path, "/")

	if len(pathParts) < 4 {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[3]
	if idStr == "" {
		http.Error(w, "Необходим ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем соединение из пула
	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		http.Error(w, "Не удалось получить соединение с БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var car models.Car
	err = conn.QueryRow(ctx,
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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

// CreateCar создает новый автомобиль (POST)
func (h *CarsHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON из тела запроса
	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация обязательных полей
	if car.Firm == "" || car.Model == "" || car.Year == 0 || car.Power == 0 || car.Price == 0 {
		http.Error(w, "Отсутствуют обязательные поля (марка, модель, год, мощность, цена)", http.StatusBadRequest)
		return
	}

	// Получаем соединение из пула
	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		http.Error(w, "Не удалось получить соединение с БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// Вставляем новую запись в БД и получаем ID
	var id int
	err = conn.QueryRow(ctx,
		`INSERT INTO cars (firm, model, year, power, color, price, dealer_id) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) 
		 RETURNING id`,
		car.Firm, car.Model, car.Year, car.Power, car.Color, car.Price, car.DealerID,
	).Scan(&id)

	if err != nil {
		http.Error(w, "Ошибка при создании автомобиля: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем созданную запись для ответа
	createdCar := models.Car{
		ID:       id,
		Firm:     car.Firm,
		Model:    car.Model,
		Year:     car.Year,
		Power:    car.Power,
		Color:    car.Color,
		Price:    car.Price,
		DealerID: car.DealerID,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCar)

	if h.Rabbit != nil {
		h.Rabbit.PublishEvent(messaging.CarEvent{
			EventType: "CREATE",
			Car:       createdCar,
		})
	}
}

// UpdateCar обновляет существующий автомобиль (PUT)
func (h *CarsHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPut {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из пути
	path := r.URL.Path
	pathParts := strings.Split(path, "/")

	if len(pathParts) < 4 {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[3]
	if idStr == "" {
		http.Error(w, "ID обязателен", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Парсим JSON из тела запроса
	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация обязательных полей
	if car.Firm == "" || car.Model == "" || car.Year == 0 || car.Power == 0 || car.Price == 0 {
		http.Error(w, "Отсутствуют обязательные поля (марка, модель, год, мощность, цена)", http.StatusBadRequest)
		return
	}

	// Получаем соединение из пула
	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		http.Error(w, "Не удалось получить соединение с БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// Проверяем существует ли автомобиль
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM cars WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Автомобиль не найден", http.StatusNotFound)
		return
	}

	// Обновляем запись
	result, err := conn.Exec(ctx,
		`UPDATE cars 
		 SET firm = $1, model = $2, year = $3, power = $4, color = $5, price = $6, dealer_id = $7
		 WHERE id = $8`,
		car.Firm, car.Model, car.Year, car.Power, car.Color, car.Price, car.DealerID, id,
	)

	if err != nil {
		http.Error(w, "Ошибка при обновлении автомобиля: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем что запись была обновлена
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Автомобиль не найден", http.StatusNotFound)
		return
	}

	// Получаем обновленную запись для ответа
	updatedCar := models.Car{
		ID:       id,
		Firm:     car.Firm,
		Model:    car.Model,
		Year:     car.Year,
		Power:    car.Power,
		Color:    car.Color,
		Price:    car.Price,
		DealerID: car.DealerID,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCar)

	if h.Rabbit != nil {
		h.Rabbit.PublishEvent(messaging.CarEvent{
			EventType: "UPDATE",
			Car:       updatedCar,
		})
	}
}

// DeleteCar удаляет автомобиль по ID (DELETE)
func (h *CarsHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из пути
	path := r.URL.Path
	pathParts := strings.Split(path, "/")

	if len(pathParts) < 4 {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[3]
	if idStr == "" {
		http.Error(w, "ID обязателен", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем соединение из пула
	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		http.Error(w, "Не удалось получить соединение с БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// Сначала получаем данные автомобиля для RabbitMQ
	var car models.Car
	err = conn.QueryRow(ctx,
		"SELECT id, firm, model, year, power, color, price, dealer_id FROM cars WHERE id = $1",
		id,
	).Scan(&car.ID, &car.Firm, &car.Model, &car.Year,
		&car.Power, &car.Color, &car.Price, &car.DealerID)

	// Проверяем существует ли автомобиль
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM cars WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Автомобиль не найден", http.StatusNotFound)
		return
	}

	// Удаляем запись
	result, err := conn.Exec(ctx, "DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Ошибка при удалении автомобиля: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем что запись была удалена
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Автомобиль не найден", http.StatusNotFound)
		return
	}

	// Возвращаем успешный ответ без тела
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNoContent)

	if h.Rabbit != nil {
		h.Rabbit.PublishEvent(messaging.CarEvent{
			EventType: "DELETE",
			Car:       car,
		})
	}
}
