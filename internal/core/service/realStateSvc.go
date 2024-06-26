package service

import (
	"context"

	"github.com/natanchagas/gin-crud/internal/core/domain"
	"github.com/natanchagas/gin-crud/internal/core/ports"
)

type realStateService struct {
	repository ports.RealStateRepository
}

func NewRealStateService(r ports.RealStateRepository) *realStateService {
	return &realStateService{
		repository: r,
	}
}

func (s *realStateService) Create(ctx context.Context, realState domain.RealState) (domain.RealState, error) {
	id, err := s.repository.CreateRealState(ctx, realState)

	if err != nil {
		return domain.RealState{}, err
	}

	realState.Id = uint64(id)

	return realState, nil
}

func (s *realStateService) Get(ctx context.Context, id uint64) (domain.RealState, error) {
	realState, err := s.repository.GetRealState(ctx, id)
	if err != nil {
		return domain.RealState{}, err
	}

	return realState, nil
}

func (s *realStateService) Update(ctx context.Context, realState domain.RealState, id uint64) (domain.RealState, error) {
	_, err := s.repository.GetRealState(ctx, id)
	if err != nil {
		return domain.RealState{}, err
	}

	realState, err = s.repository.UpdateRealState(ctx, realState, id)

	if err != nil {
		return domain.RealState{}, err
	}

	return realState, nil
}

func (s *realStateService) Delete(ctx context.Context, id uint64) error {
	return s.repository.DeleteRealState(ctx, id)
}
