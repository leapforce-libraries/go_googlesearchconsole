package googlesearchconsole

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	g_types "github.com/leapforce-libraries/go_googlesearchconsole/types"
	go_http "github.com/leapforce-libraries/go_http"
	go_types "github.com/leapforce-libraries/go_types"
)

type GetSitemapsResponse struct {
	Sitemap []Sitemap `json:"sitemap"`
}

type Sitemap struct {
	Path            string                 `json:"path"`
	LastSubmitted   g_types.DateTimeString `json:"lastSubmitted"`
	IsPending       bool                   `json:"isPending"`
	IsSitemapsIndex bool                   `json:"isSitemapsIndex"`
	Type            string                 `json:"type"`
	LastDownloaded  g_types.DateTimeString `json:"lastDownloaded"`
	Warnings        go_types.Int64String   `json:"warnings"`
	Errors          go_types.Int64String   `json:"errors"`
	Contents        []SitemapContent       `json:"contents"`
}

type SitemapContent struct {
	Type      string               `json:"type"`
	Submitted go_types.Int64String `json:"submitted"`
	Indexed   go_types.Int64String `json:"indexed"`
}

type GetSitemapsConfig struct {
	SiteURL      string
	SitemapIndex *string
}

func (service *Service) GetSitemaps(config *GetSitemapsConfig) (*[]Sitemap, *errortools.Error) {
	values := url.Values{}

	if config.SitemapIndex != nil {
		values.Set("sitemapIndex", *config.SitemapIndex)
	}

	getSitemapsResponse := GetSitemapsResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("sites/%s/sitemaps?%s", url.QueryEscape(config.SiteURL), values.Encode())),
		ResponseModel: &getSitemapsResponse,
	}
	_, _, e := service.googleService.HTTPRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &getSitemapsResponse.Sitemap, nil
}
