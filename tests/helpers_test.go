package tests

import "github.com/linspacestrom/InterShipAv/internal/domain"

func MakeTestTeam(name string, ids []string) domain.Team {
	members := []domain.TeamMember{}
	for _, id := range ids {
		members = append(members, domain.TeamMember{ID: id, Username: id + "name", IsActive: true})
	}
	return domain.Team{
		Name:    name,
		Members: members,
	}
}

func MakeTestUser(id, username, team string, active bool) domain.User {
	return domain.User{
		Id:       id,
		Username: username,
		TeamName: team,
		IsActive: active,
	}
}

func MakeTestPR(id, name, author string, reviewers bool) domain.PullRequestRead {
	ids := []string{}
	if reviewers {
		ids = []string{"rev1", "rev2"}
	}
	return domain.PullRequestRead{
		Id:                id,
		Name:              name,
		AuthorId:          author,
		Status:            domain.StatusOpen,
		AssignReviewerIds: ids,
	}
}
