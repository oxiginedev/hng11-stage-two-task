package sqlstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
)

const (
	createOrganisation = `
	INSERT INTO organisations(id, name, description, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)
	`
)

func (s *sqlstore) CreateOrganisation(ctx context.Context, o *models.Organisation) error {
	query := s.Rebind(createOrganisation)

	result, err := s.ExecContext(ctx, query,
		o.ID,
		o.Name,
		o.Description,
		o.CreatedAt,
		o.UpdatedAt,
	)
	if err != nil {
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

func (s *sqlstore) FetchOrganisationByID(ctx context.Context, id string) (*models.Organisation, error) {
	org := &models.Organisation{}
	query := s.Rebind("SELECT * FROM organisations WHERE id = ? LIMIT 1")

	err := s.QueryRowxContext(ctx, query, id).StructScan(org)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datastore.ErrRecordNotFound
		}
		return nil, err
	}

	return org, nil
}
