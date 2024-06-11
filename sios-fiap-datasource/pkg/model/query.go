package model

import (
	"encoding/json"
	"time"
)

type FiapQuery struct {
	PointIDs  []string      `json:"point_ids"`
	DataRange DataRangeType `json:"data_range"`
	StartTime LinkedTime    `json:"start_time"`
	EndTime   LinkedTime    `json:"end_time"`
}

type DataRangeType string

const Period DataRangeType = "period"
const Latest DataRangeType = "latest"
const Oldest DataRangeType = "oldest"

type LinkedTime struct {
	FixedTime     *time.Time
	LinkDashboard bool
}

func (l *LinkedTime) UnmarshalJSON(b []byte) error {
	var r struct {
		FixedTime     string `json:"time"`
		LinkDashboard bool   `json:"link_dashboard"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	if r.FixedTime == "" {
		l.FixedTime = nil
	} else if dt, err := time.Parse(time.RFC3339, r.FixedTime); err == nil {
		l.FixedTime = &dt
	} else {
		return err
	}
	l.LinkDashboard = r.LinkDashboard
	return nil
}
