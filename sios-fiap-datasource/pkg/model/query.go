package model

import (
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
	RawTime       string `json:"time"`
	LinkDashboard bool   `json:"link_dashboard"`
}

const frontendDatetimeLayout = "2006-01-02 15:04:05"

func (l *LinkedTime) GetTime(_ string) (*time.Time, error) {
	if l.RawTime == "" {
		return nil, nil
	} else {
		if dt, err := time.Parse(frontendDatetimeLayout, l.RawTime); err == nil {
			return &dt, nil
		} else {
			return nil, err
		}
	}
}
