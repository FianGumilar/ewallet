package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrAuthFailed = errors.New("authentcation failed")
var ErrUsernameExists = errors.New("username already exists")
var ErrOtpInvalid = errors.New("otp code invalid")
var ErrAccountNotFound = errors.New("account not found")
var ErrInquiryNotFound = errors.New("inquiry not found")
var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrCodeNotFound = errors.New("code not found")
