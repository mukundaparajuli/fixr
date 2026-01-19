package repository

import "github.com/mukundaparajuli/fixr/internal/server"

type Repositories struct {
	Category *CategoryRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Category: NewCategoriesRepository(s),
	}
}
