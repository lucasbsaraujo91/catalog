package entity

type ComboName struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	IsAvailable bool   `json:"is_available"`
}

// Enable ativa o combo
func (c *ComboName) Enable() {
	c.IsAvailable = true
}

// Disable desativa o combo
func (c *ComboName) Disable() {
	c.IsAvailable = false
}
