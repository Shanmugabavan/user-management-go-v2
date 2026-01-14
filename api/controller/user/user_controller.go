package user

import (
	"encoding/json"
	"net/http"
	"user-management/api/controller/user/create"
	"user-management/api/controller/user/update"
	"user-management/api/responses"
	"user-management/bootstrap"
	"user-management/domain"
	"user-management/internal/validator"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserController struct {
	domain.UserRepository
	Env *bootstrap.Env
}

// CreateUser godoc
// @Summary Create user
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body create.UserRequest true "User data"
// @Success 201 {object} domain.User
// @Failure 400 {object} responses.Response "Validation failed"
// @Failure 404 {object} responses.Response "User not found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /users [post]
func (u *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserRequest create.UserRequest

	_ = json.NewDecoder(r.Body).Decode(&createUserRequest)

	valError := validator.Validate.Struct(createUserRequest)
	if valError != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.Response{
			Message: "validation failed",
			Errors:  valError.Error(),
		})
		return
	}

	if createUserRequest.Status == domain.UserStatusDefault {
		createUserRequest.Status = domain.UserStatusActive
	}

	user := domain.User{
		UserId:    uuid.New(),
		FirstName: createUserRequest.FirstName,
		LastName:  createUserRequest.LastName,
		Email:     createUserRequest.Email,
		Phone:     createUserRequest.Phone,
		Age:       createUserRequest.Age,
		Status:    createUserRequest.Status,
	}

	createdUser, err2 := u.Create(r.Context(), &user)

	createUserResponse := create.UserResponse{
		UserID: createdUser.UserID,
		Email:  createdUser.Email,
		Status: createdUser.Status,
	}
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responses.Response{
			Message: "Internal Server Error",
			Errors:  err2.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(createUserResponse)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} domain.User "List of users"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /users [get]
func (u *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	userEntities, err2 := u.GetAll(r.Context())

	usersDtoResponse := make([]domain.User, 0, len(userEntities))

	for _, u := range userEntities {
		usersDtoResponse = append(usersDtoResponse, domain.User{
			UserId:    u.UserId,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Phone:     u.Phone,
			Age:       int(u.Age),
			Status:    u.Status,
		})
	}
	if err2 != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(responses.Response{
			Message: "Internal Server Error",
			Errors:  err2.Error(),
		})
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(usersDtoResponse)
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieve a single user by UUID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} domain.User "User found"
// @Failure 400 {object} responses.Response "Invalid user ID"
// @Failure 404 {object} responses.Response "User not found"
// @Failure 500 {object} responses.Response "Internal server error"
// @Router /users/{id} [get]
func (u *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(responses.Response{
			Message: "Internal Server Error",
			Errors:  err.Error(),
		})

		return
	}

	userEntity, err2 := u.GetById(r.Context(), userID)

	userResponse := domain.User{
		UserId:    userEntity.UserId,
		FirstName: userEntity.FirstName,
		LastName:  userEntity.LastName,
		Email:     userEntity.Email,
		Phone:     userEntity.Phone,
		Age:       userEntity.Age,
		Status:    userEntity.Status,
	}
	if err2 != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(responses.Response{
			Message: "user not found",
			Errors:  err2.Error(),
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(userResponse)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param user body update.UserRequest true "Update user payload"
// @Success 200 {object} create.UserResponse "User updated successfully"
// @Failure 400 {object} responses.Response "Invalid request / Validation failed"
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response "Internal server error"
// @Router /users/{id} [put]
func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateUserRequest update.UserRequest
	idParam := chi.URLParam(r, "id")

	userID, errId := uuid.Parse(idParam)

	if errId != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.Response{
			Message: "user not found",
			Errors:  errId.Error(),
		})
		return
	}

	err := json.NewDecoder(r.Body).Decode(&updateUserRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.Response{
			Message: "Json Conversion Issue",
			Errors:  err.Error(),
		})
		return
	}

	valError := validator.Validate.Struct(updateUserRequest)
	if valError != nil {
		http.Error(w, valError.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.Response{
			Message: "validation failed",
			Errors:  valError.Error(),
		})
		return
	}

	user := domain.User{
		FirstName: updateUserRequest.FirstName,
		LastName:  updateUserRequest.LastName,
		Email:     updateUserRequest.Email,
		Phone:     updateUserRequest.Phone,
		Age:       updateUserRequest.Age,
		Status:    domain.UserStatus(updateUserRequest.Status),
	}

	updatedUser, err2 := u.Update(r.Context(), userID, &user)

	createUserResponse := create.UserResponse{
		UserID: updatedUser.UserID,
		Email:  updatedUser.Email,
		Status: updatedUser.Status,
	}
	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responses.Response{
			Message: "user not found",
			Errors:  err2.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(createUserResponse)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 202 {string} string "User deleted successfully"
// @Failure 400 {object} responses.Response "Invalid user ID"
// @Failure 404 {object} responses.Response "User not found"
// @Failure 500 {object} responses.Response "Internal server error"
// @Router /users/{id} [delete]
func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, errId := uuid.Parse(idParam)

	if errId != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responses.Response{
			Message: "Invalid user id",
			Errors:  errId.Error(),
		})
		return
	}

	_, err2 := u.Delete(r.Context(), userID)

	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responses.Response{
			Message: "user not found",
			Errors:  err2.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
