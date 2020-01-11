package memoryrepo

import (
	"errors"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	memdb "github.com/hashicorp/go-memdb"
	"gitlab.com/erikwu09/yamlr/models"
	"gitlab.com/erikwu09/yamlr/utils"
)

type memoryRepository struct {
	//memoryDB is a in-memory schema-based db that supports data indexing for optimized querying
	memoryDb *memdb.MemDB
}

const metadataTable string = "metadata"
const maintainerTable string = "maintainers"

var repo *memoryRepository

var schema *memdb.DBSchema

// Returns a memoryRepository
func GetMemoryRepository() (*memoryRepository, error) {

	if repo != nil {
		return repo, nil
	}

	var table = make(map[string]*memdb.TableSchema)

	// set up search indexes for metadata table
	var metadataIndexes = make(map[string]*memdb.IndexSchema)

	metadataIndexes["id"] = &memdb.IndexSchema{
		Name:    "id",
		Unique:  true,
		Indexer: &memdb.StringFieldIndex{Field: "Id"},
	}

	metadataIndexes["Title"] = &memdb.IndexSchema{
		Name:    "Title",
		Unique:  false,
		Indexer: &memdb.StringFieldIndex{Field: "Title"},
	}

	metadataIndexes["Company"] = &memdb.IndexSchema{
		Name:    "Company",
		Unique:  false,
		Indexer: &memdb.StringFieldIndex{Field: "Company"},
	}

	metadataIndexes["Website"] = &memdb.IndexSchema{
		Name:    "Website",
		Unique:  false,
		Indexer: &memdb.StringFieldIndex{Field: "Website", Lowercase: true},
	}

	metadataIndexes["Source"] = &memdb.IndexSchema{
		Name:    "Source",
		Unique:  false,
		Indexer: &memdb.StringFieldIndex{Field: "Source", Lowercase: true},
	}

	metadataIndexes["License"] = &memdb.IndexSchema{
		Name:    "License",
		Unique:  false,
		Indexer: &memdb.StringFieldIndex{Field: "License"},
	}

	metadataIndexes["MaintainerIDs"] = &memdb.IndexSchema{
		Name:    "MaintainerIDs",
		Unique:  false,
		Indexer: &memdb.StringSliceFieldIndex{Field: "MaintainerIDs"},
	}

	table[metadataTable] = &memdb.TableSchema{
		Name:    metadataTable,
		Indexes: metadataIndexes,
	}

	//set up search indexes for maintainers
	var maintainerIndexes = make(map[string]*memdb.IndexSchema)

	maintainerIndexes["id"] = &memdb.IndexSchema{
		Name:    "id",
		Unique:  true,
		Indexer: &memdb.StringFieldIndex{Field: "Id"},
	}

	maintainerIndexes["Name"] = &memdb.IndexSchema{
		Name:    "Name",
		Unique:  false,
		Indexer: &memdb.StringFieldIndex{Field: "Name"},
	}

	maintainerIndexes["Email"] = &memdb.IndexSchema{
		Name: "Email",
		// email should be unique index key?
		Unique:  true,
		Indexer: &memdb.StringFieldIndex{Field: "Email", Lowercase: true},
	}

	table[maintainerTable] = &memdb.TableSchema{
		Name:    maintainerTable,
		Indexes: maintainerIndexes,
	}

	schema = &memdb.DBSchema{Tables: table}

	cache, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	repo = &memoryRepository{memoryDb: cache}

	return repo, nil
}

func (r *memoryRepository) Get(id uuid.UUID) (metadata *models.Metadata, err error) {
	tx := r.memoryDb.Txn(false)
	defer abortOrCommit(err, tx)

	result, err := r.getMetadata(id, tx)
	if err != nil {
		return nil, err
	}
	return &result.Metadata, nil
}

func (r *memoryRepository) Insert(metadata *models.Metadata) (id uuid.UUID, err error) {
	tx := r.memoryDb.Txn(true)
	defer abortOrCommit(err, tx)
	//first get maintainers
	id = uuid.New()
	_, err = r.insertMetadata(metadata, tx, id)
	return id, err
}

func (r *memoryRepository) Update(id uuid.UUID, metadata *models.Metadata) (err error) {
	tx := r.memoryDb.Txn(true)
	defer abortOrCommit(err, tx)

	result, err := r.getMetadata(id, tx)
	if err != nil {
		return
	}
	if result == nil {
		return
	}
	// we can also just change the reference of the returned pointer to the new model
	tx.Delete(metadataTable, result)

	_, err = r.insertMetadata(metadata, tx, id)
	return
}

func (r *memoryRepository) Search(metadata *models.Metadata) (results []models.Metadata, err error) {
	params := getQueryParams(*metadata)
	temp, err := r.queryByParams(params)
	if err != nil {
		return nil, err
	}
	results = make([]models.Metadata, 0, len(temp))
	for _, result := range temp {
		results = append(results, *result)
	}
	return
}

func getQueryParams(metadata models.Metadata) map[string]interface{} {
	params := structs.Map(metadata)
	if metadata.Maintainers != nil && len(metadata.Maintainers) > 0 {
		emails := make([]string, 0, len(metadata.Maintainers))
		for _, m := range metadata.Maintainers {
			emails = append(emails, m.Email)
		}
		params["Email"] = emails
	}
	if _, k := params["Maintainers"]; k {
		delete(params, "Maintainers")
	}
	return params
}

