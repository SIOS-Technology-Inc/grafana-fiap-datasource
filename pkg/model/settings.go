package model

import "time"

type FiapDatasourceSettings struct {
	Url            string `json:"url"`
	ServerTimezone string `json:"server_timezone"`
}

const serverTimezoneLayout = "-07:00"

func (s *FiapDatasourceSettings) GetLocation() (*time.Location, error) {
	if s.ServerTimezone == "" {
		return time.UTC, nil
	} else {
		if dt, err := time.Parse(serverTimezoneLayout, s.ServerTimezone); err == nil {
			return dt.Location(), nil
		} else {
			return nil, err
		}
	}
}
