package dto

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

type PRReadRequest struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
	Status   string `json:"status"`
}
