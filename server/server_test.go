package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/erikwu09/yamlr/models"
	"gopkg.in/yaml.v2"
)

var data = `
company: "Random Inc."
description: aasasdfadfs
license: Apache-2.0
maintainers: 
  - 
    email: firstmaintainer@hotmail.com
    name: "firstmaintainer app1"
  - 
    email: secondmaintainer@gmail.com
    name: "secondmaintainer app1"
source: "https://github.com/random/repo"
title: "Valid App 1"
version: "0.0.1"
website: "https://website.com"
`

func Test_YamlUnmarshal(t *testing.T) {
	metadata := models.Metadata{}
	err := yaml.Unmarshal([]byte(data), &metadata)
	assert.NoError(t, err)
	assert.True(t, len(metadata.Maintainers) == 2)
	assert.True(t, metadata.Maintainers[0].Name == "firstmaintainer app1")
}
