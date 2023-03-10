package service

import (
	"context"

	"github.com/CallumKerson/loggerrific"
)

type Service struct {
	logger loggerrific.Logger
}

func New(logger loggerrific.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) IsReady(ctx context.Context) bool {
	return true
}
