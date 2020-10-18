package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) profile.Repository {
	return ProfileRepository{
		db: db,
	}
}

func (p ProfileRepository) Create(userID string) error {
	_, err := p.db.Exec(
		"INSERT INTO profiles (userID) VALUES ($1)",
		userID,
	)

	if err != nil {
		return fmt.Errorf("couldn't create profile: %w", err)
	}

	return nil
}

func (p ProfileRepository) Get(userID string) (models.Profile, error) {
	row := p.db.QueryRow(
		"SELECT userID, name, phone, addresses, points FROM profiles WHERE userID=$1",
		userID,
	)

	profile := models.Profile{}
	switch err := row.Scan(&profile.UserID, &profile.Name, &profile.Phone, &profile.Addresses, &profile.Points); err {
	case sql.ErrNoRows:
		return models.Profile{}, fmt.Errorf("profile doesn't exist")
	case nil:
		return profile, nil
	default:
		return models.Profile{}, fmt.Errorf("error with db: %w", err)
	}
}

func (p ProfileRepository) Update(profile models.Profile) error {
	_, err := p.db.Exec(
		`UPDATE profiles
		SET name=$1, phone=$2, addresses=$3, points=$4
		WHERE userID=$5`,
		profile.Name, profile.Phone, profile.Addresses, profile.Points, profile.UserID,
	)

	if err != nil {
		return fmt.Errorf("couln't update profile: %w", err)
	}

	return nil
}

func (p ProfileRepository) Delete(userID string) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("couldn't convert to string: %w", err)
	}

	_, err = p.db.Exec(
		"DELETE FROM profiles WHERE userID=$1",
		id,
	)

	if err != nil {
		return fmt.Errorf("couldn't delete profile from Postgres: %w", err)
	}

	return nil
}
