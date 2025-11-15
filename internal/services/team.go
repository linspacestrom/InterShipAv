package services

import (
	"context"
	"errors"

	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/repositories"
	"github.com/linspacestrom/InterShipAv/internal/transaction"
	"github.com/linspacestrom/InterShipAv/internal/validateError"
)

type TeamSer interface {
	Create(ctx context.Context, team domain.Team) (domain.Team, error)
	GetByName(ctx context.Context, name string) (domain.Team, error)
}

type TeamService struct {
	teamRepo repositories.TeamRepo
	userRepo repositories.UserRepo
	tm       *transaction.Manager
}

func NewTeamService(teamRepo repositories.TeamRepo, userRepo repositories.UserRepo, tm *transaction.Manager) TeamService {
	return TeamService{teamRepo: teamRepo, userRepo: userRepo, tm: tm}
}

func (s *TeamService) Create(ctx context.Context, team domain.Team) (domain.Team, error) {
	var createdTeam domain.Team

	err := s.tm.Do(ctx, func(ctx context.Context) error {
		_, err := s.teamRepo.GetByName(ctx, team.Name)
		if err == nil {
			return validateError.ErrTeamExists
		}
		if !errors.Is(err, validateError.TeamNotFound) {
			return err
		}

		createdTeam, err = s.teamRepo.Create(ctx, team)
		if err != nil {
			return err
		}

		if err := s.userRepo.AddUsersToTeam(ctx, team.Members, createdTeam.Name); err != nil {
			return err
		}

		users, err := s.userRepo.GetUserByTeamName(ctx, createdTeam.Name)
		if err != nil {
			return err
		}

		createdTeam.Members = users

		return nil
	})

	if err != nil {
		return domain.Team{}, err
	}

	return createdTeam, nil
}

func (s *TeamService) GetByName(ctx context.Context, name string) (domain.Team, error) {
	var team domain.Team

	err := s.tm.Do(ctx, func(ctx context.Context) error {
		var errInner error
		team, errInner = s.teamRepo.GetByName(ctx, name)
		if errors.Is(errInner, validateError.TeamNotFound) {
			return errInner
		}
		users, err := s.userRepo.GetUserByTeamName(ctx, team.Name)
		if err != nil {
			return nil
		}
		team.Members = users

		return nil
	})

	if err != nil {
		return domain.Team{}, err
	}

	return team, err
}
