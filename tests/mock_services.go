package tests

import (
	"context"
	"errors"
	"sync"

	"github.com/linspacestrom/InterShipAv/internal/domain"
	"github.com/linspacestrom/InterShipAv/internal/services"
)

type FakeTeamService struct {
	createCalls map[string]int
	lock        sync.Mutex
}

func NewFakeTeamService() *FakeTeamService {
	return &FakeTeamService{createCalls: make(map[string]int)}
}

var _ services.TeamSer = (*FakeTeamService)(nil)

func (s *FakeTeamService) Create(ctx context.Context, team domain.Team) (domain.Team, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.createCalls[team.Name]++
	if s.createCalls[team.Name] > 1 {
		return domain.Team{}, errors.New("team already exists")
	}
	return team, nil
}

func (s *FakeTeamService) GetByName(ctx context.Context, name string) (domain.Team, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.createCalls[name] == 0 {
		return domain.Team{}, errors.New("team not found")
	}
	return domain.Team{Name: name, Members: []domain.TeamMember{}}, nil
}

type FakeUserService struct {
	registeredUsers map[string]domain.User
	lock            sync.Mutex
}

func NewFakeUserService() *FakeUserService {
	return &FakeUserService{registeredUsers: make(map[string]domain.User)}
}

var _ services.UserSer = (*FakeUserService)(nil)

func (s *FakeUserService) SetActive(ctx context.Context, id string, isActive bool) (domain.User, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, ok := s.registeredUsers[id]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}
	user.IsActive = isActive
	return user, nil
}

func (s *FakeUserService) GetReview(ctx context.Context, userId string) (domain.UserReview, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.registeredUsers[userId]
	if !ok {
		return domain.UserReview{}, errors.New("user not found")
	}
	return domain.UserReview{Id: userId, PullRequests: []domain.PullRequestReviewRead{}}, nil
}

type FakePRService struct {
	createCalls map[string]int
	createdPRs  map[string]domain.PullRequestRead
	lock        sync.Mutex
	users       *FakeUserService
}

func NewFakePRServiceWithUsers(users *FakeUserService) *FakePRService {
	return &FakePRService{
		createCalls: make(map[string]int),
		createdPRs:  make(map[string]domain.PullRequestRead),
		users:       users,
	}
}

var _ services.PRSer = (*FakePRService)(nil)

func (s *FakePRService) Create(ctx context.Context, createPr domain.PullRequestCreate) (domain.PullRequestRead, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.createCalls[createPr.Id]++
	if s.createCalls[createPr.Id] > 1 {
		return domain.PullRequestRead{}, errors.New("pr already exists")
	}
	if s.users != nil {
		s.users.lock.Lock()
		_, ok := s.users.registeredUsers[createPr.AuthorId]
		s.users.lock.Unlock()
		if !ok {
			return domain.PullRequestRead{}, errors.New("pr not found")
		}
	}
	pr := domain.PullRequestRead{
		Id: createPr.Id, Name: createPr.Name, AuthorId: createPr.AuthorId, Status: domain.StatusOpen, AssignReviewerIds: []string{"rev1", "rev2"},
	}
	s.createdPRs[createPr.Id] = pr
	return pr, nil
}

func (s *FakePRService) Merge(ctx context.Context, prMerger domain.PRMerge) (domain.PRMergeRead, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	pr, ok := s.createdPRs[prMerger.Id]
	if !ok {
		return domain.PRMergeRead{}, errors.New("pr not found")
	}
	pr.Status = domain.StatusMerged
	return domain.PRMergeRead{Id: pr.Id, Name: pr.Name, AuthorId: pr.AuthorId, Status: pr.Status, AssignReviewerIds: pr.AssignReviewerIds}, nil
}

func (s *FakePRService) Reassign(ctx context.Context, pr domain.PRReassign) (domain.PrReassignRead, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	origin, ok := s.createdPRs[pr.Id]
	if !ok {
		return domain.PrReassignRead{}, errors.New("pr not found")
	}
	if len(origin.AssignReviewerIds) == 0 {
		return domain.PrReassignRead{}, errors.New("no reviewers")
	}
	replacedId := origin.AssignReviewerIds[0]
	origin.AssignReviewerIds[0] = "new_user"
	return domain.PrReassignRead{PullRequest: origin, ReplacedId: replacedId}, nil
}
