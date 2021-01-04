package googlesearchconsole

import (
	google "github.com/leapforce-libraries/go_google"
)

const (
	apiName string = "GoogleSearchConsole"
	apiURL  string = "https://www.googleapis.com/webmasters/v3"
)

// GoogleSearchConsole stores GoogleSearchConsole configuration
//
type GoogleSearchConsole struct {
	Client *google.GoogleClient
}

// methods
//
func NewGoogleSearchConsole(clientID string, clientSecret string, scope string, bigQuery *google.BigQuery) *GoogleSearchConsole {
	config := google.GoogleClientConfig{
		APIName:      apiName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
	}

	googleClient := google.NewGoogleClient(config, bigQuery)

	return &GoogleSearchConsole{googleClient}
}
