package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// --------------------------------------------------------------------------

type CreateServicePayload struct {
	Name            string     `json:"name" validate:"required,min=1,max=100"`
	Description     *string    `json:"description" validate:"omitempty,min=1,max=255"`
	Rate            *int       `json:"rate" validate:"min=1"`
	Method          *Method    `json:"method" validate:"oneof=hourly"`
	ParentServiceID *uuid.UUID `json:"parentServiceId" validate:"omitempty,uuid"`
	CategoryID      *uuid.UUID `json:"categoryId" validate:"omitempty,uuid"`
	Metadata        *Metadata  `json:"metadata"`
}

func (p *CreateServicePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// --------------------------------------------------------------------------

type UpdateServicePayload struct {
	ID              uuid.UUID  `param:"id" validate:"required,uuid"`
	Name            *string    `json:"name" validate:"required,min=1,max=100"`
	Description     *string    `json:"description" validate:"omitempty,min=1,max=255"`
	ParentServiceID *uuid.UUID `json:"parentServiceId" validate:"omitempty,uuid"`
	CategoryID      *uuid.UUID `json:"categoryId" validate:"omitempty,uuid"`
	Metadata        *Metadata  `json:"metadata"`
}

func (p *UpdateServicePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// --------------------------------------------------------------------------

type GetServicesQuery struct {
	Page            *int       `query:"page" validate:"omitempty,min=1"`
	Limit           *int       `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort            *string    `query:"sort" validate:"omitempty,oneof=created_at updated_at name"`
	Order           *string    `query:"order" validate:"omitempty,oneof = asc desc"`
	Search          *string    `query:"search" validate:"omitempty,min=1"`
	Status          *Status    `query:"status" validate:"omitempty, oneof = active inactive"`
	ParentServiceID *uuid.UUID `query:"parentServiceId" validate:"omitempty, uuid"`
	CategoryID      *uuid.UUID `query:"categoryId" validate:"omitempty,uuid"`
}

func (q *GetServicesQuery) Validate() error {
	validate := validator.New()

	if err := validate.Struct(q); err != nil {
		return err
	}

	// Set defaults for pagination
	if q.Page == nil {
		defaultPage := 1
		q.Page = &defaultPage
	}
	if q.Limit == nil {
		defaultLimit := 20
		q.Limit = &defaultLimit
	}
	if q.Sort == nil {
		defaultSort := "created_at"
		q.Sort = &defaultSort
	}
	if q.Order == nil {
		defaultOrder := "desc"
		q.Order = &defaultOrder
	}

	return nil
}

// --------------------------------------------------------------------------

type GetServiceByIDPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetServiceByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// --------------------------------------------------------------------------

type DeleteServiceByIDPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteServiceByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
