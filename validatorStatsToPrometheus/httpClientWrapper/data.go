package httpClientWrapper

import "github.com/TerraDharitri/drt-go-chain-core/data/validator"

// ValidatorStatistics defines the validator statistics api response
type ValidatorStatistics = validator.ValidatorStatistics

// ValidatorStatisticsApiResponse defines the response received when calling /validator/statistics endpoint
type ValidatorStatisticsApiResponse struct {
	Data struct {
		Statistics map[string]*ValidatorStatistics `json:"statistics"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}
