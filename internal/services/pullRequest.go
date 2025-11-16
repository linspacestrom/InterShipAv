package services

import (
	"context"
	"errors"

	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/repositories"
	"github.com/linspacestrom/InterShipAv/internal/transaction"
	"github.com/linspacestrom/InterShipAv/internal/validateError"
)

type PRSer interface {
	Create(ctx context.Context, createPr domain.PullRequestCreate) (domain.PullRequestRead, error)
}

type PRService struct {
	prRepo   repositories.PrRepo
	userRepo repositories.UserRepo
	teamRepo repositories.TeamRepo
	tm       *transaction.Manager
}

func NewPRService(prRepo repositories.PrRepo, userRepo repositories.UserRepo, teamRepo repositories.TeamRepo, tm *transaction.Manager) PRService {
	return PRService{prRepo: prRepo, userRepo: userRepo, teamRepo: teamRepo, tm: tm}
}

func (s *PRService) Create(ctx context.Context, createPr domain.PullRequestCreate) (domain.PullRequestRead, error) {
	var pr domain.PullRequestRead

	err := s.tm.Do(ctx, func(ctx context.Context) error {
		_, err := s.prRepo.GetById(ctx, createPr.Id)
		if err == nil {
			return validateError.ErrPRExist
		}
		if !errors.Is(err, validateError.ErrPrNotExist) {
			return err
		}

		author, err := s.userRepo.GetById(ctx, createPr.AuthorId)
		if err != nil {
			if errors.Is(err, validateError.UserNotFound) {
				return validateError.UserNotFound
			}
			return err
		}

		pr, err = s.prRepo.Create(ctx, createPr)
		if err != nil {
			return err
		}

		users, err := s.userRepo.GetReviewers(ctx, author.TeamName, author.Id)
		if err != nil {
			return err
		}

		reviewerIds, err := s.prRepo.AssignReviewers(ctx, createPr.Id, users)
		if err != nil {
			return err
		}

		pr.AssignReviewerIds = reviewerIds
		return nil
	})

	if err != nil {
		return domain.PullRequestRead{}, err
	}
	return pr, nil
}
