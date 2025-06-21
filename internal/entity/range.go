package entity

type Range struct {
    ID          int64   `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    MinValue    float64 `json:"min_value"`
    MaxValue    float64 `json:"max_value"`
}