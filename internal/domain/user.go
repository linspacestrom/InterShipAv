package domain

type UserShort struct {
	Id       string
	IsActive bool
}

type User struct {
	Id       string
	Username string
	TeamName string
	IsActive bool
}

type UserReview struct {
	Id           string
	PullRequests []PullRequestReviewRead
}
