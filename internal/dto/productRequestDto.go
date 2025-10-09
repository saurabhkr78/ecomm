package dto

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryId  uint   `json:"category_id"`
	ImageUrl    string `json:"image_url"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
}

type UpdateStocksRequest struct {
	Stock int `json:"stock"`
}
