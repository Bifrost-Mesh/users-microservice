package grpc

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Bifrost-Mesh/users-microservice/pkg/logger"
	"github.com/Bifrost-Mesh/users-microservice/pkg/utils"
)

// Returns a gRPC request logger (which under the hood invokes slog).
// TODO : Filter out critical fields (like password).
func newGRPCRequestLogger(slogLogger *slog.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(ctx context.Context, logLevel logging.Level, message string, fields ...any) {
			slogLogger.Log(ctx, slog.Level(logLevel), message, fields...)
		},
	)
}

// Constructs and returns unary server interceptor for error handling.
// The error handler interceptor converts any error (returned from the usecases layer) to gRPC
// specific error.
func errorHandlerUnaryServerInterceptor(
	toGRPCErrorStatusCodeFn ToGRPCErrorStatusCodeFn,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		request any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		response, err := handler(ctx, request)
		return response, toGRPCError(ctx, err, toGRPCErrorStatusCodeFn)
	}
}

// Constructs and returns stream server interceptor for error handling.
// The error handler interceptor converts any error (returned from the usecases layer) to gRPC
// specific error.
func errorHandlerStreamServerInterceptor(
	toGRPCErrorStatusCodeFn ToGRPCErrorStatusCodeFn,
) grpc.StreamServerInterceptor {
	return func(
		server any,
		stream grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(server, stream)
		return toGRPCError(stream.Context(), err, toGRPCErrorStatusCodeFn)
	}
}

// Converts any error (returned from the usecases layer) to gRPC specific error.
// If the error is unexpected (not of type APIError), then that gets logged.
func toGRPCError(ctx context.Context,
	err error,
	toGRPCErrorStatusCodeFn ToGRPCErrorStatusCodeFn,
) error {
	if err == nil {
		return nil
	}

	switch err.(type) {
	case utils.APIError:
		return status.Error(toGRPCErrorStatusCodeFn(err), err.Error())

	default:
		// Log unexpected error.
		slog.ErrorContext(ctx, "Unexpected error occurred", logger.Error(err))

		return status.Error(codes.Internal, utils.ErrInternalServer.Error())
	}
}
