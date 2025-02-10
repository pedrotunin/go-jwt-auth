package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type PSQLUserRepository struct {
	db *sql.DB
}

func NewPSQLUserRepository(db *sql.DB) *PSQLUserRepository {
	return &PSQLUserRepository{
		db: db,
	}
}

func (repo *PSQLUserRepository) GetUserByEmail(email models.UserEmail) (*models.User, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("GetUserByEmail: error creating transaction: %s", err.Error())
		return nil, fmt.Errorf("GetUserByEmail: error creating transaction: %w", err)
	}

	stmt, err := tx.Prepare("SELECT * FROM users WHERE email=$1;")
	if err != nil {
		log.Printf("GetUserByEmail: error creating statement: %s", err.Error())

		tx.Rollback()
		return nil, fmt.Errorf("GetUserByEmail: error creating prepared statement: %w", err)
	}
	defer stmt.Close()

	var resId int
	var resEmail, resPassword string
	err = stmt.QueryRow(email).Scan(&resId, &resEmail, &resPassword)
	if err != nil {
		log.Printf("GetUserByEmail: error executing query: %s", err.Error())

		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrUserNotFound
		}

		return nil, fmt.Errorf("GetUserByEmail: error scanning query result: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("GetUserByEmail: error during commit: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	log.Print("GetUserByEmail: user found in users table")
	return &models.User{
		ID:       resId,
		Email:    resEmail,
		Password: resPassword,
	}, nil
}

func (repo *PSQLUserRepository) CreateUser(u *models.User) (id int, err error) {
	user, _ := repo.GetUserByEmail(u.Email)
	if user != nil {
		log.Printf("GetUserByEmail: user email found in database")
		return -1, utils.ErrUserEmailAlreadyExists
	}

	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("CreateUser: error starting transaction: %s", err.Error())
		return -1, err
	}

	stmt, err := tx.Prepare("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id;")
	if err != nil {
		log.Printf("CreateUser: error starting statement: %s", err.Error())
		tx.Rollback()
		return -1, err
	}
	defer stmt.Close()

	var insertedID int
	err = stmt.QueryRow(u.Email, u.Password).Scan(&insertedID)
	if err != nil {
		log.Printf("CreateUser: error executing query: %s", err.Error())
		tx.Rollback()
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("CreateUser: error during commit: %s", err.Error())
		tx.Rollback()
		return -1, err
	}

	log.Print("CreateUser: user created in users table")
	return insertedID, nil
}
