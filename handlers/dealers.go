package handlers

import (
	"CarDealership/database/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type DealersHandler struct {
	DB *pgx.Conn
}

func NewDealersHandler(db *pgx.Conn) *DealersHandler {
	return &DealersHandler{DB: db}
}

// GetAllDealers возвращает всех дилеров
func (h *DealersHandler) GetAllDealers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.DB.Query(ctx,
		"SELECT id, name, city, address, area, rating FROM dealers")
	if err != nil {
		http.Error(w, "Ошибка при получении дилеров: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dealers []models.Dealer
	for rows.Next() {
		var dealer models.Dealer
		if err := rows.Scan(&dealer.ID, &dealer.Name, &dealer.City,
			&dealer.Address, &dealer.Area, &dealer.Rating); err != nil {
			http.Error(w, "Ошибка при чтении данных дилера: "+err.Error(), http.StatusInternalServerError)
			return
		}
		dealers = append(dealers, dealer)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка при обработке результатов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dealers)
}

// GetDealerByID возвращает дилера по ID
func (h *DealersHandler) GetDealerByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	var dealer models.Dealer
	err = h.DB.QueryRow(ctx,
		"SELECT id, name, city, address, area, rating FROM dealers WHERE id = $1", id).
		Scan(&dealer.ID, &dealer.Name, &dealer.City,
			&dealer.Address, &dealer.Area, &dealer.Rating)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Дилер не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dealer)
}

// CreateDealer создает нового дилера (POST)
func (h *DealersHandler) CreateDealer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON из тела запроса
	var dealer models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&dealer); err != nil {
		http.Error(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация обязательных полей
	if dealer.Name == "" || dealer.City == "" || dealer.Address == "" {
		http.Error(w, "Отсутствуют обязательные поля (название, город, адрес)", http.StatusBadRequest)
		return
	}

	// Валидация рейтинга (если передан)
	if dealer.Rating < 0 || dealer.Rating > 5 {
		http.Error(w, "Рейтинг должен быть от 0 до 5", http.StatusBadRequest)
		return
	}

	// Вставляем новую запись в БД и получаем ID
	var id int
	err := h.DB.QueryRow(ctx,
		`INSERT INTO dealers (name, city, address, area, rating) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id`,
		dealer.Name, dealer.City, dealer.Address, dealer.Area, dealer.Rating,
	).Scan(&id)

	if err != nil {
		http.Error(w, "Ошибка при создании дилера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем созданную запись для ответа
	createdDealer := models.Dealer{
		ID:      id,
		Name:    dealer.Name,
		City:    dealer.City,
		Address: dealer.Address,
		Area:    dealer.Area,
		Rating:  dealer.Rating,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDealer)
}

// UpdateDealer обновляет существующего дилера (PUT)
func (h *DealersHandler) UpdateDealer(w http.ResponseWriter, r *http.Request) {
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
	var dealer models.Dealer
	if err := json.NewDecoder(r.Body).Decode(&dealer); err != nil {
		http.Error(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация обязательных полей
	if dealer.Name == "" || dealer.City == "" || dealer.Address == "" {
		http.Error(w, "Отсутствуют обязательные поля (название, город, адрес)", http.StatusBadRequest)
		return
	}

	// Валидация рейтинга (если передан)
	if dealer.Rating < 0 || dealer.Rating > 5 {
		http.Error(w, "Рейтинг должен быть от 0 до 5", http.StatusBadRequest)
		return
	}

	// Проверяем существует ли дилер
	var exists bool
	err = h.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM dealers WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Дилер не найден", http.StatusNotFound)
		return
	}

	// Обновляем запись
	result, err := h.DB.Exec(ctx,
		`UPDATE dealers 
		 SET name = $1, city = $2, address = $3, area = $4, rating = $5
		 WHERE id = $6`,
		dealer.Name, dealer.City, dealer.Address, dealer.Area, dealer.Rating, id,
	)

	if err != nil {
		http.Error(w, "Ошибка при обновлении дилера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем что запись была обновлена
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Дилер не найден", http.StatusNotFound)
		return
	}

	// Получаем обновленную запись для ответа
	updatedDealer := models.Dealer{
		ID:      id,
		Name:    dealer.Name,
		City:    dealer.City,
		Address: dealer.Address,
		Area:    dealer.Area,
		Rating:  dealer.Rating,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedDealer)
}

// DeleteDealer удаляет дилера по ID (DELETE)
func (h *DealersHandler) DeleteDealer(w http.ResponseWriter, r *http.Request) {
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

	// Проверяем существует ли дилер
	var exists bool
	err = h.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM dealers WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Дилер не найден", http.StatusNotFound)
		return
	}

	// Удаляем запись
	result, err := h.DB.Exec(ctx, "DELETE FROM dealers WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Ошибка при удалении дилера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем что запись была удалена
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Дилер не найден", http.StatusNotFound)
		return
	}

	// Возвращаем успешный ответ без тела
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNoContent)
}
