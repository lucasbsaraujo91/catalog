package entity

type Product struct {
    ID          int64   `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    CategoryID  int64   `json:"category_id"`
    ProductTypeID int64 `json:"product_type_id"`
}