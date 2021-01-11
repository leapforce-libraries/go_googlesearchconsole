package googlesearchconsole

import (
	"bytes"
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
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

func (service *Service) Query(body []byte, siteURL string) (*QueryResponse, *errortools.Error) {
	url := fmt.Sprintf("%s/sites/%s/searchAnalytics/query", apiURL, url.QueryEscape(siteURL))
	//fmt.Println(url)

	response := QueryResponse{}

	_, _, e := service.googleService.Post(url, bytes.NewBuffer(body), &response)
	if e != nil {
		return nil, e
	}

	return &response, nil
}
