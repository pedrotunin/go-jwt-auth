package models

type RefreshTokenID = int
type RefreshTokenContent = string
type RefreshTokenStatus = string

var RefreshTokenStatusActive RefreshTokenStatus = "active"
var RefreshTokenStatusInactive RefreshTokenStatus = "inactive"

type RefreshToken struct {
	ID      RefreshTokenID
	Content RefreshTokenContent
	Status  RefreshTokenStatus
}
