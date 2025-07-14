package combonameusecase

import (
	"context"
	"errors"
	"log"

	comboname "catalog/internal/entity/comboname"
	"catalog/pkg/events"
)

type ComboNameServiceInterface interface {
	GetByID(ctx context.Context, id int64) (*comboname.ComboName, error)
	GetAll(ctx context.Context, page, limit int) ([]comboname.ComboName, int64, error)
	Update(ctx context.Context, combo *comboname.ComboName) error
	Create(ctx context.Context, combo *comboname.ComboName) (int64, error)
	Disable(ctx context.Context, id int64) error
}

type ComboNameService struct {
	repo       comboname.Repository
	event      events.EventInterface
	dispatcher events.EventDispatcherInterface
}

func NewComboNameService(r comboname.Repository, e events.EventInterface, d events.EventDispatcherInterface) *ComboNameService {
	return &ComboNameService{
		repo:       r,
		event:      e,
		dispatcher: d,
	}
}

func (uc *ComboNameService) GetByID(ctx context.Context, id int64) (*comboname.ComboName, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ComboNameService) GetAll(ctx context.Context, page, limit int) ([]comboname.ComboName, int64, error) {
	return uc.repo.GetAll(ctx, page, limit)
}

func (uc *ComboNameService) Update(ctx context.Context, combo *comboname.ComboName) error {
	return uc.repo.Update(ctx, combo)
}

func (uc *ComboNameService) Create(ctx context.Context, combo *comboname.ComboName) (int64, error) {
	if combo.Name == "" || combo.Nickname == "" {
		return 0, errors.New("name and nickname are required")
	}

	if combo.ComboNameUuid == "" {
		combo.GenerateUUID()
	}

	id, err := uc.repo.Create(ctx, combo)
	if err != nil {
		return 0, err
	}
	combo.ID = id

	if uc.dispatcher != nil {
		event := events.NewBaseEvent("ComboNameCreated", combo)

		log.Printf("Evento sendo despachado com payload: %+v\n", combo)

		if err := uc.dispatcher.Dispatch(event); err != nil {
			log.Printf("Erro ao despachar evento: %v\n", err)
		} else {
			log.Println("Evento despachado com sucesso")
		}
	}

	return id, nil
}

func (uc *ComboNameService) Disable(ctx context.Context, id int64) error {
	return uc.repo.Disable(ctx, id)
}
