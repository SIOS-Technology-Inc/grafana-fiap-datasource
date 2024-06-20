package model

import (
	"encoding/json"
	"time"
)

type FiapQuery struct {
	PointIDs  []PointID     `json:"point_ids"`
	DataRange DataRangeType `json:"data_range"`
	StartTime LinkedTime    `json:"start_time"`
	EndTime   LinkedTime    `json:"end_time"`
}

type PointID struct {
	Value string `json:"point_id"`
}

type DataRangeType string

const Period DataRangeType = "period"
const Latest DataRangeType = "latest"
const Oldest DataRangeType = "oldest"

type LinkedTime struct {
	FixedTime     *time.Time
	LinkDashboard bool
}

const frontendDatetimeLayout = "2006-01-02 15:04:05"

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
	} else if dt, err := time.Parse(frontendDatetimeLayout, r.FixedTime); err == nil {
		l.FixedTime = &dt
	} else {
		return err
	}
	l.LinkDashboard = r.LinkDashboard
	return nil
}
