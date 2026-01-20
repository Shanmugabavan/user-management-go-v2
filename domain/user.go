package domain

import (
	"context"
	"user-management/enum/user"
	"user-management/internal/db"

	"github.com/google/uuid"
)

type User struct {
	UserId    uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Age       int
	Status    user.UserStatus
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (db.CreateUserRow, error)
	GetAll(c context.Context) ([]User, error)
	GetById(c context.Context, id uuid.UUID) (User, error)
	Update(c context.Context, id uuid.UUID, user *User) (db.UpdateUserRow, error)
	Delete(c context.Context, id uuid.UUID) (uuid.UUID, error)
}
