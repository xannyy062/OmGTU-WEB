package router

import (
	"CarDealership/handlers"
	"net/http"
)

func SetupRoutes(carsHandler *handlers.CarsHandler, dealersHandler *handlers.DealersHandler) {
	// Важно: оба обработчика должны быть зарегистрированы
	// Один для списка, другой для конкретного элемента
	http.HandleFunc("/api/cars", carsHandler.GetAllCars)
	http.HandleFunc("/api/cars/", carsHandler.GetCarByID) // Обратите внимание на слеш в конце!
	http.HandleFunc("/api/dealers", dealersHandler.GetAllDealers)
	http.HandleFunc("/api/dealers/", dealersHandler.GetDealerByID) // И здесь тоже
}
