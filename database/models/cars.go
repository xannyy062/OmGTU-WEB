package models

type Car struct {
	ID       int    `json:"id"`
	Firm     string `json:"firm"`
	Model    string `json:"model"`
	Year     int    `json:"year"`
	Power    int    `json:"power"`
	Color    string `json:"color"`
	Price    int    `json:"price"`
	DealerID int    `json:"dealer_id"`
}
