package models

import "time"

type Organisation struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"created_at"`
}

type OrganisationUser struct {
	ID             string       `db:"id"`
	OrganisationID string       `db:"organisation_id"`
	UserID         string       `db:"user_id"`
	UserMetadata   UserMetadata `db:"user_metadata"`
	CreatedAt      time.Time    `db:"created_at"`
	UpdatedAt      time.Time    `db:"created_at"`
}

type UserMetadata struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Phone     string `db:"phone"`
}
