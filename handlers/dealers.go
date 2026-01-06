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
		http.Error(w, "Ошибка извлечения диллеров: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dealers []models.Dealer
	for rows.Next() {
		var dealer models.Dealer
		if err := rows.Scan(&dealer.ID, &dealer.Name, &dealer.City,
			&dealer.Address, &dealer.Area, &dealer.Rating); err != nil {
			http.Error(w, "Не удалось прочитать диллеров: "+err.Error(), http.StatusInternalServerError)
			return
		}
		dealers = append(dealers, dealer)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Ошибка итерации диллеров: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dealers)
}

// GetDealerByID возвращает дилера по ID
func (h *DealersHandler) GetDealerByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	path := r.URL.Path
	pathParts := strings.Split(path, "/")

	// Проверяем что у нас минимум 4 части: ["", "api", "dealers", "id"]
	if len(pathParts) < 4 {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[3] // Это будет ID
	if idStr == "" {
		http.Error(w, "Необходим ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID формат: "+err.Error(), http.StatusBadRequest)
		return
	}

	var dealer models.Dealer
	err = h.DB.QueryRow(ctx,
		"SELECT id, name, city, address, area, rating FROM dealers WHERE id = $1", id).
		Scan(&dealer.ID, &dealer.Name, &dealer.City,
			&dealer.Address, &dealer.Area, &dealer.Rating)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Диллер не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dealer)
}
