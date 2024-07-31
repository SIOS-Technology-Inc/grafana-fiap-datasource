package plugin

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/sios/fiap/pkg/model"
)

var (
	_ model.FiapApiClient        = (*MockClient)(nil)
	_ model.FiapApiClientCreator = createDefaultMockClient
)

type MockClient struct {
	actualArguments *fetchFuncArguments

	checkHealthFunc        func() (*backend.CheckHealthResult, error)
	fetchWithDateRangeFunc func(resp *backend.DataResponse, dataRange model.DataRangeType, fromTime *time.Time, toTime *time.Time, pointIDs []model.PointID, query *backend.DataQuery) error
}

type fetchFuncArguments struct {
	dataRange model.DataRangeType
	fromTime  *time.Time
	toTime    *time.Time
	pointIDs  []model.PointID
}

func createDefaultMockClient(settings *model.FiapDatasourceSettings) (model.FiapApiClient, error) {
	return &MockClient{
		checkHealthFunc: func() (*backend.CheckHealthResult, error) {
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusOk,
				Message: "Data source is working",
			}, nil
		},
		fetchWithDateRangeFunc: func(resp *backend.DataResponse, _ model.DataRangeType, fromTime *time.Time, toTime *time.Time, pointIDs []model.PointID, query *backend.DataQuery) error {
			for _, pointID := range pointIDs {
				// create data frame response.
				// For an overview on data frames and how grafana handles them:
				// https://grafana.com/developers/plugin-tools/introduction/data-frames
				frame := data.NewFrame(fmt.Sprintf("%s:%s", query.RefID, pointID.Value))

				// add fields.
				frame.Fields = append(frame.Fields,
					data.NewField("time", nil, []time.Time{*fromTime, *toTime}),
					data.NewField(pointID.Value, nil, []int64{10, 20}),
				)

				resp.Frames = append(resp.Frames, frame)
			}

			return nil
		},
	}, nil
}

func (cli *MockClient) CheckHealth() (*backend.CheckHealthResult, error) {
	return cli.checkHealthFunc()
}

func (cli *MockClient) FetchWithDateRange(resp *backend.DataResponse, dataRange model.DataRangeType, fromTime *time.Time, toTime *time.Time, pointIDs []model.PointID, query *backend.DataQuery) error {
	cli.actualArguments = &fetchFuncArguments{
		dataRange: dataRange,
		fromTime:  fromTime,
		toTime:    toTime,
		pointIDs:  pointIDs,
	}
	return cli.fetchWithDateRangeFunc(resp, dataRange, fromTime, toTime, pointIDs, query)
}

func TestNewDatasource(t *testing.T) {
	originalCreateClient := createClient
	t.Run("Normal", func(t *testing.T) {
		createClient = createDefaultMockClient

		expectedURL := "http://test.url:12345"
		expectedTimezone := "+09:00"
		inst, err := NewDatasource(context.Background(), backend.DataSourceInstanceSettings{
			JSONData: []byte(fmt.Sprintf(`{"url":"%s","server_timezone":"%s"}`, expectedURL, expectedTimezone)),
		})
		if err != nil {
			t.Error("failed to create new datasource")
		} else if inst == nil {
			t.Error("NewDatasource must return new datasource")
		} else {
			if ds, ok := inst.(*Datasource); !ok {
				t.Error("NewDatasource must return Datasource instance")
			} else if ds.Settings.Url != expectedURL {
				t.Errorf("NewDatasource must have url. expected:%s actual %s", expectedURL, ds.Settings.Url)
			}
		}
	})
	t.Run("Error", func(t *testing.T) {
		t.Run("InvalidSettings", func(t *testing.T) {
			createClient = createDefaultMockClient

			inst, err := NewDatasource(context.TODO(), backend.DataSourceInstanceSettings{
				JSONData: []byte(`{"url":"ht`),
			})
			if inst != nil {
				t.Error("NewDatasource must not return new datasource")
			} else if err == nil {
				t.Error("NewDatasource must return an error")
			}
		})
		t.Run("ClientCreation", func(t *testing.T) {
			createClient = func(_ *model.FiapDatasourceSettings) (model.FiapApiClient, error) {
				return nil, errors.New("test client creation error")
			}

			expectedURL := "http://test.url:12345"
			expectedTimezone := "+09:00"
			inst, err := NewDatasource(context.TODO(), backend.DataSourceInstanceSettings{
				JSONData: []byte(fmt.Sprintf(`{"url":"%s","server_timezone":"%s"}`, expectedURL, expectedTimezone)),
			})
			if inst != nil {
				t.Error("NewDatasource must not return new datasource")
			} else if err == nil {
				t.Error("NewDatasource must return an error")
			}

		})
	})
	createClient = originalCreateClient
}

