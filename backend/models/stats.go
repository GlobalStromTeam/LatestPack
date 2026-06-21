package models

type DailyStats struct {
	Date     string
	Launches int64
	Updates  int64
	Traffic  int64
}

type StatItem struct {
	Value  int64   `json:"value"`
	Change float64 `json:"change"`
	Unit   string  `json:"unit,omitempty"`
}

type StatsResponse struct {
	Launches StatItem `json:"launches"`
	Updates  StatItem `json:"updates"`
	Traffic  StatItem `json:"traffic"`
}
