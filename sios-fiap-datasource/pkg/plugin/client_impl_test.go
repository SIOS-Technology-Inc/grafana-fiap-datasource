package plugin

import (
	"fmt"
	"strings"
	"testing"
	"time"

	fiapmodel "github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	dsmodel "github.com/sios/fiap/pkg/model"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap"
	"github.com/cockroachdb/errors"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var _ fiap.Fetcher = (*mockFetchClient)(nil)

type fetchClientArguments struct {
	fromDate  *time.Time
	untilDate *time.Time
	ids       []string
}

type fetchClientResults struct {
	pointSets map[string](fiapmodel.ProcessedPointSet)
	points    map[string]([]fiapmodel.Value)
	fiapErr   *fiapmodel.Error
}

type mockFetchClient struct {
	failLatest, failOldest, failDateRange bool

	actualArguments *fetchClientArguments
	results         *fetchClientResults
}

func (*mockFetchClient) Fetch(keys []fiapmodel.UserInputKey, option *fiapmodel.FetchOption) (pointSets map[string]fiapmodel.ProcessedPointSet, points map[string][]fiapmodel.Value, fiapErr *fiapmodel.Error, err error) {
	return nil, nil, nil, errors.New("unimplemented")
}

func (*mockFetchClient) FetchByIdsWithKey(key fiapmodel.UserInputKeyNoID, ids ...string) (pointSets map[string]fiapmodel.ProcessedPointSet, points map[string][]fiapmodel.Value, fiapErr *fiapmodel.Error, err error) {
	return nil, nil, nil, errors.New("unimplemented")
}

func (*mockFetchClient) FetchOnce(keys []fiapmodel.UserInputKey, option *fiapmodel.FetchOnceOption) (pointSets map[string]fiapmodel.ProcessedPointSet, points map[string][]fiapmodel.Value, cursor string, fiapErr *fiapmodel.Error, err error) {
	return nil, nil, "", nil, errors.New("unimplemented")
}

func (f *mockFetchClient) FetchDateRange(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string]fiapmodel.ProcessedPointSet, points map[string][]fiapmodel.Value, fiapErr *fiapmodel.Error, err error) {
	if f.failDateRange {
		return nil, nil, nil, errors.New("test FetchLatest error")
	}
	f.actualArguments = &fetchClientArguments{
		fromDate:  fromDate,
		untilDate: untilDate,
		ids:       ids,
	}
	return f.results.pointSets, f.results.points, f.results.fiapErr, nil
}

func (f *mockFetchClient) FetchLatest(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string]fiapmodel.ProcessedPointSet, points map[string][]fiapmodel.Value, fiapErr *fiapmodel.Error, err error) {
	if f.failLatest {
		return nil, nil, nil, errors.New("test FetchLatest error")
	}
	f.actualArguments = &fetchClientArguments{
		fromDate:  fromDate,
		untilDate: untilDate,
		ids:       ids,
	}
	return f.results.pointSets, f.results.points, f.results.fiapErr, nil
}

func (f *mockFetchClient) FetchOldest(fromDate *time.Time, untilDate *time.Time, ids ...string) (pointSets map[string]fiapmodel.ProcessedPointSet, points map[string][]fiapmodel.Value, fiapErr *fiapmodel.Error, err error) {
	if f.failOldest {
		return nil, nil, nil, errors.New("test FetchOldest error")
	}
	f.actualArguments = &fetchClientArguments{
		fromDate:  fromDate,
		untilDate: untilDate,
		ids:       ids,
	}
	return f.results.pointSets, f.results.points, f.results.fiapErr, nil
}

