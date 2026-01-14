package create

import "github.com/jackc/pgx/v5/pgtype"

type UserResponse struct {
	UserID pgtype.UUID
	Email  string
	Status int32
}
