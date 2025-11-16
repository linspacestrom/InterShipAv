package domain

import "time"

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

type PullRequestReviewRead struct {
	Id       string
	Name     string
	AuthorId string
	Status   PullRequestStatus
}

type PRMerge struct {
	Id string
}

type PRMergeRead struct {
	Id                string
	Name              string
	AuthorId          string
	Status            PullRequestStatus
	AssignReviewerIds []string
	MergedAt          *time.Time
}

type PRReassign struct {
	Id        string
	OldUserId string
}

type PrReassignRead struct {
	PullRequest PullRequestRead
	ReplacedId  string
}
