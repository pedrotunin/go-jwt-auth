package repositories

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type PSQLEmailVerificationTokenRepository struct {
	db *sql.DB
}

func NewPSQLEmailVerificationTokenRepository(db *sql.DB) *PSQLEmailVerificationTokenRepository {
	return &PSQLEmailVerificationTokenRepository{
		db: db,
	}
}

func (repo *PSQLEmailVerificationTokenRepository) CreateVerificationToken(token *models.EmailVerificationToken) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("CreateVerificationToken: error creating transaction: %s", err.Error())
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO email_verification_tokens (content, user_id, expires_at) VALUES ($1, $2, $3);")
	if err != nil {
		log.Printf("CreateVerificationToken: error creating statement: %s", err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.Content, token.UserID, token.ExpiresAt)
	if err != nil {
		log.Printf("CreateVerificationToken: error executing query: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("CreateVerificationToken: error during commmit: %s", err.Error())
		tx.Rollback()
		return err
	}

	log.Printf("CreateVerificationToken: email verification token created")
	return nil
}

func (repo *PSQLEmailVerificationTokenRepository) GetVerificationTokenByContent(content models.EmailVerificationTokenContent) (*models.EmailVerificationToken, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("GetVerificationTokenByContent: error creating transaction: %s", err.Error())
		return nil, err
	}

	query := "SELECT id, content, user_id, expires_at FROM email_verification_tokens WHERE content=$1 AND is_used=FALSE;"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Printf("GetVerificationTokenByContent: error creating statement: %s", err.Error())
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	var resId, resUserId int
	var resContent string
	var resExpiresAt time.Time
	err = stmt.QueryRow(content).Scan(&resId, &resContent, &resUserId, &resExpiresAt)
	if err != nil {
		log.Printf("GetVerificationTokenByContent: error executing query: %s", err.Error())
		tx.Rollback()

		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrVerifyTokenNotFound
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("GetVerificationTokenByContent: error during commmit: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	evToken := models.EmailVerificationToken{
		ID:        resId,
		Content:   resContent,
		UserID:    resUserId,
		ExpiresAt: resExpiresAt,
	}

	log.Printf("GetVerificationTokenByContent: email verification token created")
	return &evToken, nil

}

func (repo *PSQLEmailVerificationTokenRepository) SetTokenToUsed(token *models.EmailVerificationToken) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("SetTokenToUsed: error creating transaction: %s", err.Error())
		return err
	}

	stmt, err := tx.Prepare("UPDATE email_verification_tokens SET is_used=TRUE WHERE id=$1;")
	if err != nil {
		log.Printf("SetTokenToUsed: error creating statement: %s", err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.ID)
	if err != nil {
		log.Printf("SetTokenToUsed: error executing query: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("SetTokenToUsed: error during commmit: %s", err.Error())
		tx.Rollback()
		return err
	}

	log.Printf("SetTokenToUsed: token set to used")
	return nil
}
