package repository

import "github.com/mukundaparajuli/fixr/internal/server"

type Repositories struct {
	Category *CategoryRepository
	Service  *ServiceRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Category: NewCategoryRepository(s),
		Service:  NewServiceRepository(s),
	}
}