func TestQueryData(t *testing.T) {
	ds := Datasource{Client: &MockClient{
		checkHealthFunc: func() (*backend.CheckHealthResult, error) {
			return nil, errors.New("not expected to call this function")
		},
		fetchWithDateRangeFunc: func(resp *backend.DataResponse, _ model.DataRangeType, _ *time.Time, _ *time.Time, pointIDs []model.PointID, query *backend.DataQuery) error {
			for _, pointID := range pointIDs {
				// create data frame response.
				// For an overview on data frames and how grafana handles them:
				// https://grafana.com/developers/plugin-tools/introduction/data-frames
				frame := data.NewFrame(fmt.Sprintf("%s:%s", query.RefID, pointID.Value))

				// add fields.
				frame.Fields = append(frame.Fields,
					data.NewField("time", nil, []time.Time{}),
					data.NewField(pointID.Value, nil, []int64{}),
				)

				resp.Frames = append(resp.Frames, frame)
			}

			return nil
		},
	}, Settings: model.FiapDatasourceSettings{
		Url:            "http://test.url:12345",
		ServerTimezone: "",
	}}
	t.Run("Normal", func(t *testing.T) {
		t.Run("EmptyTimeQuery", func(t *testing.T) {
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_a"}],"data_range":"period","start_time":{"time":"","link_dashboard":false},"end_time":{"time":"","link_dashboard":false}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if respA.Error != nil {
				t.Errorf("failed query of RefID '%s': %s", "A", respA.Error.Error())
			}

			if cli, ok := ds.Client.(*MockClient); !ok {
				t.Fatal("client is not mock")
			} else {
				if cli.actualArguments.dataRange != "period" {
					t.Errorf("expected datarange is %s but %s", "period", cli.actualArguments.dataRange)
				}
				if cli.actualArguments.fromTime != nil {
					t.Errorf("expected fromTime is nil but %s", cli.actualArguments.fromTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if cli.actualArguments.toTime != nil {
					t.Errorf("expected toTime is nil but %s", cli.actualArguments.toTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if len(cli.actualArguments.pointIDs) != 1 {
					t.Errorf("expected pointIDs' length is %d but %d", 1, len(cli.actualArguments.pointIDs))
				} else if pointID := cli.actualArguments.pointIDs[0].Value; pointID != "id_a" {
					t.Errorf("expected pointID[0] is %s but %s", "id_a", pointID)
				}
			}
		})
		t.Run("FixedTimeQuery", func(t *testing.T) {
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_b"}],"data_range":"latest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":false},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":false}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if respA.Error != nil {
				t.Errorf("failed query of RefID '%s': %s", "A", respA.Error.Error())
			}

			if cli, ok := ds.Client.(*MockClient); !ok {
				t.Fatal("client is not mock")
			} else {
				if cli.actualArguments.dataRange != "latest" {
					t.Errorf("expected datarange is %s but %s", "latest", cli.actualArguments.dataRange)
				}
				if cli.actualArguments.fromTime == nil {
					t.Errorf("expected fromTime is %s but nil", "2024-06-01 00:00:00 +00:00")
				} else if !cli.actualArguments.fromTime.Equal(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)) {
					t.Errorf("expected fromTime is %s but %s", "2024-06-01 00:00:00 +00:00", cli.actualArguments.fromTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if cli.actualArguments.toTime == nil {
					t.Errorf("expected toTime is %s but nil", "2024-06-30 23:59:59 +00:00")
				} else if !cli.actualArguments.toTime.Equal(time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC)) {
					t.Errorf("expected toTime is %s but %s", "2024-06-30 23:59:59 +00:00", cli.actualArguments.toTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if len(cli.actualArguments.pointIDs) != 1 {
					t.Errorf("expected pointIDs' length is %d but %d", 1, len(cli.actualArguments.pointIDs))
				} else if pointID := cli.actualArguments.pointIDs[0].Value; pointID != "id_b" {
					t.Errorf("expected pointID[0] is %s but %s", "id_b", pointID)
				}
			}
		})
		t.Run("LinkedTimeQuery", func(t *testing.T) {
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_c"}],"data_range":"oldest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":true},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":true}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if respA.Error != nil {
				t.Errorf("failed query of RefID '%s': %s", "A", respA.Error.Error())
			}

			if cli, ok := ds.Client.(*MockClient); !ok {
				t.Fatal("client is not mock")
			} else {
				if cli.actualArguments.dataRange != "oldest" {
					t.Errorf("expected datarange is %s but %s", "oldest", cli.actualArguments.dataRange)
				}
				if cli.actualArguments.fromTime == nil {
					t.Errorf("expected fromTime is %s but nil", "2024-03-01 00:00:00 +00:00")
				} else if !cli.actualArguments.fromTime.Equal(time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)) {
					t.Errorf("expected fromTime is %s but %s", "2024-03-01 00:00:00 +00:00", cli.actualArguments.fromTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if cli.actualArguments.toTime == nil {
					t.Errorf("expected toTime is %s but nil", "2024-03-31 23:59:59 +00:00")
				} else if !cli.actualArguments.toTime.Equal(time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC)) {
					t.Errorf("expected toTime is %s but %s", "2024-03-31 23:59:59 +00:00", cli.actualArguments.toTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if len(cli.actualArguments.pointIDs) != 1 {
					t.Errorf("expected pointIDs' length is %d but %d", 1, len(cli.actualArguments.pointIDs))
				} else if pointID := cli.actualArguments.pointIDs[0].Value; pointID != "id_c" {
					t.Errorf("expected pointID[0] is %s but %s", "id_c", pointID)
				}
			}
		})
		t.Run("MultipleQuery", func(t *testing.T) {
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_b"}],"data_range":"latest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":false},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":false}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
						{
							RefID: "B",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_c"}],"data_range":"oldest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":true},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":true}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 2 {
				t.Error("QueryData must return just two response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if respA.Error != nil {
				t.Errorf("failed query of RefID '%s': %s", "A", respA.Error.Error())
			}
			if respB, ok := resp.Responses["B"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "B")
			} else if respB.Error != nil {
				t.Errorf("failed query of RefID '%s': %s", "B", respB.Error.Error())
			}
		})
	})
	t.Run("Error", func(t *testing.T) {
		t.Run("InvalidQuery", func(t *testing.T) {
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_b"}],"data_r`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if !strings.Contains(respA.Error.Error(), "json unmarshal") {
				t.Errorf("expected error is %s but %s", "json unmarshal", respA.Error.Error())
			}
		})
		t.Run("InvalidTime", func(t *testing.T) {
			t.Run("Start", func(t *testing.T) {
				resp, err := ds.QueryData(
					context.Background(),
					&backend.QueryDataRequest{
						Queries: []backend.DataQuery{
							{
								RefID: "A",
								JSON:  []byte(`{"point_ids":[{"point_id":"id_b"}],"data_range":"latest","start_time":{"time":"2024-06-01","link_dashboard":false},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":false}}`),
								TimeRange: backend.TimeRange{
									From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
									To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
								},
							},
						},
					},
				)
				if err != nil {
					t.Fatal(err)
				}
				if len(resp.Responses) != 1 {
					t.Error("QueryData must return just one response")
				}
				if respA, ok := resp.Responses["A"]; !ok {
					t.Errorf("QueryData must return response of RefID '%s'", "A")
				} else if !strings.Contains(respA.Error.Error(), "start time parse") {
					t.Errorf("expected error is %s but %s", "start time parse", respA.Error.Error())
				}
			})
			t.Run("End", func(t *testing.T) {
				resp, err := ds.QueryData(
					context.Background(),
					&backend.QueryDataRequest{
						Queries: []backend.DataQuery{
							{
								RefID: "A",
								JSON:  []byte(`{"point_ids":[{"point_id":"id_b"}],"data_range":"latest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":false},"end_time":{"time":"2024/06/30T23:59:59","link_dashboard":false}}`),
								TimeRange: backend.TimeRange{
									From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
									To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
								},
							},
						},
					},
				)
				if err != nil {
					t.Fatal(err)
				}
				if len(resp.Responses) != 1 {
					t.Error("QueryData must return just one response")
				}
				if respA, ok := resp.Responses["A"]; !ok {
					t.Errorf("QueryData must return response of RefID '%s'", "A")
				} else if !strings.Contains(respA.Error.Error(), "end time parse") {
					t.Errorf("expected error is %s but %s", "end time parse", respA.Error.Error())
				}
			})
		})
		t.Run("FetchFailed", func(t *testing.T) {
			ds := Datasource{Client: &MockClient{
				checkHealthFunc: func() (*backend.CheckHealthResult, error) {
					return nil, errors.New("not expected to call this function")
				},
				fetchWithDateRangeFunc: func(_ *backend.DataResponse, _ model.DataRangeType, _ *time.Time, _ *time.Time, _ []model.PointID, _ *backend.DataQuery) error {
					return errors.New("test fetch error")
				},
			}, Settings: model.FiapDatasourceSettings{
				Url:            "http://test.url:12345",
				ServerTimezone: "+09:00",
			}}
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_b"}],"data_range":"latest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":false},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":false}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if !strings.Contains(respA.Error.Error(), "test fetch error") {
				t.Errorf("expected error is %s but %s", "test fetch error", respA.Error.Error())
			}
		})
	})
}

