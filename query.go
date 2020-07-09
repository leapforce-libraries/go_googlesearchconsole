package googlesearchconsole

import (
	"fmt"
	"net/url"
)

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
	url := fmt.Sprintf("%ssites/%s/searchAnalytics/query", gsc.BaseURL, url.QueryEscape(gsc.ClientURL))
	fmt.Println(url)

	response := QueryResponse{}

	err := gsc.PostBytes(url, body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
