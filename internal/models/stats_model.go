package models

type Stats struct {
	ByUser map[string]int `json:"user_stats"`
	ByPR   map[string]int `json:"pr_stats"`
}