func TestFetchWithDateRange(t *testing.T) {
	query := &backend.DataQuery{
		RefID: "A",
	}
	fromTime := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC)
	fetchClient := mockFetchClient{}
	cli := ClientImpl{Client: &fetchClient}
	t.Run("Normal", func(t *testing.T) {
		t.Run("Period", func(t *testing.T) {
			dataRange := dsmodel.Period
			fetchClient.failLatest, fetchClient.failOldest, fetchClient.failDateRange = true, true, false
			fetchClient.results = &fetchClientResults{
				pointSets: map[string]fiapmodel.ProcessedPointSet{},
				fiapErr:   nil,
			}
			t.Run("FloatValues", func(t *testing.T) {
				pointIDs := []dsmodel.PointID{
					{Value: "id_a"},
				}
				fetchClient.results.points = map[string][]fiapmodel.Value{
					"id_a": {
						{
							Time:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
							Value: "33.4",
						},
						{
							Time:  time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC),
							Value: "-25.4",
						},
						{
							Time:  time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
							Value: "-130",
						},
						{
							Time:  time.Date(2024, 5, 22, 0, 0, 0, 0, time.UTC),
							Value: "0",
						},
						{
							Time:  time.Date(2024, 5, 29, 0, 0, 0, 0, time.UTC),
							Value: "-1.602e-19",
						},
						{
							Time:  time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
							Value: "0.0",
						},
					},
				}

				resp := &backend.DataResponse{}
				err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
				if err != nil {
					t.Error(err)
				}

				checkFrame(resp, fetchClient.results.points, query, data.FieldTypeFloat64, func(message string) {
					t.Error(message)
				})
				if fetchClient.actualArguments.fromDate == nil {
					t.Errorf("expected fromDate is %s but nil", fromTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.fromDate.Equal(fromTime) {
					t.Errorf("expected fromDate is %s but %s", fromTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.fromDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.untilDate == nil {
					t.Errorf("expected untilDate is %s but nil", toTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.untilDate.Equal(toTime) {
					t.Errorf("expected untilDate is %s but %s", toTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.untilDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.ids == nil {
					t.Errorf("point IDs are nil")
				} else {
					if len(fetchClient.actualArguments.ids) != len(pointIDs) {
						t.Errorf("expected pointID length is %d but %d", len(pointIDs), len(fetchClient.actualArguments.ids))
					}
					for i := range fetchClient.actualArguments.ids {
						if fetchClient.actualArguments.ids[i] != pointIDs[i].Value {
							t.Errorf("expected pointID[%d] is %s but %s", i, pointIDs[i].Value, fetchClient.actualArguments.ids[i])
						}
					}
				}
			})
			t.Run("CompoundValues", func(t *testing.T) {
				pointIDs := []dsmodel.PointID{
					{Value: "id_a"},
				}
				fetchClient.results.points = map[string][]fiapmodel.Value{
					"id_a": {
						{
							Time:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
							Value: "33.4",
						},
						{
							Time:  time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
							Value: "piyo",
						},
					},
				}

				resp := &backend.DataResponse{}
				err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
				if err != nil {
					t.Error(err)
				}

				checkFrame(resp, fetchClient.results.points, query, data.FieldTypeString, func(message string) {
					t.Error(message)
				})
				if fetchClient.actualArguments.fromDate == nil {
					t.Errorf("expected fromDate is %s but nil", fromTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.fromDate.Equal(fromTime) {
					t.Errorf("expected fromDate is %s but %s", fromTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.fromDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.untilDate == nil {
					t.Errorf("expected untilDate is %s but nil", toTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.untilDate.Equal(toTime) {
					t.Errorf("expected untilDate is %s but %s", toTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.untilDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.ids == nil {
					t.Errorf("point IDs are nil")
				} else {
					if len(fetchClient.actualArguments.ids) != len(pointIDs) {
						t.Errorf("expected pointID length is %d but %d", len(pointIDs), len(fetchClient.actualArguments.ids))
					}
					for i := range fetchClient.actualArguments.ids {
						if fetchClient.actualArguments.ids[i] != pointIDs[i].Value {
							t.Errorf("expected pointID[%d] is %s but %s", i, pointIDs[i].Value, fetchClient.actualArguments.ids[i])
						}
					}
				}
			})
			t.Run("StringValues", func(t *testing.T) {
				pointIDs := []dsmodel.PointID{
					{Value: "id_a"},
				}
				fetchClient.results.points = map[string][]fiapmodel.Value{
					"id_a": {
						{
							Time:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
							Value: "hoge",
						},
						{
							Time:  time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
							Value: "piyo",
						},
					},
				}

				resp := &backend.DataResponse{}
				err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
				if err != nil {
					t.Error(err)
				}

				checkFrame(resp, fetchClient.results.points, query, data.FieldTypeString, func(message string) {
					t.Error(message)
				})
				if fetchClient.actualArguments.fromDate == nil {
					t.Errorf("expected fromDate is %s but nil", fromTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.fromDate.Equal(fromTime) {
					t.Errorf("expected fromDate is %s but %s", fromTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.fromDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.untilDate == nil {
					t.Errorf("expected untilDate is %s but nil", toTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.untilDate.Equal(toTime) {
					t.Errorf("expected untilDate is %s but %s", toTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.untilDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.ids == nil {
					t.Errorf("point IDs are nil")
				} else {
					if len(fetchClient.actualArguments.ids) != len(pointIDs) {
						t.Errorf("expected pointID length is %d but %d", len(pointIDs), len(fetchClient.actualArguments.ids))
					}
					for i := range fetchClient.actualArguments.ids {
						if fetchClient.actualArguments.ids[i] != pointIDs[i].Value {
							t.Errorf("expected pointID[%d] is %s but %s", i, pointIDs[i].Value, fetchClient.actualArguments.ids[i])
						}
					}
				}
			})
		})
		t.Run("Latest", func(t *testing.T) {
			dataRange := dsmodel.Latest
			fetchClient.failLatest, fetchClient.failOldest, fetchClient.failDateRange = false, true, true
			fetchClient.results = &fetchClientResults{
				pointSets: map[string]fiapmodel.ProcessedPointSet{},
				fiapErr:   nil,
			}
			t.Run("SinglePointID", func(t *testing.T) {
				pointIDs := []dsmodel.PointID{
					{Value: "id_c"},
				}
				fetchClient.results.points = map[string][]fiapmodel.Value{
					"id_c": {
						{
							Time:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
							Value: "33.4",
						},
						{
							Time:  time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
							Value: "piyo",
						},
					},
				}

				resp := &backend.DataResponse{}
				err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
				if err != nil {
					t.Error(err)
				}

				checkFrame(resp, fetchClient.results.points, query, data.FieldTypeString, func(message string) {
					t.Error(message)
				})
				if fetchClient.actualArguments.fromDate == nil {
					t.Errorf("expected fromDate is %s but nil", fromTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.fromDate.Equal(fromTime) {
					t.Errorf("expected fromDate is %s but %s", fromTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.fromDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.untilDate == nil {
					t.Errorf("expected untilDate is %s but nil", toTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.untilDate.Equal(toTime) {
					t.Errorf("expected untilDate is %s but %s", toTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.untilDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.ids == nil {
					t.Errorf("point IDs are nil")
				} else {
					if len(fetchClient.actualArguments.ids) != len(pointIDs) {
						t.Errorf("expected pointID length is %d but %d", len(pointIDs), len(fetchClient.actualArguments.ids))
					}
					for i := range fetchClient.actualArguments.ids {
						if fetchClient.actualArguments.ids[i] != pointIDs[i].Value {
							t.Errorf("expected pointID[%d] is %s but %s", i, pointIDs[i].Value, fetchClient.actualArguments.ids[i])
						}
					}
				}
			})
			t.Run("MultiplePointIDs", func(t *testing.T) {
				pointIDs := []dsmodel.PointID{
					{Value: "id_a"},
					{Value: "id_b"},
				}
				fetchClient.results.points = map[string][]fiapmodel.Value{
					"id_a": {
						{
							Time:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
							Value: "hoge",
						},
						{
							Time:  time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
							Value: "piyo",
						},
					},
					"id_b": {
						{
							Time:  time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC),
							Value: "33.4",
						},
						{
							Time:  time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
							Value: "fuga",
						},
					},
				}

				resp := &backend.DataResponse{}
				err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
				if err != nil {
					t.Error(err)
				}

				checkFrame(resp, fetchClient.results.points, query, data.FieldTypeString, func(message string) {
					t.Error(message)
				})
				if fetchClient.actualArguments.fromDate == nil {
					t.Errorf("expected fromDate is %s but nil", fromTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.fromDate.Equal(fromTime) {
					t.Errorf("expected fromDate is %s but %s", fromTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.fromDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.untilDate == nil {
					t.Errorf("expected untilDate is %s but nil", toTime.Format("2006-01-02 15:04:05"))
				} else if !fetchClient.actualArguments.untilDate.Equal(toTime) {
					t.Errorf("expected untilDate is %s but %s", toTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.untilDate.Format("2006-01-02 15:04:05"))
				}
				if fetchClient.actualArguments.ids == nil {
					t.Errorf("point IDs are nil")
				} else {
					if len(fetchClient.actualArguments.ids) != len(pointIDs) {
						t.Errorf("expected pointID length is %d but %d", len(pointIDs), len(fetchClient.actualArguments.ids))
					}
					for i := range fetchClient.actualArguments.ids {
						if fetchClient.actualArguments.ids[i] != pointIDs[i].Value {
							t.Errorf("expected pointID[%d] is %s but %s", i, pointIDs[i].Value, fetchClient.actualArguments.ids[i])
						}
					}
				}
			})
		})
		t.Run("Oldest", func(t *testing.T) {
			dataRange := dsmodel.Oldest
			fetchClient.failLatest, fetchClient.failOldest, fetchClient.failDateRange = true, false, true
			fetchClient.results = &fetchClientResults{
				pointSets: map[string]fiapmodel.ProcessedPointSet{},
				fiapErr:   nil,
			}

			pointIDs := []dsmodel.PointID{
				{Value: "id_a"},
			}
			fetchClient.results.points = map[string][]fiapmodel.Value{
				"id_a": {
					{
						Time:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
						Value: "33.4",
					},
					{
						Time:  time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
						Value: "piyo",
					},
				},
			}

			resp := &backend.DataResponse{}
			err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
			if err != nil {
				t.Error(err)
			}

			checkFrame(resp, fetchClient.results.points, query, data.FieldTypeString, func(message string) {
				t.Error(message)
			})
			if fetchClient.actualArguments.fromDate == nil {
				t.Errorf("expected fromDate is %s but nil", fromTime.Format("2006-01-02 15:04:05"))
			} else if !fetchClient.actualArguments.fromDate.Equal(fromTime) {
				t.Errorf("expected fromDate is %s but %s", fromTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.fromDate.Format("2006-01-02 15:04:05"))
			}
			if fetchClient.actualArguments.untilDate == nil {
				t.Errorf("expected untilDate is %s but nil", toTime.Format("2006-01-02 15:04:05"))
			} else if !fetchClient.actualArguments.untilDate.Equal(toTime) {
				t.Errorf("expected untilDate is %s but %s", toTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.untilDate.Format("2006-01-02 15:04:05"))
			}
			if fetchClient.actualArguments.ids == nil {
				t.Errorf("point IDs are nil")
			} else {
				if len(fetchClient.actualArguments.ids) != len(pointIDs) {
					t.Errorf("expected pointID length is %d but %d", len(pointIDs), len(fetchClient.actualArguments.ids))
				}
				for i := range fetchClient.actualArguments.ids {
					if fetchClient.actualArguments.ids[i] != pointIDs[i].Value {
						t.Errorf("expected pointID[%d] is %s but %s", i, pointIDs[i].Value, fetchClient.actualArguments.ids[i])
					}
				}
			}
		})

	})
	t.Run("Error", func(t *testing.T) {
		t.Run("FetchFailed", func(t *testing.T) {
			dataRange := dsmodel.Period
			fetchClient.failLatest, fetchClient.failOldest, fetchClient.failDateRange = true, true, true

			pointIDs := []dsmodel.PointID{
				{Value: "id_a"},
			}
			fetchClient.results = nil

			resp := &backend.DataResponse{}
			err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
			if expectedErr := "test FetchLatest error"; err == nil {
				t.Errorf("expected error is %s but nil", expectedErr)
			} else if !strings.Contains(err.Error(), expectedErr) {
				t.Errorf("expected error is %s but %s", expectedErr, err.Error())

			}

			if len(resp.Frames) != 0 {
				t.Errorf("frames expects to have no fields")
			}
		})
		t.Run("FiapError", func(t *testing.T) {
			dataRange := dsmodel.Period
			fetchClient.failLatest, fetchClient.failOldest, fetchClient.failDateRange = true, true, false

			pointIDs := []dsmodel.PointID{
				{Value: "id_a"},
			}
			fetchClient.results = &fetchClientResults{
				pointSets: map[string]fiapmodel.ProcessedPointSet{},
				points: map[string][]fiapmodel.Value{
					"id_a": {
						{
							Time:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
							Value: "33.4",
						},
						{
							Time:  time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC),
							Value: "piyo",
						},
					},
				},
				fiapErr: &fiapmodel.Error{Type: "test_type", Value: "test_value"},
			}

			resp := &backend.DataResponse{}
			err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
			if expectedErr := "fiap error: type test_type, value test_value"; err == nil {
				t.Errorf("expected error is %s but nil", expectedErr)
			} else if !strings.Contains(err.Error(), expectedErr) {
				t.Errorf("expected error is %s but %s", expectedErr, err.Error())

			}

			checkFrame(resp, fetchClient.results.points, query, data.FieldTypeString, func(message string) {
				t.Error(message)
			})
			if fetchClient.actualArguments.fromDate == nil {
				t.Errorf("expected fromDate is %s but nil", fromTime.Format("2006-01-02 15:04:05"))
			} else if !fetchClient.actualArguments.fromDate.Equal(fromTime) {
				t.Errorf("expected fromDate is %s but %s", fromTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.fromDate.Format("2006-01-02 15:04:05"))
			}
			if fetchClient.actualArguments.untilDate == nil {
				t.Errorf("expected untilDate is %s but nil", toTime.Format("2006-01-02 15:04:05"))
			} else if !fetchClient.actualArguments.untilDate.Equal(toTime) {
				t.Errorf("expected untilDate is %s but %s", toTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.untilDate.Format("2006-01-02 15:04:05"))
			}
			if fetchClient.actualArguments.ids == nil {
				t.Errorf("point IDs are nil")
			} else {
				if len(fetchClient.actualArguments.ids) != len(pointIDs) {
					t.Errorf("expected pointID length is %d but %d", len(pointIDs), len(fetchClient.actualArguments.ids))
				}
				for i := range fetchClient.actualArguments.ids {
					if fetchClient.actualArguments.ids[i] != pointIDs[i].Value {
						t.Errorf("expected pointID[%d] is %s but %s", i, pointIDs[i].Value, fetchClient.actualArguments.ids[i])
					}
				}
			}
		})
		t.Run("ReturnPointSets", func(t *testing.T) {
			dataRange := dsmodel.Period
			fetchClient.failLatest, fetchClient.failOldest, fetchClient.failDateRange = true, true, false

			pointIDs := []dsmodel.PointID{
				{Value: "id_w"},
			}
			fetchClient.results = &fetchClientResults{
				pointSets: map[string]fiapmodel.ProcessedPointSet{
					"id_w": {
						PointSetID: []string{"id_x", "id_y", "id_z"},
						PointID:    []string{"id_a", "id_b"},
					},
				},
				points:  map[string][]fiapmodel.Value{},
				fiapErr: nil,
			}

			resp := &backend.DataResponse{}
			err := cli.FetchWithDateRange(resp, dataRange, &fromTime, &toTime, pointIDs, query)
			if expectedErr1, expectedErr2 := "point id 'id_w' provides point sets", "point id 'id_w' not provides point data"; err == nil {
				t.Errorf("expected error is %s and %s but nil", expectedErr1, expectedErr2)
			} else if !strings.Contains(err.Error(), expectedErr1) {
				t.Errorf("expected error is %s and %s but %s", expectedErr1, expectedErr2, err.Error())
			} else if !strings.Contains(err.Error(), expectedErr2) {
				t.Errorf("expected error is %s and %s but %s", expectedErr1, expectedErr2, err.Error())
			}

			if len(resp.Frames) != 0 {
				t.Errorf("frames expects to have no fields")
			}
			if fetchClient.actualArguments.fromDate == nil {
				t.Errorf("expected fromDate is %s but nil", fromTime.Format("2006-01-02 15:04:05"))
			} else if !fetchClient.actualArguments.fromDate.Equal(fromTime) {
				t.Errorf("expected fromDate is %s but %s", fromTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.fromDate.Format("2006-01-02 15:04:05"))
			}
			if fetchClient.actualArguments.untilDate == nil {
				t.Errorf("expected untilDate is %s but nil", toTime.Format("2006-01-02 15:04:05"))
			} else if !fetchClient.actualArguments.untilDate.Equal(toTime) {
				t.Errorf("expected untilDate is %s but %s", toTime.Format("2006-01-02 15:04:05"), fetchClient.actualArguments.untilDate.Format("2006-01-02 15:04:05"))
			}
			if fetchClient.actualArguments.ids == nil {
				t.Errorf("point IDs are nil")
			} else {
				if len(fetchClient.actualArguments.ids) != len(pointIDs) {
					t.Errorf("expected pointID length is %d but %d", len(pointIDs), len(fetchClient.actualArguments.ids))
				}
				for i := range fetchClient.actualArguments.ids {
					if fetchClient.actualArguments.ids[i] != pointIDs[i].Value {
						t.Errorf("expected pointID[%d] is %s but %s", i, pointIDs[i].Value, fetchClient.actualArguments.ids[i])
					}
				}
			}
		})
	})
}

func checkFrame(resp *backend.DataResponse, points map[string]([]fiapmodel.Value), query *backend.DataQuery, expectedType data.FieldType, cb func(string)) {
	for pointID := range points {
		expectedFrameName := query.RefID + ":" + pointID
		hasFrame := false
		for _, frame := range resp.Frames {
			if frame.Name == expectedFrameName {
				hasFrame = true
				hasTime, hasValue := false, false
				for _, field := range frame.Fields {
					switch field.Name {
					case "time":
						if hasTime {
							cb(fmt.Sprintf("field '%s' is duplicated in the frame '%s'", "time", expectedFrameName))
						} else {
							hasTime = true
						}
					case pointID:
						if hasValue {
							cb(fmt.Sprintf("field '%s' is duplicated in the frame '%s'", pointID, expectedFrameName))
						} else {
							if field.Type() == expectedType {
								hasValue = true
							} else {
								cb(fmt.Sprintf("field '%s' expects type '%s' but '%s' in the frame '%s'", pointID, expectedType.String(), field.Type().String(), expectedFrameName))
							}
						}
					default:
						cb(fmt.Sprintf("there is unexpected field '%s' in the frame '%s'", field.Name, expectedFrameName))
					}
				}
				if !hasTime {
					cb(fmt.Sprintf("field '%s' is not found in the frame '%s'", "time", expectedFrameName))
				}
				if !hasValue {
					cb(fmt.Sprintf("field '%s' is not found in the frame '%s'", pointID, expectedFrameName))
				}
			}
		}
		if !hasFrame {
			cb(fmt.Sprintf("point id '%s' is not found in frames", pointID))
		}
	}
}
