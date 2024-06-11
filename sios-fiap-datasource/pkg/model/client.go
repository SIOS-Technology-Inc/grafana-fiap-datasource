package model

import (
	"math/rand"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type FiapApiClient interface {
	CheckHealth() (*backend.CheckHealthResult, error)
	FetchWithDateRange(dataRange DataRangeType, fromTime *time.Time, toTime *time.Time, pointID string) (*data.Frame, error)
}

type FiapApiClientCreator func(settings *FiapDatasourceSettings) (FiapApiClient, error)

var (
	_ FiapApiClient        = (*MockClient)(nil)
	_ FiapApiClientCreator = CreateMockClient
)

type MockClient struct {
	ConnectionURL string
}

func CreateMockClient(settings *FiapDatasourceSettings) (FiapApiClient, error) {
	return &MockClient{ConnectionURL: settings.Url}, nil
}

func (*MockClient) CheckHealth() (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	if rand.Int()%2 == 0 {
		status = backend.HealthStatusError
		message = "randomized error"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (*MockClient) FetchWithDateRange(dataRange DataRangeType, fromTime *time.Time, toTime *time.Time, pointID string) (*data.Frame, error) {
	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/developers/plugin-tools/introduction/data-frames
	frame := data.NewFrame("response")

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{*fromTime, *toTime}),
		data.NewField("values", nil, []int64{10, 20}),
	)

	return frame, nil
}
