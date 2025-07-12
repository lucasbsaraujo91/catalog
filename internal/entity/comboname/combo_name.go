package entity

import "github.com/google/uuid"

type ComboName struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	ComboNameUuid string `json:"combo_name_uuid"`
	Nickname      string `json:"nickname"`
	IsAvailable   bool   `json:"is_available"`
}

// Enable ativa o combo
func (c *ComboName) Enable() {
	c.IsAvailable = true
}

// Disable desativa o combo
func (c *ComboName) Disable() {
	c.IsAvailable = false
}

// GenerateUUID gera e define um UUID único para o combo
func (c *ComboName) GenerateUUID() {
	c.ComboNameUuid = uuid.New().String()
}
