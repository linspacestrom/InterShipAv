package mapper

import (
	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/dto"
)

func DTOToPrCreate(req dto.PRCreateRequest) domain.PullRequestCreate {
	return domain.PullRequestCreate{
		Id:       req.Id,
		Name:     req.Name,
		AuthorId: req.AuthorId,
	}
}

func PRReadToDTO(res domain.PullRequestRead) dto.PRCreateResponse {
	return dto.PRCreateResponse{
		Id:                res.Id,
		Name:              res.Name,
		AuthorId:          res.AuthorId,
		Status:            string(res.Status),
		AssignReviewerIds: res.AssignReviewerIds,
	}
}
