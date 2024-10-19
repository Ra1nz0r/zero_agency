package models

import db "github.com/Ra1nz0r/zero_agency/db/sqlc"

type InputEditNews struct {
	Title      string  `json:"Title,omitempty"`
	Content    string  `json:"Content,omitempty"`
	Categories []int64 `json:"Categories,omitempty"`
}

type WriteResponse struct {
	Success bool         `json:"Success,omitempty"`
	News    []db.ListRow `json:"News,omitempty"`
}
