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
