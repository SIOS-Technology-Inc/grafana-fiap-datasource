package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sios/fiap/pkg/model"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
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
	log.DefaultLogger.Info("Start handle queries", "method", "QueryData", "user", req.PluginContext.User, "pluginId", req.PluginContext.PluginID)
	log.DefaultLogger.Debug("Start handle queries (more info)", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, &q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	log.DefaultLogger.Info("Finish handle queries", "method", "QueryData", "user", req.PluginContext.User, "pluginId", req.PluginContext.PluginID)
	log.DefaultLogger.Debug("Finish handle queries (more info)", "response", response)
	return response, nil
}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query *backend.DataQuery) backend.DataResponse {
	log.DefaultLogger.Info("Start handle query", "method", "query", "refID", query.RefID)
	log.DefaultLogger.Debug("Start handle query (more info)", "query", query)
	var response backend.DataResponse

	// Unmarshal the JSON into our query model.
	var qm model.FiapQuery
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		log.DefaultLogger.Error("Error parse json queries", "json", query.JSON, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	var serverTimezone *time.Location
	if tz, err := d.Settings.GetLocation(); err == nil {
		serverTimezone = tz
	} else {
		log.DefaultLogger.Error("Error parse server timezone in settings", "timezone", d.Settings.ServerTimezone, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("server timezone parse: %v", err.Error()))
	}
	var fromTime *time.Time
	if qm.StartTime.LinkDashboard {
		dt := query.TimeRange.From.In(serverTimezone)
		fromTime = &dt
	} else if dt, err := qm.StartTime.GetTime(d.Settings.ServerTimezone); err == nil {
		fromTime = dt
	} else {
		log.DefaultLogger.Error("Error parse start time in query", "time", qm.StartTime.RawTime, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("start time parse: %v", err.Error()))
	}
	var toTime *time.Time
	if qm.EndTime.LinkDashboard {
		dt := query.TimeRange.To.In(serverTimezone)
		toTime = &dt
	} else if dt, err := qm.EndTime.GetTime(d.Settings.ServerTimezone); err == nil {
		toTime = dt
	} else {
		log.DefaultLogger.Error("Error parse end time in query", "time", qm.EndTime.RawTime, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("end time parse: %v", err.Error()))
	}

	log.DefaultLogger.Info("Start fetch point data", "connectionURL", d.Settings.Url, "pointIDs", qm.PointIDs)
	log.DefaultLogger.Debug("Start fetch point data (more info)", "dataRange", qm.DataRange, "fromTime", fromTime, "toTime", toTime)
	err := d.Client.FetchWithDateRange(&response, qm.DataRange, fromTime, toTime, qm.PointIDs, query)
	if err != nil {
		log.DefaultLogger.Error("Error fetch point data", "json", query.JSON, "error", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("fiap fetch: %v", err.Error()))
	}
	log.DefaultLogger.Info("Finish fetch point data normally")

	log.DefaultLogger.Info("Finish handle query normally", "method", "query")
	log.DefaultLogger.Debug("Finish handle query normally (more info)", "response", response)
	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("Start health check", "method", "CheckHealth", "user", req.PluginContext.User, "pluginId", req.PluginContext.PluginID)
	log.DefaultLogger.Debug("Start health check (more info)", "request", req)

	result, err := d.Client.CheckHealth()
	if err != nil {
		log.DefaultLogger.Error("Unexpected error in health check", "error", err)
		return nil, err
	}

	log.DefaultLogger.Info("Finish health check normally", "method", "CheckHealth", "user", req.PluginContext.User, "pluginId", req.PluginContext.PluginID)
	log.DefaultLogger.Debug("Finish health check normally (more info)", "Status", result.Status, "Message", result.Message)
	return result, nil
}
