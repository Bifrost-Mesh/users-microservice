package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"

	"github.com/Bifrost-Mesh/users-microservice/pkg/config"
	"github.com/Bifrost-Mesh/users-microservice/pkg/connectors"
	"github.com/Bifrost-Mesh/users-microservice/pkg/constants"
	"github.com/Bifrost-Mesh/users-microservice/pkg/core/token"
	"github.com/Bifrost-Mesh/users-microservice/pkg/core/users"
	"github.com/Bifrost-Mesh/users-microservice/pkg/grpc"
	"github.com/Bifrost-Mesh/users-microservice/pkg/healthcheck"
	"github.com/Bifrost-Mesh/users-microservice/pkg/logger"
	"github.com/Bifrost-Mesh/users-microservice/pkg/utils"
	"github.com/Bifrost-Mesh/users-microservice/proto/generated"
)

var configFilePath string

func main() {
	// Get CLI flag values.
	{
		flagSet := flag.NewFlagSet("", flag.ExitOnError)

		flagSet.StringVar(&configFilePath,
			constants.FLAG_CONFIG_FILE, constants.FLAG_CONFIG_FILE_DEFAULT,
			"Config file path",
		)

		flagSet.VisitAll(utils.CreateGetFlagOrEnvValueFn(""))

		cmdArgs := os.Args[1:]
		if err := flagSet.Parse(cmdArgs); err != nil {
			slog.Error("Failed parsing command line flags", logger.Error(err))
			os.Exit(1)
		}
	}

	// When the program receives any interruption / SIGKILL / SIGTERM signal, the cancel function is
	// automatically invoked. The cancel function is responsible for freeing all the resources
	// associated with the context.
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)

	// Construct validator with custom validators.
	validator := utils.NewValidator(ctx)

	// Get config.
	config := config.MustParseConfigFile(ctx, configFilePath, validator)

	if err := run(ctx, config, validator); err != nil {
		slog.ErrorContext(ctx, err.Error())

		cancel()

		// Give some time for remaining resources (if any) to be cleaned up.
		time.Sleep(constants.RESOURCES_CLEANUP_TIMEOUT)

		os.Exit(1)
	}
}

func run(ctx context.Context, config *config.Config, validator *validator.Validate) error {
	waitGroup, ctx := errgroup.WithContext(ctx)

	// Construct connectors.

	postgresConnector := connectors.NewPostgresConnector(ctx, config.Postgres.URL)
	defer postgresConnector.Shutdown()

	// Construct services.

	usersService := users.NewUsersService(
		validator,
		users.NewUsersPostgresRepository(ctx, postgresConnector),
		token.NewJWTService(config.JWTSigningKey),
	)

	// Run gRPC server.

	gRPCServer := grpc.NewGRPCServer(ctx, grpc.NewGRPCServerArgs{
		DevModeEnabled: config.DevMode,

		Healthcheckables: []healthcheck.Healthcheckable{
			postgresConnector,
		},

		ToGRPCErrorStatusCodeFn: getGRPCErrorStatusCode,
	})

	usersAPI := users.NewUsersAPI(usersService)
	generated.RegisterUsersServiceServer(gRPCServer, usersAPI)

	waitGroup.Go(func() error {
		return gRPCServer.Run(ctx, config.ServerPort)
	})

	/*
		The returned channel gets closed when either of this happens :

			(1) A program termination signal is received, because of which the parent context's done
				  channel gets closed.

			(2) Any of the go-routines registered under the wait-group, finishes running.
	*/
	<-ctx.Done()
	slog.DebugContext(ctx, "Gracefully shutting down program....")

	// Gracefully shutdown the gRPC server, ensuring that it finishes ongoing processing of requests.
	gRPCServer.GracefulShutdown()

	return waitGroup.Wait()
}
