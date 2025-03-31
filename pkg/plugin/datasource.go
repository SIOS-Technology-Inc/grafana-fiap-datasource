package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sios/fiap/pkg/model"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

var createClient model.FiapApiClientCreator = CreateFiapApiClient

// NewDatasource creates a new datasource instance.
func NewDatasource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var ds *Datasource = &Datasource{}
	if err := json.Unmarshal(settings.JSONData, &(ds.Settings)); err != nil {
		return nil, err
	}
	if cli, err := createClient(&(ds.Settings)); err != nil {
		return nil, err
	} else {
		ds.Client = cli
	}
	return ds, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	Settings model.FiapDatasourceSettings
	Client   model.FiapApiClient
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctxLogger := backend.Logger.FromContext(ctx)
	ctxLogger.Info("Start QueryData in fiap datasource")
	ctxLogger.Debug("Start handle queries", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, &q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	ctxLogger.Debug("Finish handle queries", "response", response)
	return response, nil
}

func (d *Datasource) query(ctx context.Context, pCtx backend.PluginContext, query *backend.DataQuery) backend.DataResponse {
	ctxLogger := backend.Logger.FromContext(ctx)
	ctxLogger.Debug("Start handle query", "refID", query.RefID, "query", query)

	var response backend.DataResponse

	// Unmarshal the JSON into our query model.
	var qm model.FiapQuery
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		ctxLogger.Error("Error parse json queries", "json", query.JSON, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	var serverTimezone *time.Location
	if tz, err := d.Settings.GetLocation(); err == nil {
		serverTimezone = tz
	} else {
		ctxLogger.Error("Error parse server timezone in settings", "timezone", d.Settings.ServerTimezone, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("server timezone parse: %v", err.Error()))
	}
	var fromTime *time.Time
	if qm.StartTime.LinkDashboard {
		dt := query.TimeRange.From.In(serverTimezone)
		fromTime = &dt
	} else if dt, err := qm.StartTime.GetTime(d.Settings.ServerTimezone); err == nil {
		fromTime = dt
	} else {
		ctxLogger.Error("Error parse start time in query", "time", qm.StartTime.RawTime, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("start time parse: %v", err.Error()))
	}
	var toTime *time.Time
	if qm.EndTime.LinkDashboard {
		dt := query.TimeRange.To.In(serverTimezone)
		toTime = &dt
	} else if dt, err := qm.EndTime.GetTime(d.Settings.ServerTimezone); err == nil {
		toTime = dt
	} else {
		ctxLogger.Error("Error parse end time in query", "time", qm.EndTime.RawTime, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("end time parse: %v", err.Error()))
	}

	ctxLogger.Debug("Start fetch point data", "connectionURL", d.Settings.Url, "dataRange", qm.DataRange, "fromTime", fromTime, "toTime", toTime, "pointIDs", qm.PointIDs)
	err := d.Client.FetchWithDateRange(&response, qm.DataRange, fromTime, toTime, qm.PointIDs, query)
	if err != nil {
		ctxLogger.Error("Error fetch point data", "json", query.JSON, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("fiap fetch: %v", err.Error()))
	}

	ctxLogger.Debug("Finish handle query normally", "response", response)
	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	ctxLogger := backend.Logger.FromContext(ctx)
	ctxLogger.Info("Start CheckHealth in fiap datasource")
	ctxLogger.Debug("Start health check", "request", req)

	result, err := d.Client.CheckHealth()
	if err != nil {
		ctxLogger.Error("Unexpected error in health check", "error", err)
		return nil, err
	}

	ctxLogger.Debug("Finish health check normally", "Status", result.Status, "Message", result.Message)
	return result, nil
}
