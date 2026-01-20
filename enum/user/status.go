package user

type UserStatus int

const (
	UserStatusDefault UserStatus = iota
	UserStatusActive
	UserStatusInactive
)
