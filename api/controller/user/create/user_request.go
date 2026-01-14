package create

import (
	"user-management/domain"
)

type UserRequest struct {
	FirstName string            `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string            `json:"lastName" validate:"required,min=2,max=50"`
	Email     string            `json:"email" validate:"required,email"`
	Phone     string            `json:"phone" validate:"required,e164"`
	Age       int               `json:"age" validate:"required,gt=0"`
	Status    domain.UserStatus `json:"status" validate:"omitempty,oneof=0 1"`
}
