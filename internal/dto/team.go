package dto

type TeamMemberDTO struct {
	ID       string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type CreateTeamRequest struct {
	Name    string          `json:"team_name" binding:"required"`
	Members []TeamMemberDTO `json:"members" binding:"required,dive"`
}

type CreateTeamResponse struct {
	Name    string          `json:"team_name"`
	Members []TeamMemberDTO `json:"members"`
}

type GetTeamResponse struct {
	Name    string          `json:"team_name"`
	Members []TeamMemberDTO `json:"members"`
}
