package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrAuthFailed = errors.New("authentcation failed")
