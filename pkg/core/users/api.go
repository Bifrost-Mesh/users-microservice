package users

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Bifrost-Mesh/users-microservice/proto/generated"
)

type UsersAPI struct {
	generated.UnimplementedUsersServiceServer

	usersService *UsersService
}

func NewUsersAPI(usersService *UsersService) *UsersAPI {
	//nolint:exhaustruct
	return &UsersAPI{
		usersService: usersService,
	}
}

func (*UsersAPI) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (u *UsersAPI) Signup(ctx context.Context,
	request *generated.SignupRequest,
) (*generated.SigninResponse, error) {
	output, err := u.usersService.Signup(ctx, &SignupArgs{
		Name:     request.GetName(),
		Email:    request.GetEmail(),
		Username: request.GetUsername(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	response := &generated.SigninResponse{
		Jwt: output.JWT,
	}
	return response, nil
}

func (u *UsersAPI) Signin(ctx context.Context,
	request *generated.SigninRequest,
) (*generated.SigninResponse, error) {
	//nolint:exhaustruct
	args := &SigninArgs{
		Password: request.GetPassword(),
	}
	switch request.GetIdentifier().(type) {
	case *generated.SigninRequest_Email:
		args.Email = aws.String(request.GetEmail())

	case *generated.SigninRequest_Username:
		args.Username = aws.String(request.GetUsername())
	}

	output, err := u.usersService.Signin(ctx, args)
	if err != nil {
		return nil, err
	}

	response := &generated.SigninResponse{
		Jwt: output.JWT,
	}
	return response, nil
}
