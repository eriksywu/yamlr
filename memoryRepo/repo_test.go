package memoryrepo

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/erikwu09/yamlr/models"
	"strconv"
	"testing"
)

func Test_Setup(t *testing.T) {
	_, err := GetMemoryRepository()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func Test_InsertAndGet(t *testing.T) {
	repo, _ := GetMemoryRepository()
	metadata := dummyMetadata("", "", "")
	id, err := repo.Insert(metadata)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	assert.NotNil(t, id)
	result, err := repo.Get(id)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	assert.Equal(t, metadata, result)
}

func Test_InsertAndUpdate(t *testing.T) {
	repo, _ := GetMemoryRepository()
	metadata := dummyMetadata("", "", "")
	id, err := repo.Insert(metadata)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	assert.NotNil(t, id)
	newMetadata := dummyMetadata("new", "", "")
	err = repo.Update(id, newMetadata)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	result, err := repo.Get(id)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	assert.Equal(t, newMetadata.Title, result.Title)
	assert.Equal(t, newMetadata, result)
}

func Test_InsertMultiple(t *testing.T) {
	repo, _ := GetMemoryRepository()
	metadata := dummyMetadata("", "", "")
	id1, err := repo.Insert(metadata)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	assert.NotNil(t, id1)
	newMetadata := dummyMetadata("new", "", "")
	id2, err := repo.Insert(newMetadata)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	assert.NotNil(t, id2)
}

func Test_HandleNonExistingMetadata(t *testing.T) {
	repo, _ := GetMemoryRepository()
	result, err := repo.Get(uuid.New())
	assert.Error(t, err)
	assert.Nil(t, result)
}

func Test_QueryMetadata(t *testing.T) {
	repo, _ := GetMemoryRepository()
	for i := 0; i < 5; i++ {
		metadata := dummyMetadata(strconv.Itoa(i), "", "")
		repo.Insert(metadata)
	}
	queryParams := map[string]interface{}{"Title": "dummyTitle1"}
	results, err := repo.queryByParams(queryParams)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.True(t, len(results) == 1)

	queryParams = map[string]interface{}{"Title": "Does not exist"}
	results, err = repo.queryByParams(queryParams)
	assert.NoError(t, err)
	assert.Nil(t, results)

	queryParams = map[string]interface{}{"Does not eixst": "Does not exist"}
	results, err = repo.queryByParams(queryParams)
	assert.Error(t, err)
}

func Test_QueryMultipleFields(t *testing.T) {
	repo, _ := GetMemoryRepository()
	repo.Insert(dummyMetadata("", "", ""))
	repo.Insert(dummyMetadata("titleSuffix", "companySuffix", ""))
	repo.Insert(dummyMetadata("titleSuffix2", "companySuffix", ""))
	repo.Insert(dummyMetadata("titleSuffix", "", ""))
	repo.Insert(dummyMetadata("titleSuffix2", "companySuffix", ""))
	repo.Insert(dummyMetadata("titleSuffix2", "companySuffix2", ""))

	queryParams := map[string]interface{}{"Title": "dummyTitle" + "titleSuffix2",
		"Company": "dummyCompany" + "companySuffix"}
	results, err := repo.queryByParams(queryParams)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.True(t, len(results) == 2)
}

func Test_QueryByEmails(t *testing.T) {
	repo, _ := GetMemoryRepository()
	expectedResult := dummyMetadata("", "", "microsoft.com")
	repo.Insert(expectedResult)
	repo.Insert(dummyMetadata("", "", "amazon.com"))
	repo.Insert(dummyMetadata("", "", "google.com"))
	repo.Insert(dummyMetadata("", "", "apple.com"))
	repo.Insert(dummyMetadata("", "", "microsoft.com"))

	tx := repo.memoryDb.Txn(false)
	defer tx.Abort()
	results, err := repo.queryByEmails([]string{"email@microsoft.com"}, tx)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	resultList := results.ToSlice()
	assert.True(t, len(resultList) == 2)
	assert.Equal(t, expectedResult, resultList[0])
}

func Test_QueryByMultipleEmails(t *testing.T) {
	repo, _ := GetMemoryRepository()
	expectedResult := dummyMetadata("", "", "microsoft.com")
	expectedResult.Maintainers = append(expectedResult.Maintainers, &models.Maintainer{Name: "erik wu", Email: "erw@microsoft.com"})
	repo.Insert(expectedResult)
	repo.Insert(dummyMetadata("", "", "amazon.com"))
	repo.Insert(dummyMetadata("", "", "google.com"))
	repo.Insert(dummyMetadata("", "", "apple.com"))
	repo.Insert(dummyMetadata("", "", "microsoft.com"))

	tx := repo.memoryDb.Txn(false)
	defer tx.Abort()
	results, err := repo.queryByEmails([]string{"email@microsoft.com", "erw@microsoft.com"}, tx)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	resultList := results.ToSlice()
	assert.True(t, len(resultList) == 1)
}

func Test_QueryByMetadataAndMaintainer(t *testing.T) {
	repo, _ := GetMemoryRepository()
	repo.Insert(dummyMetadata("", "", ""))
	//should return this
	expectedResult1 := dummyMetadata("expectedResult1", "microsoft", "microsoft.com")
	repo.Insert(expectedResult1)
	repo.Insert(dummyMetadata("titleSuffix2", "companySuffix", "aol.com"))
	repo.Insert(dummyMetadata("titleSuffix", "", "amazon.com"))
	//should return this
	expectedResult2 := dummyMetadata("expectedResult2", "microsoft", "microsoft.com")
	repo.Insert(expectedResult2)
	repo.Insert(dummyMetadata("titleSuffix2", "apple", "apple.com"))
	queryParams := map[string]interface{}{"Company": "dummyCompany" + "microsoft",
		"Email": []string{"email@microsoft.com"}}

	results, err := repo.queryByParams(queryParams)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.True(t, len(results) == 2)
}

func Test_ParamTranslation(t *testing.T) {
	data := models.Metadata{Title: "title",
		Version:     "1.0",
		Company:     "microsoft",
		Website:     "microsoft.com",
		License:     "gpl",
		Description: "adsfafa",
		Maintainers: []*models.Maintainer{&models.Maintainer{Name: "Erik Wu", Email: "erw@microsoft.com"}},
	}
	params := getQueryParams(data)
	assert.True(t, (params["Title"]).(string) == data.Title)
	assert.True(t, (params["Version"]).(string) == data.Version)
	assert.True(t, (params["Company"]).(string) == data.Company)
	assert.True(t, (params["License"]).(string) == data.License)
	assert.True(t, (params["Description"]).(string) == data.Description)
	assert.True(t, len((params["Email"]).([]string)) == 1)
	assert.True(t, (params["Email"]).([]string)[0] == data.Maintainers[0].Email)
}

func dummyMetadata(titleSuffix string, companySuffix string, emailDomain string) *models.Metadata {

	maintainers := []*models.Maintainer{&models.Maintainer{Name: "Tester", Email: "email@" + emailDomain}}
	return &models.Metadata{
		Title:       "dummyTitle" + titleSuffix,
		Version:     "1.0",
		Maintainers: maintainers,
		Company:     "dummyCompany" + companySuffix,
		Website:     "dummyCompany.com",
		Source:      "github.com/dummyCompany",
		License:     "GNU",
		Description: "Lorem ipsum",
	}
}

func xor(a, b bool) bool {
	return (a || b) && !(a && b)
}
