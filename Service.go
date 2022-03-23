package googlesearchconsole

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	google "github.com/leapforce-libraries/go_google"
)

const (
	apiName    string = "GoogleSearchConsole"
	apiUrl     string = "https://www.googleapis.com/webmasters/v3"
	dateLayout string = "2006-01-02"
)

type Service google.Service

func NewServiceWithAccessToken(cfg *google.ServiceWithAccessTokenConfig) (*Service, *errortools.Error) {
	googleService, e := google.NewServiceWithAccessToken(cfg)
	if e != nil {
		return nil, e
	}
	service := Service(*googleService)
	return &service, nil
}

func NewServiceWithApiKey(cfg *google.ServiceWithApiKeyConfig) (*Service, *errortools.Error) {
	googleService, e := google.NewServiceWithApiKey(cfg)
	if e != nil {
		return nil, e
	}
	service := Service(*googleService)
	return &service, nil
}

func NewServiceWithOAuth2(cfg *google.ServiceWithOAuth2Config) (*Service, *errortools.Error) {
	googleService, e := google.NewServiceWithOAuth2(cfg)
	if e != nil {
		return nil, e
	}
	service := Service(*googleService)
	return &service, nil
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiUrl, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.googleService().ApiKey()
}

func (service *Service) ApiCallCount() int64 {
	return service.googleService().ApiCallCount()
}

func (service *Service) ApiReset() {
	service.googleService().ApiReset()
}

func (service *Service) googleService() *google.Service {
	googleService := google.Service(*service)
	return &googleService
}
