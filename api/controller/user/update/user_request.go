package update

type UserRequest struct {
	FirstName string `json:"firstName,omitempty" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"lastName,omitempty" validate:"omitempty,min=2,max=50"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,e164"`
	Age       int    `json:"age,omitempty" validate:"omitempty,gte=0,lte=150"`
	Status    int    `json:"status,omitempty"`
}
