package app

import (
	"github.com/google/uuid"
	"gitlab.com/erikwu09/yamlr/models"
)

type MetadataRepository interface {
	Insert(metadata *models.Metadata) (id uuid.UUID, err error)
	Update(id uuid.UUID, metadata *models.Metadata) (err error)
	Get(id uuid.UUID) (metadata *models.Metadata, err error)
	Search(metadata *models.Metadata) (results []models.Metadata, err error)
}

type MetadataValidator interface {
	ValidateAndSanitize(metadata *models.Metadata) error
	SanitizeURLs(metadata *models.Metadata) error
}
