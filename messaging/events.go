package messaging

import "CarDealership/database/models"

type CarEvent struct {
	EventType string     `json:"eventType"`
	Car       models.Car `json:"car"`
}
