package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	_ "github.com/lib/pq"
)

type dbProfile struct {
	UserID    string
	Name      sql.NullString
	Phone     sql.NullString
	Addresses pq.StringArray
	Points    sql.NullInt64
	Avatar    sql.NullString
}

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
		"SELECT userID, username, phone, addresses, points, avatar FROM profiles WHERE userID=$1",
		userID,
	)

	profile := dbProfile{}
	switch err := row.Scan(&profile.UserID, &profile.Name, &profile.Phone, &profile.Addresses, &profile.Points, &profile.Avatar); err {
	case sql.ErrNoRows:
		return models.Profile{}, fmt.Errorf("profile doesn't exist")
	case nil:
		return fromDBToApp(profile), nil
	default:
		return models.Profile{}, fmt.Errorf("error with db: %w", err)
	}
}

func (p ProfileRepository) Update(appProfile models.Profile) error {
	profile := fromAppToDB(appProfile)
	_, err := p.db.Exec(
		`UPDATE profiles
		SET username=$1, phone=$2, addresses=$3, points=$4
		WHERE userID=$5`,
		profile.Name, profile.Phone, profile.Addresses, profile.Points, profile.UserID,
	)

	if err != nil {
		return fmt.Errorf("couln't update profile: %w", err)
	}

	return nil
}

func (p ProfileRepository) UpdateAvatar(userID string, link string) error {
	_, err := p.db.Exec(
		`UPDATE profiles
		SET avatar=$1
		WHERE userID=$2`,
		link, userID,
	)

	if err != nil {
		return fmt.Errorf("couln't update avatar: %w", err)
	}

	return nil
}

func (p ProfileRepository) UpdateAddresses(userID string, addresses []string) error {
	_, err := p.db.Exec(
		"UPDATE profiles SET addresses = $1 WHERE userID = $2",
		pq.StringArray(addresses), userID,
	)

	if err != nil {
		return fmt.Errorf("couldn't update addresses: %w", err)
	}

	return nil
}

func (p ProfileRepository) Delete(userID string) error {
	_, err := p.db.Exec(
		"DELETE FROM profiles WHERE userID=$1",
		userID,
	)

	if err != nil {
		return fmt.Errorf("couldn't delete profile from Postgres: %w", err)
	}

	return nil
}

func fromDBToApp(dbProf dbProfile) models.Profile {
	appProf := models.Profile{
		UserID: dbProf.UserID,
		Name:   dbProf.Name.String,
		Phone:  dbProf.Phone.String,
		Points: int(dbProf.Points.Int64),
		Avatar: dbProf.Avatar.String,
	}

	for _, addr := range dbProf.Addresses {
		appProf.Addresses = append(appProf.Addresses, addr)
	}

	return appProf
}

func fromAppToDB(appProf models.Profile) dbProfile {
	return dbProfile{
		UserID:    appProf.UserID,
		Name:      sql.NullString{String: appProf.Name, Valid: true},
		Phone:     sql.NullString{String: appProf.Phone, Valid: true},
		Points:    sql.NullInt64{Int64: int64(appProf.Points), Valid: true},
		Addresses: pq.StringArray(appProf.Addresses),
		Avatar:    sql.NullString{String: appProf.Avatar, Valid: true},
	}
}
