package utils

import "errors"

// General Errors
var ErrInternalServerError = errors.New("internal server error. contact system admin")

// User Errors
var ErrUserNotFound = errors.New("user not found")
var ErrUserEmailAlreadyExists = errors.New("user's email already exists in our database")

// Refresh Token Errors
var ErrRefreshTokenInvalid = errors.New("refresh token is invalid")
var ErrRefreshTokenNotFound = errors.New("refresh token not found")
