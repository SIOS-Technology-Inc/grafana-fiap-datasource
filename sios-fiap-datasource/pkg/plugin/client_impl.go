package plugin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	fiapmodel "github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	dsmodel "github.com/sios/fiap/pkg/model"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap"
	"github.com/cockroachdb/errors"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ dsmodel.FiapApiClient        = (*ClientImpl)(nil)
	_ dsmodel.FiapApiClientCreator = CreateFiapApiClient
)

type ClientImpl struct {
	Client *fiap.FetchClient
}

func CreateFiapApiClient(settings *dsmodel.FiapDatasourceSettings) (dsmodel.FiapApiClient, error) {
	return &ClientImpl{Client: &fiap.FetchClient{ConnectionURL: settings.Url}}, nil
}

func (cli *ClientImpl) CheckHealth() (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("Start to check health", "ConnectionURL", cli.Client.ConnectionURL)
	resp, err := http.Head(cli.Client.ConnectionURL)
	if err != nil {
		log.DefaultLogger.Error("Failed to check health", "error", err, "response", resp)
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Failed to check health. Please see logs for details.",
		}, nil
	} else if resp.StatusCode > 299 {
		log.DefaultLogger.Error("URL returns bad status code", "statusCode", resp.StatusCode, "response", resp)
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("URL returns status code %d. Please see logs for details.", resp.StatusCode),
		}, nil
	}

	log.DefaultLogger.Info("Succeed to check health", "statusCode", resp.StatusCode)
	log.DefaultLogger.Debug("Succeed to check health (more info)", "response", resp)
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}

func (cli *ClientImpl) FetchWithDateRange(dataRange dsmodel.DataRangeType, fromTime *time.Time, toTime *time.Time, pointID string) (*data.Frame, error) {
	fetchErrors := make([]error, 0)

	var (
		pointSets map[string](fiapmodel.ProcessedPointSet)
		points    map[string]([]fiapmodel.Value)
		fiapErr   *fiapmodel.Error
		err       error
	)
	switch dataRange {
	case dsmodel.Period:
		pointSets, points, fiapErr, err = cli.Client.FetchDateRange(fromTime, toTime, pointID)
	case dsmodel.Latest:
		pointSets, points, fiapErr, err = cli.Client.FetchLatest(fromTime, toTime, pointID)
	case dsmodel.Oldest:
		pointSets, points, fiapErr, err = cli.Client.FetchOldest(fromTime, toTime, pointID)
	}
	if err != nil {
		fetchErrors = append(fetchErrors, err)
	}
	if fiapErr != nil {
		fetchErrors = append(fetchErrors, errors.Newf("fiap error: type %s, value %s", fiapErr.Type, fiapErr.Value))
	}
	if _, ok := pointSets[pointID]; ok {
		fetchErrors = append(fetchErrors, errors.Newf("point id '%s' provides point sets", pointID))
	}
	if _, ok := points[pointID]; !ok {
		fetchErrors = append(fetchErrors, errors.Newf("point id '%s' not provides point data", pointID))
		return nil, errors.Join(fetchErrors...)
	}

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/developers/plugin-tools/introduction/data-frames
	frame := data.NewFrame("response")

	// add fields.
	if times, values, convErr := pointsToFloatColumns(points[pointID]); convErr == nil {
		frame.Fields = append(frame.Fields,
			data.NewField("time", nil, times),
			data.NewField(pointID, nil, values),
		)
	} else {
		times, values := pointsToDefaultColumns(points[pointID])
		frame.Fields = append(frame.Fields,
			data.NewField("time", nil, times),
			data.NewField(pointID, nil, values),
		)
	}

	return frame, errors.Join(fetchErrors...)
}

func pointsToFloatColumns(pointArray []fiapmodel.Value) ([]time.Time, []float64, error) {
	if len(pointArray) == 0 {
		return nil, nil, errors.New("point is empty")
	}
	var (
		times  = make([]time.Time, len(pointArray))
		values = make([]float64, len(pointArray))
	)
	for i := range pointArray {
		times[i] = pointArray[i].Time
		if floatValue, err := strconv.ParseFloat(pointArray[i].Value, 64); err == nil {
			values[i] = floatValue
		} else {
			return nil, nil, errors.Newf("cannot parse to float: %s", pointArray[i].Value)
		}
	}
	return times, values, nil
}

func pointsToDefaultColumns(pointArray []fiapmodel.Value) ([]time.Time, []string) {
	var (
		times  = make([]time.Time, len(pointArray))
		values = make([]string, len(pointArray))
	)
	for i := range pointArray {
		times[i] = pointArray[i].Time
		values[i] = pointArray[i].Value
	}
	return times, values
}
