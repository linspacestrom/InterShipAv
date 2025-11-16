package dto

import "time"

type PRCreateRequest struct {
	Id       string `json:"pull_request_id" binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorId string `json:"author_id" binding:"required"`
}

type PRCreateResponse struct {
	Id                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignReviewerIds []string `json:"assign_reviewer"`
}

type PRReadResponse struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
	Status   string `json:"status"`
}

type PRMergeRequest struct {
	Id string `json:"pull_request_id" binding:"required"`
}

type PRMergeResponse struct {
	Id                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorId          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignReviewerIds []string   `json:"assigned_reviewers"`
	MergedAt          *time.Time `json:"mergedAt"`
}

type PRReassignRequest struct {
	Id        string `json:"pull_request_id" binding:"required"`
	OldUserId string `json:"old_user_id" binding:"required"`
}

type ReassignResponse struct {
	Id                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignReviewerIds []string `json:"assign_reviewer"`
}

type PrReassignResponse struct {
	PrRead     ReassignResponse `json:"pr"`
	ReplacedId string           `json:"replaced_by"`
}
