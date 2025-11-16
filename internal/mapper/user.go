package mapper

import (
	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/dto"
)

func UserToDTO(user domain.User) dto.UserResponse {
	return dto.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func DTOToUserShort(user dto.UserRequest) domain.UserShort {
	return domain.UserShort{
		Id:       user.Id,
		IsActive: *user.IsActive,
	}
}

func DomainReviewToDTOReview(user domain.UserReview) dto.UserReviewResponse {
	pullRequests := make([]dto.PRReadRequest, 0, len(user.PullRequests))

	for _, pr := range user.PullRequests {
		pullRequests = append(pullRequests, dto.PRReadRequest{Id: pr.Id, Name: pr.Name, AuthorId: pr.AuthorId, Status: string(pr.Status)})
	}

	return dto.UserReviewResponse{
		Id:           user.Id,
		PullRequests: pullRequests,
	}
}
