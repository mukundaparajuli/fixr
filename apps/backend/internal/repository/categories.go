package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mukundaparajuli/fixr/internal/model/category"
	"github.com/mukundaparajuli/fixr/internal/server"
)

type CategoryRepository struct {
	server *server.Server
}

func NewCategoriesRepository(server *server.Server) *CategoryRepository {
	return &CategoryRepository{server: server}
}

func (r *CategoryRepository) CreateService(
	ctx context.Context,
	UserID string,
	payload *category.CreateCategoryPayload,
) (*CategoryRepository, error) {

	stmt := `
		INSERT INTO 
			service_categories(
				user_id,
				name,
				description,
				color
			)
		VALUES
			(
				@user_id,
				@name,
				@description,
				@color
			)	
		RETURNING
		*
	`
	rows, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"user_id": UserID,
		"name": payload.Name,
		"description": payload.Description,
		"color": payload.Color
	})

	if err != nil{
		return nil, fmt.Errorf("failed to execute create category query for user_id:%s, name:%s : %w", UserID, payload.Name, err)
	}

	categoryItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[category.Category])
	if err != nil{
		return nil, fmt.Errorf("failed to execute create category query for user_id=%s, name=%s :%w", UserID, payload.Name, err)
	}

	return &categoryItem, nil
}
