package domain

type TeamMember struct {
	ID       string
	Username string
	IsActive bool
}

type Team struct {
	Name    string
	Members []TeamMember
}

type TeamMemberCreate struct {
	ID       string
	Username string
	IsActive bool
}

type TeamMemberUpdate struct {
	Username string
	IsActive bool
}

type TeamMemberDetail struct {
	ID       string
	Username string
	IsActive bool
}

type TeamDetail struct {
	Name    string
	Members []TeamMember
}
