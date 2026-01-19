package repository

import "github.com/mukundaparajuli/fixr/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
