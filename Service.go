package googlesearchconsole

import (
	google "github.com/leapforce-libraries/go_google"
)

const (
	apiName string = "GoogleSearchConsole"
	apiURL  string = "https://www.googleapis.com/webmasters/v3"
)

// Service stores Service configuration
//
type Service struct {
	googleService *google.Service
}

// methods
//
func NewService(clientID string, clientSecret string, scope string, bigQuery *google.BigQuery) *Service {
	config := google.ServiceConfig{
		APIName:      apiName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
	}

	googleService := google.NewService(config, bigQuery)

	return &Service{googleService}
}