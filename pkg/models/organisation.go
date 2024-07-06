package models

import "time"

type Organisation struct {
	ID          string    `db:"id" json:"orgId"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
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
