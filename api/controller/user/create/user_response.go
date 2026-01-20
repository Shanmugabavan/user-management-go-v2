package dtos

import "github.com/jackc/pgx/v5/pgtype"

type CreateUserResponse struct {
	UserID pgtype.UUID
	Email  string
	Status int32
}
