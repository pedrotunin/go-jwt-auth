package repositories

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type PSQLAppRepository struct {
	db *sql.DB
}

func NewPSQLAppRepository(db *sql.DB) *PSQLAppRepository {
	return &PSQLAppRepository{
		db: db,
	}
}

func (repo *PSQLAppRepository) GetAppsByUserID(userID models.UserID) ([]models.App, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("GetAppsByUserID: error creating transaction: %s", err.Error())
		return nil, err
	}

	query := "SELECT id, name, description, user_id, created_at, updated_at FROM apps WHERE user_id=$1 AND deleted_at IS NULL;"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Printf("GetAppsByUserID: error creating statement: %s", err.Error())
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		log.Printf("GetAppsByUserID: error executing query: %s", err.Error())
		tx.Rollback()

		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrAppNotFound
		}

		return nil, err
	}
	defer rows.Close()

	apps := []models.App{}

	for rows.Next() {
		var id, userId int
		var name, description string
		var createdAt, updatedAt time.Time

		err := rows.Scan(&id, &name, &description, &userId, &createdAt, &updatedAt)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		apps = append(apps, models.App{
			ID:          id,
			Name:        name,
			Description: description,
			UserID:      userId,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})

	}

	err = rows.Err()
	if err != nil {
		log.Printf("GetAppsByUserID: error during iterating rows: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("GetAppsByUserID: error during commmit: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	log.Printf("GetAppsByUserID: got apps")
	return apps, nil
}

func (repo *PSQLAppRepository) GetAppByID(appID models.AppID) (*models.App, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("GetAppByID: error creating transaction: %s", err.Error())
		return nil, err
	}

	query := "SELECT id, name, description, user_id, created_at, updated_at FROM apps WHERE id=$1 AND deleted_at IS NULL;"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Printf("GetAppByID: error creating statement: %s", err.Error())
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	var id, userId int
	var name, description string
	var createdAt, updatedAt time.Time
	err = stmt.QueryRow(appID).Scan(&id, &name, &description, &userId, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("GetAppByID: error executing query: %s", err.Error())
		tx.Rollback()

		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrAppNotFound
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("GetAppByID: error during commmit: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	app := models.App{
		ID:          id,
		Name:        name,
		Description: description,
		UserID:      userId,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	log.Printf("GetAppByID: got app")
	return &app, nil

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

func (repo *PSQLAppRepository) UpdateApp(app *models.App) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("UpdateApp: error creating transaction: %s", err.Error())
		return err
	}

	stmt, err := tx.Prepare("UPDATE apps SET name=$1, description=$2, updated_at=$3 WHERE id=$4;")
	if err != nil {
		log.Printf("UpdateApp: error creating statement: %s", err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(app.Name, app.Description, time.Now(), app.ID)
	if err != nil {
		log.Printf("UpdateApp: error executing query: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("UpdateApp: error during commmit: %s", err.Error())
		tx.Rollback()
		return err
	}

	log.Printf("UpdateApp: app updated")
	return nil
}

func (repo *PSQLAppRepository) DeleteAppByID(appID models.AppID) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Printf("DeleteAppByID: error creating transaction: %s", err.Error())
		return err
	}

	stmt, err := tx.Prepare("UPDATE apps SET deleted_at=$1 WHERE id=$2;")
	if err != nil {
		log.Printf("DeleteAppByID: error creating statement: %s", err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), appID)
	if err != nil {
		log.Printf("DeleteAppByID: error executing query: %s", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("DeleteAppByID: error during commmit: %s", err.Error())
		tx.Rollback()
		return err
	}

	log.Printf("DeleteAppByID: app deleted_at set deleted")
	return nil

}
