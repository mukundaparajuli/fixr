package category

import "github.com/mukundaparajuli/fixr/internal/model"

type Category struct {
	model.Base
	UserID      string  `json:"userId" db:"user_id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	Color       string  `json:"color" db:"color"`
}
