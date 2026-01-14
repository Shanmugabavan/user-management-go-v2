package route

import (
	"user-management/api/route/users"
	"user-management/bootstrap"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(env *bootstrap.Env, connectionPool *pgxpool.Pool, router *chi.Mux) {

	// Public APIs
	router.Group(func(r chi.Router) {
		users.UserRouter(env, connectionPool, router)
	})
}
