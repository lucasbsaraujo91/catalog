package combonameusecase

import (
	"context"

	comboname "catalog/internal/entity/comboname"
)

type ComboNameService struct {
	repo comboname.Repository
}

func NewComboNameService(r comboname.Repository) *ComboNameService {
	return &ComboNameService{repo: r}
}

func (s *ComboNameService) GetByID(ctx context.Context, id int64) (*comboname.ComboName, error) {
	return s.repo.GetByID(ctx, id)
}
