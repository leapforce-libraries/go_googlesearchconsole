package googlesearchconsole

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	google "github.com/leapforce-libraries/go_google"
	bigquery "github.com/leapforce-libraries/go_google/bigquery"
)

const (
	apiName    string = "GoogleSearchConsole"
	apiURL     string = "https://www.googleapis.com/webmasters/v3"
	dateFormat string = "2006-01-02"
)

// Service stores Service configuration
//
type Service struct {
	googleService *google.Service
}

type ServiceConfig struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

// methods
//
func NewService(config *ServiceConfig, bigQueryService *bigquery.Service) *Service {
	googleServiceConfig := google.ServiceConfig{
		APIName:      apiName,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scope:        config.Scope,
	}

	googleService := google.NewService(googleServiceConfig, bigQueryService)

	return &Service{googleService}
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) InitToken() *errortools.Error {
	return service.googleService.InitToken()
}
