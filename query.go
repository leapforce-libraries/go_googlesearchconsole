package googlesearchconsole

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Dimension string

const (
	DimensionDate             Dimension = "DATE"
	DimensionQuery            Dimension = "QUERY"
	DimensionPage             Dimension = "PAGE"
	DimensionCountry          Dimension = "COUNTRY"
	DimensionDevice           Dimension = "DEVICE"
	DimensionSearchAppearance Dimension = "SEARCH_APPEARANCE"
)

type SeachType string

const (
	SeachTypeNews  SeachType = "news"
	SeachTypeImage SeachType = "image"
	SeachTypeVideo SeachType = "video"
	SeachTypeWeb   SeachType = "web"
)

type AggregationType string

const (
	AggregationTypeAuto       AggregationType = "auto"
	AggregationTypeByPage     AggregationType = "byPage"
	AggregationTypeByProperty AggregationType = "byProperty"
)

type GroupType string

const (
	GroupTypeAnd GroupType = "AND"
)

type QueryRequest struct {
	StartDate             *time.Time
	EndDate               *time.Time
	Dimensions            *[]Dimension
	SearchType            *SeachType
	AggregationType       *AggregationType
	DimensionFilterGroups *[]struct {
		GroupType GroupType
		Filters   []struct {
			Dimension  Dimension
			Operator   string
			Expression string
		}
	}
	RowLimit *int
	StartRow *int
}

type queryRequest struct {
	StartDate             *string                 `json:"startDate,omitempty"`
	EndDate               *string                 `json:"endDate,omitempty"`
	Dimensions            *[]string               `json:"dimensions,omitempty"`
	SearchType            *string                 `json:"searchType,omitempty"`
	AggregationType       *string                 `json:"aggregationType,omitempty"`
	DimensionFilterGroups *[]dimensionFilterGroup `json:"dimensionFilterGroups,omitempty"`
	RowLimit              *int                    `json:"rowLimit,omitempty"`
	StartRow              *int                    `json:"startRow,omitempty"`
}

type dimensionFilterGroup struct {
	GroupType string            `json:"groupType"`
	Filters   []dimensionFilter `json:"filters"`
}

type dimensionFilter struct {
	Dimension  string `json:"dimension"`
	Operator   string `json:"operator"`
	Expression string `json:"expression"`
}

type QueryResponse struct {
	Rows                    []QueryResponseRow `json:"rows"`
	ResponseAggregationType string             `json:"responseAggregationType"`
}

type QueryResponseRow struct {
	Keys        []string `json:"keys"`
	Impressions float64  `json:"impressions"`
	Clicks      float64  `json:"clicks"`
	Ctr         float64  `json:"ctr"`
	Position    float64  `json:"position"`
}

func (service *Service) Query(_queryRequest *QueryRequest, siteURL string) (*QueryResponse, *errortools.Error) {
	if _queryRequest == nil {
		return nil, nil
	}

	qr := queryRequest{
		RowLimit: _queryRequest.RowLimit,
		StartRow: _queryRequest.StartRow,
	}
	if _queryRequest.StartDate != nil {
		startDate := _queryRequest.StartDate.Format(dateLayout)
		qr.StartDate = &startDate
	}
	if _queryRequest.EndDate != nil {
		endDate := _queryRequest.EndDate.Format(dateLayout)
		qr.EndDate = &endDate
	}
	if _queryRequest.Dimensions != nil {
		qr.Dimensions = &[]string{}
		for _, dimension := range *_queryRequest.Dimensions {
			*(qr.Dimensions) = append(*(qr.Dimensions), string(dimension))
		}
	}
	if _queryRequest.SearchType != nil {
		searchType := string(*_queryRequest.SearchType)
		qr.SearchType = &searchType
	}
	if _queryRequest.AggregationType != nil {
		aggregationType := string(*_queryRequest.AggregationType)
		qr.AggregationType = &aggregationType
	}
	if _queryRequest.DimensionFilterGroups != nil {
		qr.DimensionFilterGroups = &[]dimensionFilterGroup{}
		for _, dfg := range *_queryRequest.DimensionFilterGroups {
			_dimensionFilterGroup := dimensionFilterGroup{
				GroupType: string(dfg.GroupType),
			}
			for _, filter := range dfg.Filters {
				_filter := dimensionFilter{
					Dimension:  string(filter.Dimension),
					Operator:   string(filter.Operator),
					Expression: filter.Expression,
				}
				_dimensionFilterGroup.Filters = append(_dimensionFilterGroup.Filters, _filter)
			}
			*(qr.DimensionFilterGroups) = append(*(qr.DimensionFilterGroups), _dimensionFilterGroup)
		}
	}

	response := QueryResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url(fmt.Sprintf("sites/%s/searchAnalytics/query", url.QueryEscape(siteURL))),
		BodyModel:     qr,
		ResponseModel: &response,
	}
	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response, nil
}
