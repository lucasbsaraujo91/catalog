package combonamehandler

type ComboNameResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Nickname      string `json:"nickname"`
	ComboNameUuid string `json:"combo_name_uuid"`
	IsAvailable   bool   `json:"is_available"`
}

type PaginatedComboNameResponse struct {
	Data       []ComboNameResponse `json:"data"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalItems int64               `json:"total_items"`
}
