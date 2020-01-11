package validation

import (
	"net/mail"
	"net/url"
	_ "reflect"

	"gitlab.com/erikwu09/yamlr/app"
	"gitlab.com/erikwu09/yamlr/models"
)

// Simple, brute-force validation
type SimpleValidator struct {
}

func (v SimpleValidator) SanitizeURLs(metadata *models.Metadata) error {
	aggregateErr := app.NewAggregatedValidationError()
	if metadata.Website != "" {
		websiteURL, err := url.Parse(metadata.Website)
		if err != nil {
			aggregateErr.AddError(app.ValidationError{Path: "Website", Reason: "url not properly formed"})
		} else {
			metadata.Website = websiteURL.String()
		}
	}
	if metadata.Source != "" {
		sourceURL, err := url.Parse(metadata.Source)
		if err != nil {
			aggregateErr.AddError(app.ValidationError{Path: "Source", Reason: "url not properly formed"})
		} else {
			metadata.Source = sourceURL.String()
		}
	}
	if metadata.Maintainers != nil {
		for _, m := range metadata.Maintainers {
			email, err := mail.ParseAddress(m.Email)
			if err != nil {
				aggregateErr.AddError(app.ValidationError{Path: "Maintainers[].Email", Reason: "email address not properly formed"})
			} else {
				m.Email = email.String()
			}
		}
	}
	if len(aggregateErr.Errors()) == 0 {
		return nil
	}
	return aggregateErr
}

func (v SimpleValidator) ValidateAndSanitize(metadata *models.Metadata) error {
	aggregateErr := app.NewAggregatedValidationError()
	if metadata.Title == "" {
		aggregateErr.AddError(app.ValidationError{Path: "Title", Reason: "empty value"})
	}
	// TODO: check for version semantics
	if metadata.Version == "" {
		aggregateErr.AddError(app.ValidationError{Path: "Version", Reason: "empty value"})
	}
	if metadata.Company == "" {
		aggregateErr.AddError(app.ValidationError{Path: "Company", Reason: "empty value"})
	}
	if metadata.Website == "" {
		aggregateErr.AddError(app.ValidationError{Path: "Website", Reason: "empty value"})
	} else {
		websiteURL, err := url.Parse(metadata.Website)
		if err != nil {
			aggregateErr.AddError(app.ValidationError{Path: "Website", Reason: "url not properly formed"})
		} else {
			metadata.Website = websiteURL.String()
		}
	}
	if metadata.Source == "" {
		aggregateErr.AddError(app.ValidationError{Path: "Source", Reason: "empty value"})
	} else {
		sourceURL, err := url.Parse(metadata.Source)
		if err != nil {
			aggregateErr.AddError(app.ValidationError{Path: "Source", Reason: "url not properly formed"})
		} else {
			metadata.Source = sourceURL.String()
		}
	}
	if metadata.License == "" {
		aggregateErr.AddError(app.ValidationError{Path: "License", Reason: "empty value"})
	}
	if metadata.Description == "" {
		aggregateErr.AddError(app.ValidationError{Path: "Description", Reason: "empty value"})
	}
	if metadata.Maintainers == nil || len(metadata.Maintainers) == 0 {
		aggregateErr.AddError(app.ValidationError{Path: "Maintainers", Reason: "empty value"})
	} else {
		for _, m := range metadata.Maintainers {
			if m.Email == "" {
				aggregateErr.AddError(app.ValidationError{Path: "Maintainers[].Email", Reason: "empty value"})
			} else {
				email, err := mail.ParseAddress(m.Email)
				if err != nil {
					aggregateErr.AddError(app.ValidationError{Path: "Maintainers[].Email", Reason: "email address not properly formed"})
				} else {
					m.Email = email.String()
				}
			}
			if m.Name == "" {
				aggregateErr.AddError(app.ValidationError{Path: "Maintainers[].Email", Reason: "empty value"})
			}
		}
	}
	if len(aggregateErr.Errors()) == 0 {
		return nil
	}
	return aggregateErr
}
