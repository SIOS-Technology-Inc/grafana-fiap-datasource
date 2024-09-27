package model

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type FiapApiClient interface {
	CheckHealth() (*backend.CheckHealthResult, error)
	FetchWithDateRange(resp *backend.DataResponse, dataRange DataRangeType, fromTime *time.Time, toTime *time.Time, pointIDs []PointID, query *backend.DataQuery) error
}

type FiapApiClientCreator func(settings *FiapDatasourceSettings) (FiapApiClient, error)
