package cli

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	slugRegex    = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
	issueIDRegex = regexp.MustCompile(`^\d+$`)
	eventIDRegex = regexp.MustCompile(`^[a-f0-9]{32}$`)
)

const maxSlugLength = 50

var specialEventIDs = map[string]bool{
	"latest":      true,
	"oldest":      true,
	"recommended": true,
}

func validateOrgSlug(slug string) error {
	if len(slug) > maxSlugLength {
		return NewInvalidInputError(fmt.Sprintf("organization slug too long (max %d chars): %s", maxSlugLength, slug))
	}
	if !slugRegex.MatchString(slug) {
		return NewInvalidInputError(fmt.Sprintf("invalid organization slug: %q (must match [a-z0-9][a-z0-9-]*)", slug))
	}
	return nil
}

func validateProjectSlug(slug string) error {
	if len(slug) > maxSlugLength {
		return NewInvalidInputError(fmt.Sprintf("project slug too long (max %d chars): %s", maxSlugLength, slug))
	}
	if !slugRegex.MatchString(slug) {
		return NewInvalidInputError(fmt.Sprintf("invalid project slug: %q (must match [a-z0-9][a-z0-9-]*)", slug))
	}
	return nil
}

func validateIssueID(id string) error {
	if !issueIDRegex.MatchString(id) {
		return NewInvalidInputError(fmt.Sprintf("invalid issue ID: %q (must be numeric)", id))
	}
	return nil
}

func validateEventID(id string) error {
	if specialEventIDs[id] {
		return nil
	}
	if !eventIDRegex.MatchString(id) {
		return NewInvalidInputError(fmt.Sprintf("invalid event ID: %q (must be 32 hex chars or latest/oldest/recommended)", id))
	}
	return nil
}

func validateInspectURL(rawURL string) error {
	if !strings.Contains(rawURL, "sentry.io") {
		return NewInvalidInputError(fmt.Sprintf("invalid Sentry URL: %q (must contain sentry.io)", rawURL))
	}
	return nil
}
