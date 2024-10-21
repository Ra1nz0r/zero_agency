package models

import (
	db "github.com/Ra1nz0r/zero_agency/db/sqlc"
)

type InputEditNews struct {
	ID         int64   `json:"ID" validate:"required,min=1"`
	Title      string  `json:"Title" validate:"omitempty,notblank,max=100"`
	Content    string  `json:"Content" validate:"omitempty,notblank"`
	Categories []int64 `json:"Categories" validate:"omitempty,min=1"`
}

type LoginRequest struct {
	Username string `json:"Username" validate:"required,notblank,min=1,max=20"`
	Password string `json:"Password" validate:"required,notblank,min=1,max=20"`
}

type WriteResponse struct {
	Success bool         `json:"Success,omitempty"`
	News    []db.ListRow `json:"News,omitempty"`
}
