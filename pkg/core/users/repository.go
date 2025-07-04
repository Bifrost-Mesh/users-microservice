package users

import "context"

type (
	UsersRepository interface {
		Create(ctx context.Context, args *CreateUserArgs) (int32, error)

		FindByEmail(ctx context.Context, email string) (*FindUserByOperationOutput, error)
		FindByUsername(ctx context.Context, username string) (*FindUserByOperationOutput, error)
		UserIDExists(ctx context.Context, id int32) (bool, error)
	}

	CreateUserArgs struct {
		Name,
		Email,
		Username,
		HashedPassword string
	}

	FindUserByOperationOutput struct {
		ID             int32
		HashedPassword string
	}
)
