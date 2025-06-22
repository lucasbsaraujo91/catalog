package entity

type ProductConfig struct {
    ID          int64  `json:"id"`
    ProductID   int64  `json:"product_id"`
    Name        string `json:"name"`
    Value       string `json:"value"`
}