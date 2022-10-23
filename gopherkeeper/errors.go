package gopherkeeper

import "errors"

var ErrUserLoginConflict = errors.New(`данный логин уже занят`)
