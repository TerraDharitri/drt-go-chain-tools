package validatorStatsToPrometheus

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-tools/jsonToPrometheus/httpClientWrapper"
)

// HttpClientWrapper defines the behavior of wrapper over HttpClient
type HttpClientWrapper interface {
	GetValidatorStatistics(ctx context.Context) (map[string]*httpClientWrapper.ValidatorStatistics, error)
	IsInterfaceNil() bool
}
