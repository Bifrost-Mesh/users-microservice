package users

import (
	"context"

	goValidator "github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/Bifrost-Mesh/users-microservice/pkg/core/token"
)

type UsersService struct {
	validator *goValidator.Validate

	usersRepository UsersRepository

	tokenService token.TokenService
}

func NewUsersService(
	validator *goValidator.Validate,
	usersRespository UsersRepository,
	tokenService token.TokenService,
) *UsersService {
	return &UsersService{
		validator,
		usersRespository,
		tokenService,
	}
}

type SignupArgs struct {
	Name     string `validate:"name"`
	Email    string `validate:"email"`
	Username string `validate:"username"`
	Password string `validate:"password"`
}

func (u *UsersService) Signup(ctx context.Context, args *SignupArgs) (*SigninOutput, error) {
	// Validate input.
	err := u.validator.StructCtx(ctx, args)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID, err := u.usersRepository.Create(ctx, &CreateUserArgs{
		Name:           args.Name,
		Email:          args.Email,
		Username:       args.Username,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	jwt, err := u.tokenService.Issue(userID)
	if err != nil {
		return nil, err
	}

	response := &SigninOutput{
		UserID: userID,
		JWT:    *jwt,
	}
	return response, nil
}

type (
	SigninArgs struct {
		Email    *string `validate:"omitempty,email"`
		Username *string `validate:"omitempty,username"`

		Password string `validate:"password"`
	}

	SigninOutput struct {
		UserID int32
		JWT    string
	}
)

func (u *UsersService) Signin(ctx context.Context, args *SigninArgs) (*SigninOutput, error) {
	// Validate input.
	err := u.validator.StructCtx(ctx, args)
	if err != nil {
		return nil, err
	}

	var userDetails *FindUserByOperationOutput
	switch {
	case args.Email != nil:
		userDetails, err = u.usersRepository.FindByEmail(ctx, *args.Email)

	case args.Username != nil:
		userDetails, err = u.usersRepository.FindByUsername(ctx, *args.Username)

	default:
		panic("unreachable")
	}
	if err != nil {
		return nil, err
	}

	jwt, err := u.tokenService.Issue(userDetails.ID)
	if err != nil {
		return nil, err
	}

	response := &SigninOutput{
		UserID: userDetails.ID,
		JWT:    *jwt,
	}
	return response, nil
}
