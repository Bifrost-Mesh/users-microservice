package main

import (
	"github.com/Bifrost-Mesh/users-microservice/pkg/constants"
	"github.com/Bifrost-Mesh/users-microservice/pkg/utils"
	"google.golang.org/grpc/codes"
)

// Returns suitable gRPC error status code, based on the given error.
func getGRPCErrorStatusCode(err error) codes.Code {
	if _, ok := err.(utils.APIError); ok {
		switch err {
		case constants.ErrInvalidEmail:
		case constants.ErrInvalidUsername:
		case constants.ErrInvalidJWT:
			return codes.InvalidArgument

		case constants.ErrDuplicateEmail:
		case constants.ErrDuplicateUsername:
			return codes.AlreadyExists

		case constants.ErrExpiredJWT:
			return codes.Unauthenticated

		case constants.ErrUserNotFound:
			return codes.NotFound

		default:
			return codes.Unknown
		}
	}

	return codes.Internal
}
