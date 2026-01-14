package integration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"user-management/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestDatabase() (testcontainers.Container, *pgxpool.Pool, error) {
	ctx := context.Background()
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:13.2",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "25621",
			"POSTGRES_DB":       "test_db",
			"POSTGRES_USER":     "postgres",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		).WithStartupTimeout(60 * time.Second),
	}

	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	if err != nil {
		return nil, nil, err
	}

	port, err := dbContainer.MappedPort(ctx, "5432")

	if err != nil {
		return nil, nil, err
	}

	host, err := dbContainer.Host(ctx)

	if err != nil {
		return nil, nil, err
	}

	connectionString := fmt.Sprintf("postgres://postgres:25621@%v:%v/test_db", host, port.Port())
	err = MigrateDb(connectionString)

	connectionPool, err := pgxpool.New(ctx, connectionString)

	if err != nil {
		return nil, nil, err
	}

	return dbContainer, connectionPool, err

}

func MigrateDb(connectionString string) (err error) {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, strings.Replace(connectionString, "postgres://", "pgx5://", 1))

	if err != nil {
		return err
	}

	err = m.Up()

	defer m.Close()

	if err != nil {
		return err
	}

	return nil
}
