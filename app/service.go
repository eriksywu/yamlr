package app

import (
	"errors"
	"github.com/google/uuid"
	"gitlab.com/erikwu09/yamlr/models"
	"log"
)

type MetadataManager struct {
	repository MetadataRepository
	validator  MetadataValidator
	logger     *log.Logger
}

func (m MetadataManager) CreateMetadata(metadata models.Metadata) (uuid.UUID, error) {
	var id uuid.UUID
	if err := m.validator.ValidateAndSanitize(&metadata); err != nil {
		return id, err
	}
	id, err := m.repository.Insert(&metadata)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (m MetadataManager) UpdateMetadata(metadata models.Metadata, id uuid.UUID) error {
	if err := m.validator.ValidateAndSanitize(&metadata); err != nil {
		return err
	}
	if err := m.repository.Update(id, &metadata); err != nil {
		return err
	}
	return nil
}

func (m MetadataManager) SearchMetadata(metadata models.Metadata) (models.SearchResults, error) {
	if err := m.validator.SanitizeURLs(&metadata); err != nil {
		return models.SearchResults{}, err
	}
	results, err := m.repository.Search(&metadata)
	if err != nil {
		return models.SearchResults{}, err
	}
	searchResults := models.SearchResults{Results: results}
	return searchResults, nil
}

func (m MetadataManager) GetMetadata(id uuid.UUID) (*models.Metadata, error) {
	return m.repository.Get(id)
}

func BuildMetadataManager(repository MetadataRepository, validator MetadataValidator, logger *log.Logger) (MetadataManager, error) {
	var mgr MetadataManager
	if repository == nil {
		return mgr, errors.New("repo not set")
	}
	if validator == nil {
		return mgr, errors.New("validator not set")
	}

	return MetadataManager{repository: repository, validator: validator, logger: logger}, nil
}
