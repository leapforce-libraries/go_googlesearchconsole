package googlesearchconsole

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type QueryRequest struct {
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
	Dimensions []string `json:"dimensions"`
	RowLimit   int      `json:"rowLimit"`
}

type QueryResponse struct {
	Rows                    []QueryResponseRow `json:"rows"`
	ResponseAggregationType string             `json:"responseAggregationType"`
}

type QueryResponseRow struct {
	Keys        []string `json:"keys"`
	Impressions int      `json:"impressions"`
	Clicks      int      `json:"clicks"`
	CTR         float64  `json:"ctr"`
	Position    float64  `json:"position"`
}

func (service *Service) Query(queryRequest *QueryRequest, siteURL string) (*QueryResponse, *errortools.Error) {
	if queryRequest == nil {
		return nil, nil
	}

	response := QueryResponse{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("sites/%s/searchAnalytics/query", url.QueryEscape(siteURL))),
		BodyModel:     *queryRequest,
		ResponseModel: &response,
	}
	_, _, e := service.googleService.Post(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response, nil
}
