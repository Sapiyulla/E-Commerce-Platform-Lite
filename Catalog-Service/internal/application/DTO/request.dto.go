package dto

type Request struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float32 `json:"price"`
}
