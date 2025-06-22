package entity

type Category struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    BusinessGroupID int64 `json:"business_group_id"`
}