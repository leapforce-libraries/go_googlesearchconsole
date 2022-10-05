package googlesearchconsole

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type GetSitesResponse struct {
	Site []Site `json:"siteEntry"`
}

type Site struct {
	SiteUrl         string `json:"siteUrl"`
	PermissionLevel string `json:"permissionLevel"`
}

func (service *Service) GetSites() (*[]Site, *errortools.Error) {
	getSitesResponse := GetSitesResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url("sites"),
		ResponseModel: &getSitesResponse,
	}
	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &getSitesResponse.Site, nil
}
