package validation

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/erikwu09/yamlr/app"
	"gitlab.com/erikwu09/yamlr/models"
)

func Test_ValidatesCorrectData(t *testing.T) {
	maintainers := []*models.Maintainer{&models.Maintainer{Name: "Tester", Email: "erikwu@microsoft.com"}}
	metadata := &models.Metadata{
		Title:       "dummyTitle",
		Version:     "1.0",
		Maintainers: maintainers,
		Company:     "dummyCompany",
		Website:     "http://www.dummyCompany.com",
		Source:      "github.com/dummyCompany",
		License:     "GNU",
		Description: "Lorem ipsum",
	}
	testee := SimpleValidator{}
	err := testee.ValidateAndSanitize(metadata)
	assert.NoError(t, err)
	log.Println(metadata.Website)
	log.Println(metadata.Source)
}

func Test_InvalidatesBadData(t *testing.T) {
	maintainers := []*models.Maintainer{&models.Maintainer{Name: "Tester", Email: "emicrosoft.com"}}
	metadata := &models.Metadata{
		Title:       "",
		Version:     "1.0",
		Maintainers: maintainers,
		Company:     "dummyCompany",
		Website:     "http://www.dummyCompany.com",
		Source:      "dummyCompany",
		License:     "GNU",
		Description: "Lorem ipsum",
	}
	testee := SimpleValidator{}
	err := testee.ValidateAndSanitize(metadata)
	assert.Error(t, err)
	aggregatedErr, k := err.(app.AggregatedValidationError)
	assert.True(t, k)
	assert.True(t, len(aggregatedErr.Errors()) == 2)
	log.Println(err.Error())
}
