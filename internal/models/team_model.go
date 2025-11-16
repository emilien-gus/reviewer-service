package models

type Team struct {
	Name    string        `json:"team_name" binding:"required"`
	Members *[]TeamMember `json:"members,omitempty" binding:"required"`
}

type TeamMember struct {
	ID       string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}