// TODO: refactor this
func (r *memoryRepository) queryByParams(params map[string]interface{}) (results []*models.Metadata, err error) {
	var resultSet *utils.Set
	tx := r.memoryDb.Txn(true)
	defer abortOrCommit(err, tx)
	for index, queryVal := range params {
		if isQueryValEmptyOrNull(queryVal) {
			continue
		}
		var temp *utils.Set
		if index == "Email" {
			temp, err = r.queryByEmails(queryVal.([]string), tx)
		} else {
			temp, err = r.query(index, queryVal, tx)
		}
		if err != nil {
			continue
			//return nil, err
		}
		if resultSet == nil {
			resultSet = temp
		} else {
			resultSet = resultSet.Intersect(*temp)
		}
	}
	if resultSet == nil {
		return
	}
	for result := range *resultSet {
		results = append(results, result.(*models.Metadata))
	}
	return
}

func (r *memoryRepository) query(index string, queryVal interface{}, tx *memdb.Txn) (*utils.Set, error) {
	result := utils.NewSet()
	it, err := tx.Get(metadataTable, index, queryVal)
	if err != nil {
		return nil, err
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		metadata, k := obj.(*MetadataDAO)
		if !k {
			return nil, errors.New("something went wrong")
		}
		result.Add(&(metadata.Metadata))
	}
	return result, nil
}

// TODO: use reflection instead
func isQueryValEmptyOrNull(queryVal interface{}) bool {
	switch queryVal.(type) {
	case string:
		return queryVal.(string) == ""
	// assume it's a slice otherwise
	default:
		slice, k := queryVal.([]interface{})
		if !k {
			return false
		}
		return slice == nil
	}
}

func (r *memoryRepository) queryByEmails(emails []string, tx *memdb.Txn) (results *utils.Set, err error) {
	//resultSet here is a set of metadata IDs
	var metadataIds *utils.Set

	for _, email := range emails {
		temp, err := r.queryByEmail(email, tx)
		if err != nil {
			return nil, err
		}
		if metadataIds == nil {
			metadataIds = temp
		} else {
			metadataIds = metadataIds.Intersect(*temp)
		}
	}
	results = utils.NewSet()
	for metadataID := range *metadataIds {
		idString, _ := metadataID.(string)
		id, err := uuid.Parse(idString)
		if err != nil {
			return nil, err
		}
		metadata, err := r.getMetadata(id, tx)
		if err != nil {
			return nil, err
		}
		results.Add(&(metadata.Metadata))
	}
	return
}

func (r *memoryRepository) queryByEmail(email string, tx *memdb.Txn) (*utils.Set, error) {
	obj, err := tx.First(maintainerTable, "Email", email)
	if err != nil {
		return nil, err
	}
	maintainer, k := obj.(*MaintainerDAO)
	if !k {
		return nil, errors.New("something went wrong")
	}
	results := utils.NewSet()
	for _, metadataID := range maintainer.MetadataIds {
		results.Add(metadataID)
	}

	return results, nil
}

func (r *memoryRepository) getMetadata(id uuid.UUID, tx *memdb.Txn) (metadataDAO *MetadataDAO, err error) {
	var res interface{}
	res, err = tx.First(metadataTable, "id", id.String())
	if err != nil {
		return
	}
	if res == nil {
		return nil, errors.New("no metadata found")
	}
	metadataDAO, k := res.(*MetadataDAO)
	if !k {
		return
	}
	return
}

func (r *memoryRepository) addNewMaintainer(maintainer *models.Maintainer, metadataID *uuid.UUID, tx *memdb.Txn) (*MaintainerDAO, error) {
	maintainerID := uuid.New()
	maintainerDAO := &MaintainerDAO{Maintainer: *maintainer, Id: maintainerID.String()}
	if metadataID != nil {
		maintainerDAO.MetadataIds = append(maintainerDAO.MetadataIds, metadataID.String())
	}
	err := tx.Insert(maintainerTable, maintainerDAO)
	if err != nil {
		return nil, err
	}
	return maintainerDAO, nil
}

func (r *memoryRepository) insertMetadata(metadata *models.Metadata, tx *memdb.Txn, id uuid.UUID) (metadataDAO *MetadataDAO, err error) {
	//first get maintainers
	maintainerIds := make([]string, 0)
	for _, m := range metadata.Maintainers {
		var res interface{}
		res, err = tx.First(maintainerTable, "Email", m.Email)
		if err != nil {
			return
		}
		var maintainerID string
		if res == nil {
			var maintainerDAO *MaintainerDAO
			maintainerDAO, err = r.addNewMaintainer(m, &id, tx)
			if err != nil {
				return
			}
			maintainerID = maintainerDAO.Id
		} else {
			maintainer, k := res.(*MaintainerDAO)
			if !k {
				return nil, errors.New("something went wrong")
			}
			maintainer.MetadataIds = append(maintainer.MetadataIds, id.String())
			maintainerID = maintainer.Id
		}
		maintainerIds = append(maintainerIds, maintainerID)
	}

	//persist
	metadataDAO = &MetadataDAO{Metadata: *metadata, Id: id.String(), MaintainerIDs: maintainerIds}

	err = tx.Insert(metadataTable, metadataDAO)
	if err != nil {
		return
	}
	return
}

func abortOrCommit(err error, tx *memdb.Txn) {
	if err != nil {
		tx.Abort()
	} else {
		tx.Commit()
	}
}
