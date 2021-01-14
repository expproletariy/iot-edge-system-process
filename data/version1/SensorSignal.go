package version1

import "time"

type SensorSignal struct {
	Id          string     `json:"id"`
	SiteID      string     `json:"site_id"`
	PartitionID string     `json:"partition_id"`
	CreatedAt   time.Time  `json:"created_at"`
	Type        SensorType `json:"type"`
	Name        string     `json:"name"`
	Value       float64    `json:"value"`
}
