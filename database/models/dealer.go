package models

type Dealer struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	City    string  `json:"city"`
	Address string  `json:"address"`
	Area    string  `json:"area"`
	Rating  float64 `json:"rating"`
}
