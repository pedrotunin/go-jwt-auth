package repositories

import (
	"database/sql"
	"log"

	"github.com/pedrotunin/go-jwt-auth/internal/models"
)

type PSQLAppRepository struct {
	db *sql.DB
}

func NewPSQLAppRepository(db *sql.DB) *PSQLAppRepository {
	return &PSQLAppRepository{
		db: db,
	}
}

func (repo *PSQLAppRepository) CreateApp(app *models.App) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("CreateApp: error creating transaction: %s", err.Error())
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO apps (name, description, user_id) VALUES ($1, $2, $3);")
	if err != nil {
		log.Printf("CreateApp: error creating statement: %s", err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(app.Name, app.Description, app.UserID)
	if err != nil {
		log.Printf("CreateApp: error executing query: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("CreateApp: error during commmit: %s", err.Error())
		tx.Rollback()
		return err
	}

	log.Printf("CreateApp: app created")
	return nil
}
