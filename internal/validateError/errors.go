package validateError

import (
	"errors"
)

var ErrTeamExists = errors.New("team already exist")
var TeamNotFound = errors.New("команда не найдена")
var UserNotFound = errors.New("пользователь не найден")
var ErrPRExist = errors.New("PR id already exists")
var ErrPrNotExist = errors.New("PR not found")
