package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrAuthFailed = errors.New("authentcation failed")
var ErrUsernameExists = errors.New("username already exists")
var ErrOtpInvalid = errors.New("otp code invalid")
