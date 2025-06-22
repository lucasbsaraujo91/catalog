package entity

type Recurrence struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Period      string `json:"period"`
}