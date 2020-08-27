package googlesearchconsole

import (
	"fmt"
	"net/url"
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

func (gsc *GoogleSearchConsole) Query(body []byte) (*QueryResponse, error) {
	err := gsc.Validate()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%ssites/%s/searchAnalytics/query", gsc.baseURL, url.QueryEscape(gsc.SiteURL))
	//fmt.Println(url)

	response := QueryResponse{}

	err = gsc.PostBytes(url, body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
