package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-management/api/controller/user"
	"user-management/api/controller/user/create"
	"user-management/api/controller/user/update"
	"user-management/api/responses"
	"user-management/domain"
	"user-management/internal/db"
	"user-management/internal/validator"
	"user-management/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
}

func (m *mockRepo) Create(ctx context.Context, user *domain.User) (db.CreateUserRow, error) {
	return db.CreateUserRow{
		UserID: repository.ToPgUUID(user.UserId),
		Email:  user.Email,
		Status: int32(user.Status),
	}, nil
}

func (m *mockRepo) GetAll(c context.Context) ([]domain.User, error) {
	var users []domain.User

	users = append(users, domain.User{UserId: uuid.New()})
	users = append(users, domain.User{UserId: uuid.New()})

	return users, nil
}

func (m *mockRepo) GetById(c context.Context, id uuid.UUID) (domain.User, error) {
	return domain.User{UserId: id}, nil
}

func (m *mockRepo) Update(c context.Context, id uuid.UUID, user *domain.User) (db.UpdateUserRow, error) {
	return db.UpdateUserRow{UserID: repository.ToPgUUID(id)}, nil
}

func (m *mockRepo) Delete(c context.Context, id uuid.UUID) (uuid.UUID, error) {
	return uuid.New(), nil
}

func TestCreateUserWithValidData(t *testing.T) {
	mockUserController := user.UserController{
		&mockRepo{},
		nil,
	}

	createRequest := create.UserRequest{
		Email:     "s@gmail.com",
		Phone:     "+94776463619",
		Age:       2,
		Status:    1,
		FirstName: "ss",
		LastName:  "ss",
	}

	serializedObject, _ := json.Marshal(createRequest)
	requestBody := bytes.NewBuffer(serializedObject)
	request, _ := http.NewRequest(http.MethodPost, "", requestBody)

	request.Header.Set("Content-Type", "application/json")
	validator.Init()

	rr := httptest.NewRecorder()
	mockUserController.CreateUser(rr, request)

	var resp create.UserResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, createRequest.Email, resp.Email)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
}

func TestCreateUserWithInValidJsonData(t *testing.T) {
	mockUserController := user.UserController{
		&mockRepo{},
		nil,
	}

	createRequest := create.UserRequest{
		Email: "invalidEmail",
	}

	serializedObject, _ := json.Marshal(createRequest)
	requestBody := bytes.NewBuffer(serializedObject)
	request, _ := http.NewRequest(http.MethodPost, "", requestBody)

	request.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	validator.Init()
	mockUserController.CreateUser(rr, request)

	var resp responses.Response
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.Errors)
}

func TestGetAllUsers(t *testing.T) {
	mockUserController := user.UserController{
		&mockRepo{},
		nil,
	}

	request, _ := http.NewRequest(http.MethodPost, "", nil)

	request.Header.Set("Content-Type", "application/json")
	validator.Init()

	rr := httptest.NewRecorder()
	mockUserController.GetAllUsers(rr, request)

	var resp []domain.User
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestGetUserById(t *testing.T) {
	mockUserController := user.UserController{
		&mockRepo{},
		nil,
	}

	r := chi.NewRouter()
	r.Get("/users/{id}", mockUserController.GetUserById)

	id := uuid.New().String()
	request, _ := http.NewRequest(http.MethodGet, "/users/"+id, nil)

	validator.Init()

	var resp domain.User

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, request)

	err := json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestUpdateUser(t *testing.T) {
	mockUserController := user.UserController{
		&mockRepo{},
		nil,
	}

	updateRequest := update.UserRequest{
		Email: "sample@gmail.com",
	}

	serializedObject, _ := json.Marshal(updateRequest)
	requestBody := bytes.NewBuffer(serializedObject)

	r := chi.NewRouter()
	r.Put("/users/{id}", mockUserController.UpdateUser)

	id := uuid.New().String()
	request, _ := http.NewRequest(http.MethodPut, "/users/"+id, requestBody)

	validator.Init()

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusOK, rr.Code)

}

func TestUpdateUserWithInvalidEmail(t *testing.T) {
	mockUserController := user.UserController{
		&mockRepo{},
		nil,
	}

	updateRequest := update.UserRequest{
		Email: "sample.com",
	}

	serializedObject, _ := json.Marshal(updateRequest)
	requestBody := bytes.NewBuffer(serializedObject)

	r := chi.NewRouter()
	r.Put("/users/{id}", mockUserController.UpdateUser)

	id := uuid.New().String()
	request, _ := http.NewRequest(http.MethodPut, "/users/"+id, requestBody)

	validator.Init()

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
