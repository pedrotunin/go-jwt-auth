package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/pedrotunin/jwt-auth/internal/models"
	"github.com/pedrotunin/jwt-auth/internal/utils"
)

type PSQLRefreshTokenRepository struct {
	db *sql.DB
}

func NewPSQLRefreshTokenRepository(db *sql.DB) *PSQLRefreshTokenRepository {
	return &PSQLRefreshTokenRepository{
		db: db,
	}
}

func (repo *PSQLRefreshTokenRepository) CreateRefreshToken(token *models.RefreshToken) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("CreateRefreshToken: error creating transaction: %s", err.Error())
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO refresh_tokens (content, status) VALUES ($1, $2);")
	if err != nil {
		log.Printf("CreateRefreshToken: error creating statement: %s", err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.Content, token.Status)
	if err != nil {
		log.Printf("CreateRefreshToken: error executing query: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("CreateRefreshToken: error during commmit: %s", err.Error())
		tx.Rollback()
		return err
	}

	log.Printf("CreateRefreshToken: refresh token created")
	return nil
}

func (repo *PSQLRefreshTokenRepository) GetRefreshTokenByContent(content models.RefreshTokenContent) (*models.RefreshToken, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("GetRefreshTokenByContent: error creating transaction: %s", err.Error())
		return nil, fmt.Errorf("GetRefreshTokenByContent: error creating transaction: %w", err)
	}

	stmt, err := tx.Prepare("SELECT * FROM refresh_tokens WHERE content=$1;")
	if err != nil {
		log.Printf("GetRefreshTokenByContent: error creating statement: %s", err.Error())
		tx.Rollback()
		return nil, fmt.Errorf("GetRefreshTokenByContent: error creating prepared statement: %w", err)
	}
	defer stmt.Close()

	var resId int
	var resContent, resStatus string
	err = stmt.QueryRow(content).Scan(&resId, &resContent, &resStatus)
	if err != nil {
		log.Printf("GetRefreshTokenByContent: error executing query: %s", err.Error())
		tx.Rollback()

		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrRefreshTokenNotFound
		}

		return nil, fmt.Errorf("GetUserByEmail: error scanning query result: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("GetRefreshTokenByContent: error during commit: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	log.Printf("GetRefreshTokenByContent: refresh token found")
	return &models.RefreshToken{
		ID:      resId,
		Content: resContent,
		Status:  resStatus,
	}, nil
}

func (repo *PSQLRefreshTokenRepository) InvalidateRefreshTokenByContent(content models.RefreshTokenContent) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("InvalidateRefreshTokenByContent: error creating transaction: %s", err.Error())
		return fmt.Errorf("InvalidateRefreshTokenByContent: error creating transaction: %w", err)
	}

	stmt, err := tx.Prepare("UPDATE refresh_tokens SET status='inactive' WHERE content=$1;")
	if err != nil {
		log.Printf("InvalidateRefreshTokenByContent: error creating statement: %s", err.Error())
		tx.Rollback()
		return fmt.Errorf("InvalidateRefreshTokenByContent: error creating prepared statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(content)
	if err != nil {
		log.Printf("InvalidateRefreshTokenByContent: error executing query: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("InvalidateRefreshTokenByContent: error during commit: %s", err.Error())
		tx.Rollback()
		return err
	}

	log.Printf("InvalidateRefreshTokenByContent: invalidated refresh token")
	return nil
}
