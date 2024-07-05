package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
)

func (s *sqlstore) CreateUser(ctx context.Context, u *models.User) error {
	var query = s.Rebind(
		"INSERT INTO users(first_name, last_name, email, phone, password) VALUES(?, ?, ?, ?, ?)")

	result, err := s.ExecContext(ctx, query,
		u.ID,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		u.Password,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return datastore.ErrDuplicate
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return datastore.ErrRecordNotCreated
	}

	return nil
}

func (s *sqlstore) FetchUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	query := s.Rebind("SELECT * FROM users WHERE id = ? LIMIT 1")

	err := s.QueryRowxContext(ctx, query, id).StructScan(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datastore.ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *sqlstore) FetchUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := s.Rebind("SELECT * FROM users WHERE email = ? LIMIT 1")

	err := s.QueryRowxContext(ctx, query, email).StructScan(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datastore.ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}
