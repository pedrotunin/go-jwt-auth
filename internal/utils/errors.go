package utils

import "errors"

// General Errors
var ErrInternalServerError = errors.New("internal server error. contact system admin")
var ErrAuthorizationHeaderNotFound = errors.New("Authorization header not found")
var ErrMultipleAuthorizationHeaders = errors.New("multiple Authorization headers not accepted")
var ErrAuthorizationHeaderMalformed = errors.New("Authorization header malformed")

// User Errors
var ErrUserNotFound = errors.New("user not found")
var ErrUserEmailAlreadyExists = errors.New("user's email already exists in our database")
var ErrUserInactive = errors.New("user is not active")
var ErrUserPending = errors.New("user is pending activation")
var ErrInvalidUserStatus = errors.New("user status is not valid")
var ErrInvalidUserID = errors.New("user ID is invalid")
var ErrUserIDsDoNotMatch = errors.New("user IDs do not match")

// Refresh Token Errors
var ErrRefreshTokenInvalid = errors.New("refresh token is invalid")
var ErrRefreshTokenNotFound = errors.New("refresh token not found")

// Token Errors
var ErrTokenInvalid = errors.New("invalid token")

// Password Errors
var ErrPasswordsNotMatch = errors.New("passwords don't match")
var ErrEmailPasswordIncorrect = errors.New("email or password incorrect")
var ErrPasswordTooShort = errors.New("password must be at least 8 characters long")

// E-mail Errors
var ErrInvalidEmail = errors.New("email is invalid")

// Verify Token Errors
var ErrVerifyTokenNotFound = errors.New("verify token not found")
var ErrVerifyTokenExpired = errors.New("verify token expired")
