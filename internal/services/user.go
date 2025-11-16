package services

import (
	"context"

	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/repositories"
	"github.com/linspacestrom/InterShipAv/internal/transaction"
	"github.com/linspacestrom/InterShipAv/internal/validateError"
)

type UserSer interface {
	SetActive(ctx context.Context, id string, isActive bool) (domain.User, error)
	GetReview(ctx context.Context, userId string) (domain.UserReview, error)
}

type UserService struct {
	userRepo repositories.UserRepo
	prRepo   repositories.PrRepo
	tm       *transaction.Manager
}

func NewUserService(userRepo repositories.UserRepo, prRepo repositories.PrRepo, tm *transaction.Manager) UserService {
	return UserService{userRepo: userRepo, prRepo: prRepo, tm: tm}
}

func (s *UserService) SetActive(ctx context.Context, id string, isActive bool) (domain.User, error) {
	var user domain.User
	err := s.tm.Do(ctx, func(ctx context.Context) error {
		_, err := s.userRepo.GetById(ctx, id)
		if err != nil {
			return validateError.UserNotFound
		}

		user, err = s.userRepo.SetActiveById(ctx, id, isActive)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) GetReview(ctx context.Context, userId string) (domain.UserReview, error) {
	var reviewer domain.UserReview

	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return reviewer, nil
	}

	prIds, err := s.prRepo.GetPullRequestIdsByUserId(ctx, userId)
	if err != nil {
		return reviewer, nil
	}

	prReviews, err := s.prRepo.GetPullRequestsByIds(ctx, prIds)
	if err != nil {
		return reviewer, nil
	}

	reviewer.Id = user.Id
	reviewer.PullRequests = prReviews

	return reviewer, nil

}
