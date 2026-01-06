package loader

import (
	"CarDealership/database/models"
	"encoding/json"
	"os"
)

func LoadCarsFromJSON(filename string) ([]models.Car, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data struct {
		Cars []models.Car `json:"cars"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data.Cars, nil
}

func LoadDealersFromJSON(filename string) ([]models.Dealer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dealers []models.Dealer
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dealers); err != nil {
		return nil, err
	}

	return dealers, nil
}
