package integration

import (
	"context"
	"testing"
	"user-management/domain"
	"user-management/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRespository(t *testing.T) {
	_, connectionPool, err := SetupTestDatabase()
	if err != nil {
		return
	}

	userRepository := repository.NewUserRepository(connectionPool)

	newUser := domain.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "abc@gmail.com",
		Phone:     "1234567890",
		Age:       25,
		Status:    domain.UserStatusActive,
		UserId:    uuid.New(),
	}

	t.Run("CreateNewUser", func(t *testing.T) {
		created, err := userRepository.Create(context.Background(), &newUser)
		assert.NoError(t, err)
		assert.NotNil(t, created)
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		users, err := userRepository.GetAll(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
	})

	t.Run("GetUserById", func(t *testing.T) {
		user, err := userRepository.GetById(context.Background(), newUser.UserId)
		assert.NoError(t, err)
		assert.Equal(t, newUser, user)
	})

	t.Run("GetUserByIdForNonExistingUserId", func(t *testing.T) {
		user, err := userRepository.GetById(context.Background(), uuid.New())
		assert.Error(t, err)
		assert.Empty(t, user)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		deletedID, err := userRepository.Delete(context.Background(), newUser.UserId)
		assert.NoError(t, err)
		assert.Equal(t, newUser.UserId, deletedID)

		user, err := userRepository.GetById(context.Background(), newUser.UserId)
		assert.Error(t, err)
		assert.Empty(t, user)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		userRepository.Create(context.Background(), &newUser)

		updatedUserRequest := domain.User{
			FirstName: "UpdatedFirstName",
			LastName:  "UpdatedLastName",
			Email:     "updatedEmail@email.com",
			Phone:     "updatedPhone",
			Age:       26,
			Status:    domain.UserStatusInactive,
			UserId:    newUser.UserId,
		}
		updatedUserRow, _ := userRepository.Update(context.Background(), newUser.UserId, &updatedUserRequest)
		assert.NotEmpty(t, updatedUserRow)
		assert.Equal(t, updatedUserRequest.Email, updatedUserRow.Email)
		assert.Equal(t, updatedUserRequest.UserId, repository.ToUUIDFromPgUUID(updatedUserRow.UserID))
	})
}
