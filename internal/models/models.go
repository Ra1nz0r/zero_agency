package models

import (
	db "github.com/Ra1nz0r/zero_agency/db/sqlc"
)

type InputEditNews struct {
	Id         int64   `json:"Id" validate:"required"`
	Title      string  `json:"Title" validate:"omitempty,notblank,max=100"`
	Content    string  `json:"Content" validate:"omitempty,notblank"`
	Categories []int64 `json:"Categories" validate:"omitempty,min=1"`
}

type WriteResponse struct {
	Success bool         `json:"Success,omitempty"`
	News    []db.ListRow `json:"News,omitempty"`
}

type ParamId struct {
	Id int64
}
