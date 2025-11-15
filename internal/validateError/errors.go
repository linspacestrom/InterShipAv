package validateError

import "errors"

var ErrTeamExists = errors.New("team already exist")
var TeamNotFound = errors.New("команда не найдена")
