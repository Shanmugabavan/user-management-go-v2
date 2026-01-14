package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pool "github.com/jackc/pgx/v5/pgxpool"
)

func GetConnectionPool(env *Env) *pool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		env.DBUser,
		env.DBPass,
		env.DBHost,
		env.DBPort,
		env.DBName,
		env.DBSSLMode,
	)

	var err error

	db, err := pool.New(ctx, connectionString)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}

	if err = db.Ping(ctx); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Successfully connected to database")
	migrateData(connectionString)
	fmt.Println("Successfully migration process done")

	return db
}

func CloseConnectionPool(db *pool.Pool) {
	db.Close()
}

func migrateData(connectionString string) {
	db, err := migrate.New("file://./migrations", connectionString)

	if err != nil {
		log.Fatal(err)
	}
	err = db.Up()
	if err != nil {
		fmt.Println(err)
	}

	db.Close()
}
