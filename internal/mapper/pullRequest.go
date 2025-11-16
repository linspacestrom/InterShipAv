package mapper

import (
	"time"

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

func DTOtoPRMerge(req dto.PRMergeRequest) domain.PRMerge {
	return domain.PRMerge{Id: req.Id}
}

func PRMergeToDTO(res domain.PRMergeRead) dto.PRMergeResponse {
	var mergedAt *time.Time
	if res.MergedAt != nil {
		t := res.MergedAt.In(time.Local).Truncate(time.Second)
		mergedAt = &t
	}

	return dto.PRMergeResponse{
		Id:                res.Id,
		Name:              res.Name,
		AuthorId:          res.AuthorId,
		Status:            string(res.Status),
		AssignReviewerIds: res.AssignReviewerIds,
		MergedAt:          mergedAt,
	}
}
