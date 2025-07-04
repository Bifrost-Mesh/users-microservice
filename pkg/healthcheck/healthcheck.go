package healthcheck

import "github.com/Bifrost-Mesh/users-microservice/pkg/utils"

type (
	HealthcheckFn = func() error

	Healthcheckable interface {
		Healthcheck() error
	}
)

// Checks health for each of the given health-checkable entities.
// Fails fast, i.e, when it encounters an unhealthy entity, it immediately returns error.
func Healthcheck(healthcheckables []Healthcheckable) error {
	for _, healthcheckable := range healthcheckables {
		if err := healthcheckable.Healthcheck(); err != nil {
			return utils.WrapErrorWithPrefix("Healthcheck failed", err)
		}
	}
	return nil
}
