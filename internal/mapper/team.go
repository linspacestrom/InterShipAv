package mapper

import (
	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/dto"
)

func TeamToDTO(team domain.Team) dto.GetTeamResponse {
	members := make([]dto.TeamMemberDTO, len(team.Members))
	for i, m := range team.Members {
		members[i] = TeamMemberToDTO(m)
	}
	return dto.GetTeamResponse{
		Name:    team.Name,
		Members: members,
	}
}

func CreateTeamToDTO(team domain.Team) dto.CreateTeamResponse {
	members := make([]dto.TeamMemberDTO, len(team.Members))
	for i, m := range team.Members {
		members[i] = TeamMemberToDTO(m)
	}
	return dto.CreateTeamResponse{
		Name:    team.Name,
		Members: members,
	}
}

func TeamMemberToDTO(member domain.TeamMember) dto.TeamMemberDTO {
	return dto.TeamMemberDTO{
		ID:       member.ID,
		Username: member.Username,
		IsActive: member.IsActive,
	}
}

func DTOToTeam(req dto.CreateTeamRequest) domain.Team {
	members := make([]domain.TeamMember, len(req.Members))
	for i, m := range req.Members {
		members[i] = DTOToTeamMember(m)
	}
	return domain.Team{
		Name:    req.Name,
		Members: members,
	}
}

func DTOToTeamMember(m dto.TeamMemberDTO) domain.TeamMember {
	return domain.TeamMember{
		ID:       m.ID,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}
