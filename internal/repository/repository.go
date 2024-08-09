package repository

import (
	"database/sql"

	"github.com/Vykiy/house-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(email, passwordHash string, userType models.UserType) (uuid.UUID, error) {
	var userID uuid.UUID
	if err := r.db.QueryRow("INSERT INTO users (email, password_hash, user_type) VALUES ($1, $2, $3) RETURNING id", email, passwordHash, userType).Scan(&userID); err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (r *Repository) GetUser(userID uuid.UUID) (models.User, error) {
	var user models.User
	if err := r.db.Get(&user, "SELECT id, email, password_hash, user_type FROM users WHERE id = $1", userID); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *Repository) CreateHouse(address, developer string, yearBuilt int) (models.House, error) {
	var house models.House
	if err := r.db.QueryRow("INSERT INTO houses (address, developer, year_built) VALUES ($1, $2, $3) RETURNING address, developer, year_built, id, created_at, updated_at",
		address, developer, yearBuilt).Scan(&house.Address, &house.Developer, &house.YearBuilt, &house.ID, &house.CreatedAt, &house.UpdatedAt); err != nil {
		return models.House{}, err
	}

	return house, nil
}

func (r *Repository) GetFlats(houseID int) ([]models.Flat, error) {
	var flats []models.Flat
	if err := r.db.Select(&flats, "SELECT flat_number, house_id, price, rooms FROM flats WHERE house_id = $1", houseID); err != nil {
		return nil, err
	}

	return flats, nil

}

func (r *Repository) CreateFlat(houseID int, price, rooms int) (models.Flat, error) {
	var flat models.Flat

	tx, err := r.db.Begin()
	if err != nil {
		return models.Flat{}, err
	}

	var lastFlatNumber int
	if err := tx.QueryRow("SELECT flat_number FROM flats WHERE house_id = $1 ORDER BY flat_number DESC LIMIT 1", houseID).Scan(&lastFlatNumber); err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return models.Flat{}, err
	}

	if err := tx.QueryRow("INSERT INTO flats (house_id, price, rooms, flat_number, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, house_id, price, rooms",
		houseID, price, rooms, lastFlatNumber+1, models.FlatStatusCreated).Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms); err != nil {
		tx.Rollback()
		return models.Flat{}, err
	}

	if _, err := tx.Exec("UPDATE houses SET updated_at = NOW() WHERE id = $1", flat.HouseID); err != nil {
		tx.Rollback()
		return models.Flat{}, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return models.Flat{}, err
	}

	return flat, nil
}

func (r *Repository) UpdateFlat(flatID int, status models.FlatStatus) (models.Flat, error) {
	var flat models.Flat

	if err := r.db.QueryRow("UPDATE flats SET status = $1 WHERE id = $2 RETURNING flat_number, house_id, price, rooms, status",
		status, flatID).Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
		return models.Flat{}, err
	}

	return flat, nil
}

func (r *Repository) GetFlatModerator(flatID int) (uuid.UUID, error) {
	var moderatorID uuid.UUID
	if err := r.db.QueryRow("SELECT moderator_id FROM flats WHERE id = $1", flatID).Scan(&moderatorID); err != nil {
		return uuid.Nil, err
	}

	return moderatorID, nil
}

func (r *Repository) SubscribeToNewFlats(houseID int, email string) error {
	if _, err := r.db.Exec("INSERT INTO subscriptions (house_id, email) VALUES ($1, $2)", houseID, email); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetSubscribers(houseID int) ([]string, error) {
	var subscribers []string
	if err := r.db.Select(&subscribers, "SELECT email FROM subscriptions WHERE house_id = $1", houseID); err != nil {
		return nil, err
	}

	return subscribers, nil
}
