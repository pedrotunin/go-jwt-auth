package models

import "errors"

type RefreshTokenID = int
type RefreshTokenContent = string
type RefreshTokenStatus = string

var RefreshTokenStatusActive RefreshTokenStatus = "active"
var RefreshTokenStatusInactive RefreshTokenStatus = "inactive"

var ErrRefreshTokenInvalid = errors.New("refresh token is invalid")

type RefreshToken struct {
	ID      RefreshTokenID
	Content RefreshTokenContent
	Status  RefreshTokenStatus
}
