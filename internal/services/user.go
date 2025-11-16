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
}

type UserService struct {
	userRepo repositories.UserRepo
	tm       *transaction.Manager
}

func NewUserService(userRepo repositories.UserRepo, tm *transaction.Manager) UserService {
	return UserService{userRepo: userRepo, tm: tm}
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
