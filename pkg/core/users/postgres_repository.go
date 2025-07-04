package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Bifrost-Mesh/users-microservice/pkg/connectors"
	"github.com/Bifrost-Mesh/users-microservice/pkg/constants"
	"github.com/Bifrost-Mesh/users-microservice/pkg/utils"
	"github.com/Bifrost-Mesh/users-microservice/sql/generated"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UsersPostgresRepository struct {
	*connectors.PostgresConnector
	queries *generated.Queries
}

func NewUsersPostgresRepository(ctx context.Context,
	postgresConnector *connectors.PostgresConnector,
) UsersRepository {
	queries := generated.New(postgresConnector.GetConnection())

	return &UsersPostgresRepository{
		postgresConnector,
		queries,
	}
}

func (u *UsersPostgresRepository) Create(ctx context.Context,
	args *CreateUserArgs,
) (int32, error) {
	userID, err := u.queries.CreateUser(ctx, (*generated.CreateUserParams)(args))
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && (pgErr.Code == pgerrcode.UniqueViolation) {
			switch pgErr.ColumnName {
			case "email":
				return 0, constants.ErrDuplicateEmail

			case "username":
				return 0, constants.ErrDuplicateUsername
			}
		}

		return 0, utils.WrapError(err)
	}
	return userID, nil
}

func (u *UsersPostgresRepository) FindByEmail(ctx context.Context,
	email string,
) (*FindUserByOperationOutput, error) {
	userDetails, err := u.queries.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrUserNotFound
		}

		return nil, utils.WrapError(err)
	}
	return (*FindUserByOperationOutput)(userDetails), nil
}

func (u *UsersPostgresRepository) FindByUsername(ctx context.Context,
	username string,
) (*FindUserByOperationOutput, error) {
	userDetails, err := u.queries.FindUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrUserNotFound
		}

		return nil, utils.WrapError(err)
	}
	return (*FindUserByOperationOutput)(userDetails), nil
}

func (u *UsersPostgresRepository) UserIDExists(ctx context.Context, id int32) (bool, error) {
	_, err := u.queries.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, constants.ErrUserNotFound
		}

		return false, utils.WrapError(err)
	}
	return true, nil
}
