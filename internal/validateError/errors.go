package validateError

import (
	"errors"
)

var ErrTeamExists = errors.New("team already exist")
var TeamNotFound = errors.New("team not found")
var UserNotFound = errors.New("user not found")
var ErrPRExist = errors.New("pull request id already exists")
var ErrPrNotExist = errors.New("pull request not found")
var PrMergedExist = errors.New("pull request already merged")
var UserNotAssignReviewer = errors.New("user not assign to pull request")
var NoCandidate = errors.New("no active replacement candidate in team")
var UserNotAssignToTeam = errors.New("user not assign to team")
var UserNotUniqueId = errors.New("users hasn`t unique ids")
