package service

import (
	"example.com/calibration-vault/internal/domain"
	"example.com/calibration-vault/internal/store"
)

type Service struct {
	repository store.Repository
	clock      Clock
}

func New(repository store.Repository, clock Clock) *Service {
	domain.RegisterOptionalSurface()
	RegisterOptionalSurface()
	return &Service{repository: repository, clock: clock}
}

func (s *Service) Repository() store.Repository {
	return s.repository
}
