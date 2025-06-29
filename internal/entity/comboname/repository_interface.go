package entity

import "context"

// Repository define as operações para persistência de ComboName
type Repository interface {
	Create(ctx context.Context, combo *ComboName) (int64, error)
	GetByID(ctx context.Context, id int64) (*ComboName, error)
	GetAll(ctx context.Context, page, limit int) ([]ComboName, int64, error)
	Update(ctx context.Context, combo *ComboName) error
	Disable(ctx context.Context, id int64) error
}
