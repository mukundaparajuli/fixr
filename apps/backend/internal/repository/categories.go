package repository

import "github.com/mukundaparajuli/fixr/internal/server"

type CategoryRepository struct {
	server *server.Server
}

func NewCategoriesRepository(server *server.Server) *CategoryRepository {
	return &CategoryRepository{server: server}
}
