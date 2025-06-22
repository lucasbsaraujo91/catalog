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

func (uc *ComboNameService) GetAll(ctx context.Context, page, limit int) ([]comboname.ComboName, int64, error) {
	return uc.repo.GetAll(ctx, page, limit)
}
