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

func (gsc *GoogleSearchConsole) Query(body []byte) (*QueryResponse, *errortools.Error) {
	e := gsc.Validate()
	if e != nil {
		return nil, e
	}

	url := fmt.Sprintf("%ssites/%s/searchAnalytics/query", gsc.baseURL, url.QueryEscape(gsc.SiteURL))
	//fmt.Println(url)

	response := QueryResponse{}

	_, _, e = gsc.post(url, bytes.NewBuffer(body), &response)
	if e != nil {
		return nil, e
	}

	return &response, nil
}
