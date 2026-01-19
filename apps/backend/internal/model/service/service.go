package service

import (
	"github.com/google/uuid"
	"github.com/mukundaparajuli/fixr/internal/model"
	"github.com/mukundaparajuli/fixr/internal/model/category"
)

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type Method string

const (
	Hourly Method = "hourly"
)

type Metadata struct{}

type Service struct {
	model.Base
	UserID          string     `json:"userId" db:"user_id"`
	Name            string     `json:"name" db:"name"`
	Description     *string    `json:"description" db:"description"`
	Status          Status     `json:"status" db:"status"`
	Rate            *int       `json:"rate" db:"rate"`
	Method          Method     `json:"method" db:"method"`
	ParentServiceID *uuid.UUID `json:"parentServiceId" db:"parent_service_id"`
	CategoryID      *uuid.UUID `json:"categoryId" db:"category_id"`
	Metadata        *Metadata  `json:"metadata" db:"metadata"`
	SortOrder       int        `json:"sortOrder" db:"sort_order"`
}

type PopulatedService struct {
	Service
	Category    *category.Category
	Children    []Service
	Attachments []ServiceAttachment
}
