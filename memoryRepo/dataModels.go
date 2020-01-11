package memoryrepo

import (
	"gitlab.com/erikwu09/yamlr/models"
)

//MetadataDAO is an internal access object class for memoryRepository
type MetadataDAO struct {
	models.Metadata
	Id            string
	MaintainerIDs []string
}

//MaintainerDAO is an internal access object class for memoryRepository
type MaintainerDAO struct {
	models.Maintainer
	Id          string
	MetadataIds []string
}
