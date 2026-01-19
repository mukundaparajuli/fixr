package service

import (
	"github.com/mukundaparajuli/fixr/internal/lib/job"
	"github.com/mukundaparajuli/fixr/internal/repository"
	"github.com/mukundaparajuli/fixr/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}
