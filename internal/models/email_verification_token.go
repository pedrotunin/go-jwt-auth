package models

import "time"

type EmailVerificationTokenID = int
type EmailVerificationTokenContent = string

type EmailVerificationToken struct {
	ID        EmailVerificationTokenID
	Content   EmailVerificationTokenContent
	UserID    UserID
	CreatedAt time.Time
	ExpiresAt time.Time
	IsUsed    bool
}
