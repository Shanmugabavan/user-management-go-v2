package users

import (
	"user-management/api/controller/user"
	"user-management/bootstrap"
	"user-management/repository"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRouter(env *bootstrap.Env, connectionPool *pgxpool.Pool, router *chi.Mux) {
	ur := repository.NewUserRepository(connectionPool)
	uc := &user.UserController{
		UserRepository: ur,
		Env:            env,
	}

	router.Post("/users", uc.CreateUser)
	router.Get("/users", uc.GetAllUsers)
	router.Get("/users/{id}", uc.GetUserById)
	router.Put("/users/{id}", uc.UpdateUser)
	router.Delete("/users/{id}", uc.DeleteUser)
}
