package connectors

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Bifrost-Mesh/users-microservice/pkg/assert"
	"github.com/Bifrost-Mesh/users-microservice/pkg/logger"
	"github.com/Bifrost-Mesh/users-microservice/pkg/utils"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConnector struct {
	connection *sql.DB
}

func NewPostgresConnector(ctx context.Context, url string) *PostgresConnector {
	// TODO : use a query logger.
	connection, err := sql.Open("pgx", url)
	assert.AssertErrNil(ctx, err, "Failed connecting to Postgres")

	// Ping the database, verifying that a working connection has been established.
	err = connection.Ping()
	assert.AssertErrNil(ctx, err, "Failed connecting to Postgres")

	slog.DebugContext(ctx, "Connected to Postgres")

	return &PostgresConnector{connection}
}

func (p *PostgresConnector) GetConnection() *sql.DB {
	return p.connection
}

func (p *PostgresConnector) Healthcheck() error {
	if err := p.connection.Ping(); err != nil {
		return utils.WrapErrorWithPrefix("Failed pinging Postgres", err)
	}
	return nil
}

func (p *PostgresConnector) Shutdown() {
	if err := p.connection.Close(); err != nil {
		slog.Error("Failed closing Postgres connection", logger.Error(err))
		return
	}
	slog.Debug("Shut down Postgres client")
}
