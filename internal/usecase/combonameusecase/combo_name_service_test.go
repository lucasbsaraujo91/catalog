package combonameusecase_test

import (
	"context"
	"errors"
	"testing"

	comboname "catalog/internal/entity/comboname"
	"catalog/internal/usecase/combonameusecase"
	"catalog/pkg/events"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository simula o comportamento do repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*comboname.ComboName, error) {
	args := m.Called(ctx, id)
	if obj := args.Get(0); obj != nil {
		return obj.(*comboname.ComboName), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetAll(ctx context.Context, page, limit int) ([]comboname.ComboName, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]comboname.ComboName), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) Create(ctx context.Context, combo *comboname.ComboName) (int64, error) {
	args := m.Called(ctx, combo)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, combo *comboname.ComboName) error {
	args := m.Called(ctx, combo)
	return args.Error(0)
}

func (m *MockRepository) Disable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockDispatcher simula o comportamento do dispatcher
type MockDispatcher struct {
	mock.Mock
}

func (m *MockDispatcher) Dispatch(event events.EventInterface) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockDispatcher) Clear() {
	m.Called()
}

func (m *MockDispatcher) Has(eventName string, handler events.EventHandlerInterface) bool {
	args := m.Called(eventName, handler)
	return args.Bool(0)
}

func (m *MockDispatcher) Register(eventName string, handler events.EventHandlerInterface) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *MockDispatcher) Remove(eventName string, handler events.EventHandlerInterface) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

// Teste do GetAll
func TestGetAll(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	expected := []comboname.ComboName{
		{ID: 1, Name: "Combo Teste", Nickname: "teste", IsAvailable: true, ComboNameUuid: "uuid-test"},
	}

	repo.On("GetAll", mock.Anything, 1, 10).Return(expected, int64(1), nil)

	service := combonameusecase.NewComboNameService(repo, nil, dispatcher)
	result, total, err := service.GetAll(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

// Teste do GetByID com sucesso
func TestGetByID_Success(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	expected := &comboname.ComboName{
		ID: 1, Name: "Combo X", Nickname: "x", IsAvailable: true, ComboNameUuid: "uuid-x",
	}

	repo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

	service := combonameusecase.NewComboNameService(repo, nil, dispatcher)
	result, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

// Teste do GetByID com erro
func TestGetByID_NotFound(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	repo.On("GetByID", mock.Anything, int64(99)).Return(nil, errors.New("not found"))

	service := combonameusecase.NewComboNameService(repo, nil, dispatcher)
	result, err := service.GetByID(context.Background(), 99)

	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
	repo.AssertExpectations(t)
}

// Teste do Create com sucesso
// Teste do Create com sucesso
func TestCreate_Success(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	combo := &comboname.ComboName{
		Name:        "Novo Combo",
		Nickname:    "novo",
		IsAvailable: true,
	}

	repo.On("Create", mock.Anything, combo).Return(int64(1), nil)

	// Espera o evento com nome "ComboNameCreated"
	dispatcher.On("Dispatch", mock.MatchedBy(func(e events.EventInterface) bool {
		return e.GetName() == "ComboNameCreated"
	})).Return(nil)

	service := combonameusecase.NewComboNameService(repo, nil, dispatcher)
	id, err := service.Create(context.Background(), combo)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	repo.AssertExpectations(t)
	dispatcher.AssertExpectations(t)
}

func TestCreate_InvalidInput(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	combo := &comboname.ComboName{
		Name:        "", // inválido
		Nickname:    "",
		IsAvailable: true,
	}

	service := combonameusecase.NewComboNameService(repo, nil, dispatcher)
	id, err := service.Create(context.Background(), combo)

	assert.Equal(t, int64(0), id)
	assert.EqualError(t, err, "name and nickname are required")
	repo.AssertNotCalled(t, "Create")
	dispatcher.AssertNotCalled(t, "Dispatch")
}

func TestCreate_RepoError(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	combo := &comboname.ComboName{
		Name:        "Erro Combo",
		Nickname:    "erro",
		IsAvailable: true,
	}

	repo.On("Create", mock.Anything, combo).Return(int64(0), errors.New("erro no banco"))

	service := combonameusecase.NewComboNameService(repo, nil, dispatcher)
	id, err := service.Create(context.Background(), combo)

	assert.Equal(t, int64(0), id)
	assert.EqualError(t, err, "erro no banco")
	dispatcher.AssertNotCalled(t, "Dispatch")
}

func TestCreate_DispatchError(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	combo := &comboname.ComboName{
		Name:        "Combo X",
		Nickname:    "x",
		IsAvailable: true,
	}

	repo.On("Create", mock.Anything, combo).Return(int64(1), nil)

	event := events.NewBaseEvent("ComboNameCreated", nil)
	dispatcher.On("Dispatch", mock.MatchedBy(func(e events.EventInterface) bool {
		return e.GetName() == "ComboNameCreated"
	})).Return(errors.New("erro no dispatcher"))

	service := combonameusecase.NewComboNameService(repo, event, dispatcher)
	id, err := service.Create(context.Background(), combo)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	repo.AssertExpectations(t)
	dispatcher.AssertExpectations(t)
}

// Teste do Update com sucesso
func TestUpdate_Success(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	combo := &comboname.ComboName{
		ID: 1, Name: "Atualizado", Nickname: "atual", IsAvailable: true,
	}

	repo.On("Update", mock.Anything, combo).Return(nil)

	service := combonameusecase.NewComboNameService(repo, nil, dispatcher)
	err := service.Update(context.Background(), combo)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// Teste do Disable com sucesso
func TestDisable_Success(t *testing.T) {
	repo := new(MockRepository)
	dispatcher := new(MockDispatcher)

	repo.On("Disable", mock.Anything, int64(1)).Return(nil)

	service := combonameusecase.NewComboNameService(repo, nil, dispatcher)
	err := service.Disable(context.Background(), 1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
