package models

type Team struct {
	Name    string `db:"team_name" json:"team_name"`
	Members []User `json:"members,omitempty"` // только для ответов API
}
