package services

import (
	"context"
	"errors"
	"log"

	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/repositories"
	"github.com/linspacestrom/InterShipAv/internal/transaction"
	"github.com/linspacestrom/InterShipAv/internal/validateError"
)

type PRSer interface {
	Create(ctx context.Context, createPr domain.PullRequestCreate) (domain.PullRequestRead, error)
	Merge(ctx context.Context, prMerger domain.PRMerge) (domain.PRMergeRead, error)
	Reassign(ctx context.Context, pr domain.PRReassign) (domain.PrReassignRead, error)
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

		users, err := s.userRepo.GetNewReviewers(ctx, author.TeamName, author.Id)
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

func (s *PRService) Merge(ctx context.Context, prMerger domain.PRMerge) (domain.PRMergeRead, error) {
	var pr domain.PRMergeRead
	err := s.tm.Do(ctx, func(ctx context.Context) error {
		currentPr, err := s.prRepo.GetById(ctx, prMerger.Id)
		if err != nil {
			return err
		}

		prMerged, err := s.prRepo.Merge(ctx, prMerger.Id)
		if err != nil {
			return err
		}

		users, err := s.prRepo.GetReviewersById(ctx, currentPr.Id)
		if err != nil {
			return err
		}

		pr = prMerged
		pr.AssignReviewerIds = users
		return nil
	})

	if err != nil {
		return pr, err
	}

	return pr, nil
}

func (s *PRService) Reassign(ctx context.Context, pr domain.PRReassign) (domain.PrReassignRead, error) {
	var prReassign domain.PrReassignRead
	err := s.tm.Do(ctx, func(ctx context.Context) error {
		currentPR, err := s.prRepo.GetById(ctx, pr.Id)
		log.Println(currentPR)
		if err != nil {
			return err
		}
		if currentPR.Status == domain.StatusMerged {
			return validateError.PrMergedExist
		}

		author, err := s.userRepo.GetById(ctx, currentPR.AuthorId)
		if err != nil {
			return err
		}

		oldReviewer, err := s.userRepo.GetById(ctx, pr.OldUserId)
		if err != nil {
			return err
		}

		if author.TeamName != oldReviewer.TeamName {
			return validateError.UserNotAssignToTeam
		}

		reviewerIds, err := s.prRepo.GetReviewersById(ctx, pr.Id)
		if err != nil {
			return err
		}

		found := false
		for _, id := range reviewerIds {
			if id == oldReviewer.Id {
				found = true
				break
			}
		}

		if !found {
			return validateError.UserNotAssignReviewer
		}

		newReviewerId, err := s.userRepo.GetNewReviewer(ctx, author.Id, author.TeamName, reviewerIds)
		if err != nil {
			return err
		}

		err = s.prRepo.Reassign(ctx, pr.Id, newReviewerId, pr.OldUserId)
		if err != nil {
			return err
		}

		currentPR, err = s.prRepo.GetById(ctx, pr.Id)
		if err != nil {
			return err
		}

		reviewerIds, err = s.prRepo.GetReviewersById(ctx, pr.Id)
		if err != nil {
			return err
		}

		prReassign.ReplacedId = newReviewerId
		prReassign.PullRequest = currentPR
		prReassign.PullRequest.AssignReviewerIds = reviewerIds

		return nil

	})

	if err != nil {
		return prReassign, err
	}

	return prReassign, nil

}
