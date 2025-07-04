package utils

import (
	"context"
	"log/slog"
	"os"

	"github.com/Bifrost-Mesh/users-microservice/pkg/assert"
)

// Returns the value of the given environment variable.
//
// Panics if the environment variable isn't set.
func MustGetEnv(name string) string {
	envValue, envFound := os.LookupEnv(name)
	assert.Assert(context.Background(), envFound, "Env not found", slog.String("env", name))

	return envValue
}