func TestQueryDataWithServerTz(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		ds := Datasource{Client: &MockClient{
			checkHealthFunc: func() (*backend.CheckHealthResult, error) {
				return nil, errors.New("not expected to call this function")
			},
			fetchWithDateRangeFunc: func(resp *backend.DataResponse, _ model.DataRangeType, _ *time.Time, _ *time.Time, pointIDs []model.PointID, query *backend.DataQuery) error {
				for _, pointID := range pointIDs {
					// create data frame response.
					// For an overview on data frames and how grafana handles them:
					// https://grafana.com/developers/plugin-tools/introduction/data-frames
					frame := data.NewFrame(fmt.Sprintf("%s:%s", query.RefID, pointID.Value))

					// add fields.
					frame.Fields = append(frame.Fields,
						data.NewField("time", nil, []time.Time{}),
						data.NewField(pointID.Value, nil, []int64{}),
					)

					resp.Frames = append(resp.Frames, frame)
				}

				return nil
			},
		}, Settings: model.FiapDatasourceSettings{
			Url:            "http://test.url:12345",
			ServerTimezone: "+09:00",
		}}
		expectedTimezone := time.FixedZone("+09:00", 9*60*60)
		t.Run("EmptyTimeQuery", func(t *testing.T) {
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_a"}],"data_range":"period","start_time":{"time":"","link_dashboard":false},"end_time":{"time":"","link_dashboard":false}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if respA.Error != nil {
				t.Errorf("failed query of RefID '%s': %s", "A", respA.Error.Error())
			}

			if cli, ok := ds.Client.(*MockClient); !ok {
				t.Fatal("client is not mock")
			} else {
				if cli.actualArguments.dataRange != "period" {
					t.Errorf("expected datarange is %s but %s", "period", cli.actualArguments.dataRange)
				}
				if cli.actualArguments.fromTime != nil {
					t.Errorf("expected fromTime is nil but %s", cli.actualArguments.fromTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if cli.actualArguments.toTime != nil {
					t.Errorf("expected toTime is nil but %s", cli.actualArguments.toTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if len(cli.actualArguments.pointIDs) != 1 {
					t.Errorf("expected pointIDs' length is %d but %d", 1, len(cli.actualArguments.pointIDs))
				} else if pointID := cli.actualArguments.pointIDs[0].Value; pointID != "id_a" {
					t.Errorf("expected pointID[0] is %s but %s", "id_a", pointID)
				}
			}
		})
		t.Run("FixedTimeQuery", func(t *testing.T) {
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_b"}],"data_range":"latest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":false},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":false}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if respA.Error != nil {
				t.Errorf("failed query of RefID '%s': %s", "A", respA.Error.Error())
			}

			if cli, ok := ds.Client.(*MockClient); !ok {
				t.Fatal("client is not mock")
			} else {
				if cli.actualArguments.dataRange != "latest" {
					t.Errorf("expected datarange is %s but %s", "latest", cli.actualArguments.dataRange)
				}
				if cli.actualArguments.fromTime == nil {
					t.Errorf("expected fromTime is %s but nil", "2024-06-01 00:00:00 +09:00")
				} else if !cli.actualArguments.fromTime.Equal(time.Date(2024, 6, 1, 0, 0, 0, 0, expectedTimezone)) {
					t.Errorf("expected fromTime is %s but %s", "2024-06-01 00:00:00 +09:00", cli.actualArguments.fromTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if cli.actualArguments.toTime == nil {
					t.Errorf("expected toTime is %s but nil", "2024-06-30 23:59:59 +09:00")
				} else if !cli.actualArguments.toTime.Equal(time.Date(2024, 6, 30, 23, 59, 59, 0, expectedTimezone)) {
					t.Errorf("expected toTime is %s but %s", "2024-06-30 23:59:59 +09:00", cli.actualArguments.toTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if len(cli.actualArguments.pointIDs) != 1 {
					t.Errorf("expected pointIDs' length is %d but %d", 1, len(cli.actualArguments.pointIDs))
				} else if pointID := cli.actualArguments.pointIDs[0].Value; pointID != "id_b" {
					t.Errorf("expected pointID[0] is %s but %s", "id_b", pointID)
				}
			}
		})
		t.Run("LinkedTimeQuery", func(t *testing.T) {
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_c"}],"data_range":"oldest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":true},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":true}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if respA.Error != nil {
				t.Errorf("failed query of RefID '%s': %s", "A", respA.Error.Error())
			}

			if cli, ok := ds.Client.(*MockClient); !ok {
				t.Fatal("client is not mock")
			} else {
				if cli.actualArguments.dataRange != "oldest" {
					t.Errorf("expected datarange is %s but %s", "oldest", cli.actualArguments.dataRange)
				}
				if cli.actualArguments.fromTime == nil {
					t.Errorf("expected fromTime is %s but nil", "2024-03-01 09:00:00 +09:00")
				} else if !cli.actualArguments.fromTime.Equal(time.Date(2024, 3, 1, 9, 0, 0, 0, expectedTimezone)) {
					t.Errorf("expected fromTime is %s but %s", "2024-03-01 09:00:00 +09:00", cli.actualArguments.fromTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if cli.actualArguments.toTime == nil {
					t.Errorf("expected toTime is %s but nil", "2024-04-01 8:59:59 +09:00")
				} else if !cli.actualArguments.toTime.Equal(time.Date(2024, 4, 1, 8, 59, 59, 0, expectedTimezone)) {
					t.Errorf("expected toTime is %s but %s", "2024-04-01 8:59:59 +09:00", cli.actualArguments.toTime.Format("2006-01-02 15:04:05 -07:00"))
				}
				if len(cli.actualArguments.pointIDs) != 1 {
					t.Errorf("expected pointIDs' length is %d but %d", 1, len(cli.actualArguments.pointIDs))
				} else if pointID := cli.actualArguments.pointIDs[0].Value; pointID != "id_c" {
					t.Errorf("expected pointID[0] is %s but %s", "id_c", pointID)
				}
			}
		})
	})
	t.Run("Error", func(t *testing.T) {
		t.Run("InvalidTimezone", func(t *testing.T) {
			ds := Datasource{Client: &MockClient{
				checkHealthFunc: func() (*backend.CheckHealthResult, error) {
					return nil, errors.New("not expected to call this function")
				},
				fetchWithDateRangeFunc: func(resp *backend.DataResponse, _ model.DataRangeType, _ *time.Time, _ *time.Time, pointIDs []model.PointID, query *backend.DataQuery) error {
					for _, pointID := range pointIDs {
						// create data frame response.
						// For an overview on data frames and how grafana handles them:
						// https://grafana.com/developers/plugin-tools/introduction/data-frames
						frame := data.NewFrame(fmt.Sprintf("%s:%s", query.RefID, pointID.Value))

						// add fields.
						frame.Fields = append(frame.Fields,
							data.NewField("time", nil, []time.Time{}),
							data.NewField(pointID.Value, nil, []int64{}),
						)

						resp.Frames = append(resp.Frames, frame)
					}

					return nil
				},
			}, Settings: model.FiapDatasourceSettings{
				Url:            "http://test.url:12345",
				ServerTimezone: "invalid",
			}}
			resp, err := ds.QueryData(
				context.Background(),
				&backend.QueryDataRequest{
					Queries: []backend.DataQuery{
						{
							RefID: "A",
							JSON:  []byte(`{"point_ids":[{"point_id":"id_b"}],"data_range":"latest","start_time":{"time":"2024-06-01 00:00:00","link_dashboard":false},"end_time":{"time":"2024-06-30 23:59:59","link_dashboard":false}}`),
							TimeRange: backend.TimeRange{
								From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
								To:   time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
							},
						},
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Responses) != 1 {
				t.Error("QueryData must return just one response")
			}
			if respA, ok := resp.Responses["A"]; !ok {
				t.Errorf("QueryData must return response of RefID '%s'", "A")
			} else if !strings.Contains(respA.Error.Error(), "server timezone parse") {
				t.Errorf("expected error is %s but %s", "server timezone parse", respA.Error.Error())
			}
		})
	})
}

func TestCheckHealth(t *testing.T) {
	t.Run("StatusOk", func(t *testing.T) {
		ds := Datasource{Client: &MockClient{
			checkHealthFunc: func() (*backend.CheckHealthResult, error) {
				return &backend.CheckHealthResult{
					Status:  backend.HealthStatusOk,
					Message: "Data source is working",
				}, nil
			},
			fetchWithDateRangeFunc: func(_ *backend.DataResponse, _ model.DataRangeType, _ *time.Time, _ *time.Time, _ []model.PointID, _ *backend.DataQuery) error {
				return errors.New("not expected to call this function")
			},
		}}

		res, err := ds.Client.CheckHealth()
		if err != nil {
			t.Error(err)
		}
		if res.Status != backend.HealthStatusOk {
			t.Errorf("expected status is %s but %s", backend.HealthStatusOk, res.Status)
		}
	})
	t.Run("StatusError", func(t *testing.T) {
		ds := Datasource{Client: &MockClient{
			checkHealthFunc: func() (*backend.CheckHealthResult, error) {
				return &backend.CheckHealthResult{
					Status:  backend.HealthStatusError,
					Message: "Data source is not working",
				}, nil
			},
			fetchWithDateRangeFunc: func(_ *backend.DataResponse, _ model.DataRangeType, _ *time.Time, _ *time.Time, _ []model.PointID, _ *backend.DataQuery) error {
				return errors.New("not expected to call this function")
			},
		}}

		res, err := ds.Client.CheckHealth()
		if err != nil {
			t.Error(err)
		}
		if res.Status != backend.HealthStatusError {
			t.Errorf("expected status is %s but %s", backend.HealthStatusError, res.Status)
		}
	})
	t.Run("UnexpectedError", func(t *testing.T) {
		ds := Datasource{Client: &MockClient{
			checkHealthFunc: func() (*backend.CheckHealthResult, error) {
				return nil, errors.New("test unexpected error")
			},
			fetchWithDateRangeFunc: func(_ *backend.DataResponse, _ model.DataRangeType, _ *time.Time, _ *time.Time, _ []model.PointID, _ *backend.DataQuery) error {
				return errors.New("not expected to call this function")
			},
		}}

		_, err := ds.Client.CheckHealth()
		if err == nil {
			t.Error("CheckHealth must return an error")
		} else if !strings.Contains(err.Error(), "test unexpected error") {
			t.Errorf("expected error is %s but %s", "test unexpected error", err.Error())
		}
	})
}
