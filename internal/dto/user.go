package dto

type UserRequest struct {
	Id       string `json:"user_id" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type UserResponse struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserReviewResponse struct {
	Id           string          `json:"user_id"`
	PullRequests []PRReadRequest `json:"pull_requests"`
}
