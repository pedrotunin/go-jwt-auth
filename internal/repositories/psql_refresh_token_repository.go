package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pedrotunin/jwt-auth/internal/models"
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
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO refresh_tokens (content, status) VALUES ($1, $2);")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.Content, token.Status)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (repo *PSQLRefreshTokenRepository) GetRefreshTokenByContent(content models.RefreshTokenContent) (*models.RefreshToken, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("GetRefreshTokenByContent: error creating transaction: %w", err)
	}

	stmt, err := tx.Prepare("SELECT * FROM refresh_tokens WHERE content=$1;")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("GetRefreshTokenByContent: error creating prepared statement: %w", err)
	}
	defer stmt.Close()

	var resId int
	var resContent, resStatus string
	err = stmt.QueryRow(content).Scan(&resId, &resContent, &resStatus)
	if err != nil {
		tx.Rollback()

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}

		return nil, fmt.Errorf("GetUserByEmail: error scanning query result: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &models.RefreshToken{
		ID:      resId,
		Content: resContent,
		Status:  resStatus,
	}, nil
}

func (repo *PSQLRefreshTokenRepository) InvalidateRefreshTokenByContent(content models.RefreshTokenContent) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return fmt.Errorf("InvalidateRefreshTokenByContent: error creating transaction: %w", err)
	}

	stmt, err := tx.Prepare("UPDATE refresh_tokens SET status='inactive' WHERE content=$1;")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("InvalidateRefreshTokenByContent: error creating prepared statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(content)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
