package entity

type Combo struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    ComboNameID int64  `json:"combo_name_id"`
}