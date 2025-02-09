package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pedrotunin/jwt-auth/internal/models"
)

type PSQLUserRepository struct {
	db *sql.DB
}

func NewPSQLUserRepository(db *sql.DB) *PSQLUserRepository {
	return &PSQLUserRepository{
		db: db,
	}
}

func (repo *PSQLUserRepository) GetUserByID(id models.UserID) (*models.User, error) {
	return nil, nil
}

func (repo *PSQLUserRepository) GetUserByEmail(email models.UserEmail) (*models.User, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail: error creating transaction: %w", err)
	}

	stmt, err := tx.Prepare("SELECT * FROM users WHERE email=$1;")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("GetUserByEmail: error creating prepared statement: %w", err)
	}
	defer stmt.Close()

	var resId int
	var resEmail, resPassword string
	err = stmt.QueryRow(email).Scan(&resId, &resEmail, &resPassword)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("GetUserByEmail: error scanning query result: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &models.User{
		ID:       resId,
		Email:    resEmail,
		Password: resPassword,
	}, nil
}

func (repo *PSQLUserRepository) CreateUser(u *models.User) (id int, err error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return -1, err
	}

	stmt, err := tx.Prepare("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id;")
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	defer stmt.Close()

	var insertedID int
	err = stmt.QueryRow(u.Email, u.Password).Scan(&insertedID)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return insertedID, nil
}

func (repo *PSQLUserRepository) UpdateUser(u *models.User) error {
	return nil
}

func (repo *PSQLUserRepository) DeleteUser(u *models.User) error {
	return nil
}

func (repo *PSQLUserRepository) DeleteUserByID(id models.UserID) error {
	return nil
}
