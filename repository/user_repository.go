package repository

import (
	"context"
	"user-management/domain"
	"user-management/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	connectionPool *pgxpool.Pool
	queries        *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &UserRepository{
		connectionPool: pool,
		queries:        db.New(pool),
	}
}

func (ur *UserRepository) Create(c context.Context, user *domain.User) (db.CreateUserRow, error) {
	createdd, err := ur.queries.CreateUser(c, db.CreateUserParams{
		UserID:    ToPgUUID(user.UserId),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       int32(user.Age),
		Status:    int32(user.Status),
	})

	created := db.CreateUserRow{
		UserID: createdd.UserID,
		Email:  createdd.Email,
		Status: createdd.Status,
	}

	return created, err
}

func (ur *UserRepository) GetAll(c context.Context) ([]domain.User, error) {
	dbUsers, err := ur.queries.GetAllUsers(c)

	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(dbUsers))

	for _, u := range dbUsers {
		users = append(users, domain.User{
			UserId:    ToUUIDFromPgUUID(u.UserID),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Phone:     u.Phone,
			Age:       int(u.Age),
			Status:    domain.UserStatus(u.Status),
		})
	}

	return users, nil
}

func (ur *UserRepository) GetById(c context.Context, id uuid.UUID) (domain.User, error) {
	dbUser, err := ur.queries.GetUser(c, ToPgUUID(id))

	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		UserId:    ToUUIDFromPgUUID(dbUser.UserID),
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Email:     dbUser.Email,
		Phone:     dbUser.Phone,
		Age:       int(dbUser.Age),
		Status:    domain.UserStatus(dbUser.Status),
	}

	return user, nil
}

func (ur *UserRepository) Update(c context.Context, id uuid.UUID, user *domain.User) (db.UpdateUserRow, error) {
	retrived, retError := ur.GetById(c, id)

	updateDbEntity(&retrived, user)

	if retError != nil {
		return db.UpdateUserRow{}, retError
	}

	createdd, err := ur.queries.UpdateUser(c, db.UpdateUserParams{
		UserID:    ToPgUUID(id),
		FirstName: retrived.FirstName,
		LastName:  retrived.LastName,
		Email:     retrived.Email,
		Phone:     retrived.Phone,
		Age:       int32(retrived.Age),
		Status:    int32(retrived.Status),
	})

	created := db.UpdateUserRow{
		UserID: createdd.UserID,
		Email:  createdd.Email,
		Status: createdd.Status,
	}

	return created, err
}

func (ur *UserRepository) Delete(c context.Context, id uuid.UUID) (uuid.UUID, error) {
	deletedUserId, err := ur.queries.DeleteUser(c, ToPgUUID(id))

	if err != nil {
		return uuid.New(), err
	}

	return ToUUIDFromPgUUID(deletedUserId), err
}

func ToPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func ToUUIDFromPgUUID(id pgtype.UUID) uuid.UUID {
	if !id.Valid {
		return uuid.Nil
	}
	return id.Bytes
}

func updateDbEntity(retrieved *domain.User, current *domain.User) {
	if current.FirstName != "" {
		retrieved.FirstName = current.FirstName
	}

	if current.LastName != "" {
		retrieved.LastName = current.LastName
	}

	if current.Email != "" {
		retrieved.Email = current.Email
	}

	if current.Phone != "" {
		retrieved.Phone = current.Phone
	}

	if current.Age != 0 {
		retrieved.Age = current.Age
	}

	if current.Status != 0 {
		retrieved.Status = current.Status
	}
}
