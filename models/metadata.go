package models

type SearchResults struct {
	Results []Metadata `yaml:"results,flow"`
}

type Metadata struct {
	Title       string        `yaml:"title"`
	Version     string        `yaml:"version"`
	Company     string        `yaml:"company"`
	Website     string        `yaml:"website"`
	Source      string        `yaml:"source"`
	License     string        `yaml:"license"`
	Description string        `yaml:"description"`
	Maintainers []*Maintainer `yaml:"maintainers,flow"`
}

type Maintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type Response struct {
	Id string `yaml:"id"`
}

type ErrorResponse struct {
	Message string `yaml:"message"`
}
