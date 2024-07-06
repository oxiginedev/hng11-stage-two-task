package sqlstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/oxiginedev/hng11-stage-two-task/pkg/datastore"
	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
)

const (
	createOrganisationUser = `
	INSERT INTO organisation_user (id, organisation_id, user_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)
	`

	fetchOrgUserOrganisations = `
	SELECT o.* FROM organisation_user ou
	JOIN organisations o ON ou.organisation_id = o.id
	WHERE ou.user_id = ?
	ORDER BY o.name ASC
	`

	fetchOrgUserByUserID = `
	SELECT
		o.id AS id,
		o.organisation_is AS organisation_id,
		u.id AS user_id,
		u.first_name AS user_metadata.first_name,
		u.last_name AS user_metadata.last_name,
		u.email AS user_metadata.email,
		u.phone AS user_metadata.phone
	FROM organisation_user o
	LEFT JOIN users u
		ON o.user_id = u.id
	WHERE o.user_id = ?
	AND o.organisation_id = ?
	`
)

func (s *sqlstore) CreateOrganisationUser(ctx context.Context, ou *models.OrganisationUser) error {
	query := s.Rebind(createOrganisationUser)

	result, err := s.ExecContext(ctx, query,
		ou.ID,
		ou.OrganisationID,
		ou.UserID,
		ou.CreatedAt,
		ou.UpdatedAt,
	)
	if err != nil {
		return err
	}

	aRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if aRows < 1 {
		return datastore.ErrRecordNotCreated
	}

	return err
}

func (s *sqlstore) FetchUserOrganisations(ctx context.Context, userID string) ([]*models.Organisation, error) {
	query := s.Rebind(fetchOrgUserOrganisations)

	rows, err := s.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	organisations := make([]*models.Organisation, 0)
	for rows.Next() {
		var org *models.Organisation

		err = rows.StructScan(org)
		if err != nil {
			return nil, err
		}

		organisations = append(organisations, org)
	}

	return organisations, nil
}

func (s *sqlstore) FetchOrganisationUserByUserID(ctx context.Context, userID, orgID string) (*models.OrganisationUser, error) {
	orgUser := &models.OrganisationUser{}
	err := s.QueryRowxContext(ctx, s.Rebind(fetchOrgUserByUserID), userID, orgID).StructScan(orgUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datastore.ErrRecordNotFound
		}
		return nil, err
	}

	return orgUser, nil
}
