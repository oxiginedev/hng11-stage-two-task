package datastore

import (
	"context"
	"errors"

	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
)

var (
	ErrDuplicate        = errors.New("duplicate key")
	ErrRecordNotFound   = errors.New("record not found")
	ErrRecordNotCreated = errors.New("record not created")
)

type Datastore interface {
	// users
	CreateUser(context.Context, *models.User) error
	FetchUserByID(context.Context, string) (*models.User, error)
	FetchUserByEmail(context.Context, string) (*models.User, error)

	// organisations
	CreateOrganisation(context.Context, *models.Organisation) error
	FetchOrganisationByID(context.Context, string) (*models.Organisation, error)

	// organisation_user
	CreateOrganisationUser(context.Context, *models.OrganisationUser) error
	FetchUserOrganisations(context.Context, string) ([]*models.Organisation, error)
	FetchOrganisationUserByUserID(ctx context.Context, userID, orgID string) (*models.OrganisationUser, error)

	AutoMigrate() error
	Close() error
}
