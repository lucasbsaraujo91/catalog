package entity

import "context"

// Repository define as operações para persistência de ComboName
type Repository interface {
	//Create(ctx context.Context, combo *ComboName) (int64, error)
	GetByID(ctx context.Context, id int64) (*ComboName, error)
	//GetAll(ctx context.Context) ([]*ComboName, error)
	//Update(ctx context.Context, combo *ComboName) error
	//Delete(ctx context.Context, id int64) error
}
