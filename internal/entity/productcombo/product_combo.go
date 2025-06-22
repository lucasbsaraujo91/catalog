package entity

type ProductCombo struct {
    ID        int64 `json:"id"`
    ComboID   int64 `json:"combo_id"`
    ProductID int64 `json:"product_id"`
}