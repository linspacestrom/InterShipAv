package domain

type PullRequestStatus string

const (
	StatusOpen   PullRequestStatus = "OPEN"
	StatusMerged PullRequestStatus = "MERGED"
)

type PullRequestCreate struct {
	Id       string
	Name     string
	AuthorId string
}

type PullRequestRead struct {
	Id                string
	Name              string
	AuthorId          string
	Status            PullRequestStatus
	AssignReviewerIds []string
}
