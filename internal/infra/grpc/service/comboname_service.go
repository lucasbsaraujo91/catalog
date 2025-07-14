package service

import (
	pb "catalog/internal/infra/grpc/pb"
	"catalog/internal/usecase/combonameusecase"
	context "context"
)

type ComboNameGrpcService struct {
	pb.UnimplementedComboNameServiceServer
	UseCase *combonameusecase.ComboNameService
}

func NewComboNameGrpcService(useCase *combonameusecase.ComboNameService) *ComboNameGrpcService {
	return &ComboNameGrpcService{
		UseCase: useCase,
	}
}

func (s *ComboNameGrpcService) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.ComboNameResponse, error) {
	combo, err := s.UseCase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ComboNameResponse{
		Id:            combo.ID,
		Name:          combo.Name,
		Nickname:      combo.Nickname,
		ComboNameUuid: combo.ComboNameUuid,
		IsAvailable:   combo.IsAvailable,
	}, nil
}

func (s *ComboNameGrpcService) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	combos, total, err := s.UseCase.GetAll(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, err
	}
	var pbCombos []*pb.ComboNameResponse
	for _, combo := range combos {
		pbCombos = append(pbCombos, &pb.ComboNameResponse{
			Id:            combo.ID,
			Name:          combo.Name,
			Nickname:      combo.Nickname,
			ComboNameUuid: combo.ComboNameUuid,
			IsAvailable:   combo.IsAvailable,
		})
	}
	return &pb.GetAllResponse{
		Combos: pbCombos,
		Total:  total,
	}, nil
}
