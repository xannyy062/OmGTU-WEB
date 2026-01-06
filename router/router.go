package router

import (
	"CarDealership/handlers"
	"net/http"
)

func SetupRoutes(carsHandler *handlers.CarsHandler, dealersHandler *handlers.DealersHandler) {
	// Обработчики для автомобилей
	http.HandleFunc("/api/cars", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Проверяем, запрашивается ли конкретный автомобиль
			if len(r.URL.Path) > len("/api/cars") && r.URL.Path != "/api/cars" {
				// Если путь содержит ID (/api/cars/1)
				carsHandler.GetCarByID(w, r)
			} else {
				// Если путь просто /api/cars
				carsHandler.GetAllCars(w, r)
			}
		case http.MethodPost:
			carsHandler.CreateCar(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})

	// Отдельный обработчик для PUT и DELETE автомобилей
	http.HandleFunc("/api/cars/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			carsHandler.GetCarByID(w, r)
		case http.MethodPut:
			carsHandler.UpdateCar(w, r)
		case http.MethodDelete:
			carsHandler.DeleteCar(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})

	// Обработчики для дилеров
	http.HandleFunc("/api/dealers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Проверяем, запрашивается ли конкретный дилер
			if len(r.URL.Path) > len("/api/dealers") && r.URL.Path != "/api/dealers" {
				dealersHandler.GetDealerByID(w, r)
			} else {
				dealersHandler.GetAllDealers(w, r)
			}
		case http.MethodPost:
			dealersHandler.CreateDealer(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})

	// Отдельный обработчик для PUT и DELETE дилеров
	http.HandleFunc("/api/dealers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dealersHandler.GetDealerByID(w, r)
		case http.MethodPut:
			dealersHandler.UpdateDealer(w, r)
		case http.MethodDelete:
			dealersHandler.DeleteDealer(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})
}
