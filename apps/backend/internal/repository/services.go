package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mukundaparajuli/fixr/internal/errs"
	"github.com/mukundaparajuli/fixr/internal/model"
	"github.com/mukundaparajuli/fixr/internal/model/service"
	"github.com/mukundaparajuli/fixr/internal/server"
)

type ServiceRepository struct {
	server *server.Server
}

func NewServiceRepository(s *server.Server) *ServiceRepository {
	return &ServiceRepository{}
}

func (r *ServiceRepository) CreateService(ctx context.Context, userID string, payload *service.CreateServicePayload) (*service.Service, error) {
	stmt := `
		INSERT INTO
			services(
				user_id,
				name,
				description,
				status,
				rate,
				method,
				parent_service_id,
				category_id,
				metadata,
			)
		VALUES
			(
				@user_id,
				@name,
				@description,
				@status,
				@rate,
				@method,
				@parent_service_id,
				@category_id,
				@metadata
			)
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":           userID,
		"name":              payload.Name,
		"description":       *payload.Description,
		"rate":              payload.Rate,
		"method":            payload.Method,
		"parent_service_id": *payload.ParentServiceID,
		"category_id":       payload.CategoryID,
		"metadata":          *payload.Metadata,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute create service query for name=%s, user_id=%s : %w", payload.Name, userID, err)
	}

	serviceItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[service.Service])
	if err != nil {
		return nil, fmt.Errorf("failed to execute create service query for name=%s, user_id=%s :%w", payload.Name, userID, err)
	}
	return &serviceItem, nil
}

func (r *ServiceRepository) GetServiceByID(ctx context.Context, userID string, serviceID uuid.UUID) (*service.PopulatedService, error) {
	stmt := `
	SELECT
		s.*,
		CASE
			WHEN c.id IS NOT NULL THEN to_jsonb(camel (c))
			ELSE NULL
		END AS category,
		COALESCE(
			jsonb_agg(
				to_jsonb(camel (child))
				ORDER BY
					child.sort_order ASC,
					child.created_at ASC
			) FILTER (
				WHERE
					child.id IS NOT NULL
			),
			'[]'::JSONB
		) AS children
	FROM
		services s
		LEFT JOIN service_categories c ON c.id=s.category_id
		AND c.user_id=@user_id
		LEFT JOIN services child ON child.parent_service_id=s.id
		AND child.user_id=@user_id
	GROUP BY
		s.id,
		c.id
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      serviceID,
		"user_id": userID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute get service by id query for id=%s, user_id=%s", serviceID, userID)
	}

	serviceItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[service.PopulatedService])
	if err != nil {
		return nil, fmt.Errorf("failed to execute get service by id query for id=%s, user_id=%s", serviceID, userID)
	}

	return &serviceItem, nil
}

func (r *ServiceRepository) GetServices(ctx context.Context, userID string, query *service.GetServicesQuery) (*model.PaginatedResponse[service.PopulatedService], error) {
	stmt := `
	SELECT
		s.*,
		CASE
			WHEN c.id IS NOT NULL THEN to_jsonb(camel (c))
			ELSE NULL
		END AS category,
		COALESCE(
			jsonb_agg(
				to_jsonb(camel (child))
				ORDER BY
					child.sort_order ASC,
					child.created_at ASC
			) FILTER (
				WHERE
					child.id IS NOT NULL
			),
			'[]'::JSONB
		) AS children
	FROM
		services s
		LEFT JOIN service_categories c ON c.id=s.category_id
		AND c.user_id=@user_id
		LEFT JOIN services child ON child.parent_service_id=s.id
		AND child.user_id=@user_id
		`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	conditions := []string{"s.user_id == @user_id"}
	if query.Status != nil {
		conditions = append(conditions, "s.status = @status")
		args["status"] = *query.Status
	}

	if query.CategoryID != nil {
		conditions = append(conditions, "s.category_id = @category_id")
		args["category_id"] = *query.CategoryID
	}

	if query.ParentServiceID != nil {
		conditions = append(conditions, "s.parent_service_id = @parent_service_id")
		args["parent_todo_id"] = *query.ParentServiceID
	} else {
		// By default, only show root services (no parent)
		conditions = append(conditions, "s.parent_service_id IS NULL")
	}

	if query.Search != nil {
		conditions = append(conditions, "(s.title ILIKE @search OR s.description ILIKE @search)")
		args["search"] = "%" + *query.Search + "%"
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	countStmt := "SELECT COUNT(*) FROM services s"
	if len(conditions) > 0 {
		countStmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count for services user_id=%s: %w", userID, err)
	}

	stmt += " GROUP BY s.id, c.id"

	if query.Sort != nil {
		stmt += " ORDER BY s." + *query.Sort
		if query.Order != nil && *query.Order == "desc" {
			stmt += " DESC"
		} else {
			stmt += " ASC"
		}
	} else {
		stmt += " ORDER BY s.created_at DESC"
	}

	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get todos query for user_id=%s: %w", userID, err)
	}

	services, err := pgx.CollectRows(rows, pgx.RowToStructByName[service.PopulatedService])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[service.PopulatedService]{
				Data:       []service.PopulatedService{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:todos for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[service.PopulatedService]{
		Data:       services,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil

}

func (r *ServiceRepository) UpdateService(ctx context.Context, userID string, serviceID uuid.UUID, payload *service.UpdateServicePayload) (*service.Service, error) {
	stmt := `UPDATE services SET`
	args := pgx.NamedArgs{
		"id":      serviceID,
		"user_id": userID,
	}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name=@name")
		args["name"] = *payload.Name
	}

	if payload.Description != nil {
		setClauses = append(setClauses, "description=@description")
		args["description"] = *payload.Description
	}

	if payload.CategoryID != nil {
		setClauses = append(setClauses, "category_id=@category_id")
		args["category_id"] = *payload.CategoryID
	}

	if payload.Metadata != nil {
		setClauses = append(setClauses, "metadata=@metadata")
		args["metadata"] = *payload.Metadata
	}

	if len(setClauses) == 0 {
		return nil, errs.NewBadRequestError("no fields to update", false, nil, nil, nil)
	}

	stmt += strings.Join(setClauses, ",")
	stmt += "where id=@id AND user_id=@user_id"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	updatedService, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[service.Service])
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &updatedService, nil
}

func (r *ServiceRepository) DeleteService(ctx context.Context, userID string, serviceID uuid.UUID) error {
	stmt := `
		DELETE FROM services
		WHERE 
			id=@serviceID
			AND user_id=@user_id 	
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"id":         userID,
		"service_id": serviceID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute query : %w", err)
	}
	if result.RowsAffected() == 0 {
		code := "SERVICE_NOT_FOUND"
		return errs.NewNotFoundError("services not found", false, &code)
	}
	return nil
}
