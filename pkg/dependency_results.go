package pkg

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ossf/scorecard/v4/pkg"

	sce "github.com/ossf/scorecard/v4/errors"
)

// ChangeType is the change type (added, updated, removed) of a dependency.
type ChangeType string

const (
	// Added suggests the dependency is a newly added one.
	Added ChangeType = "added"
	// Updated suggests the dependency is updated from an old version.
	Updated ChangeType = "updated"
	// Removed suggests the dependency is removed.
	Removed ChangeType = "removed"
)

// IsValid determines if a ChangeType is valid.
func (ct *ChangeType) IsValid() bool {
	switch *ct {
	case Added, Updated, Removed:
		return true
	default:
		return false
	}
}

// ScorecardResultsWithError is used for the dependency-diff module to record scorecard results and their errors.
type ScorecardResultsWithError struct {
	// ScorecardResults is the scorecard result for the dependency repo.
	ScorecardResults *pkg.ScorecardResult `json:"scorecardResults"`

	// Error is an error returned when running the scorecard checks. A nil Error indicates the run succeeded.
	Error error `json:"scorecardRunTimeError"`
}

// DependencyCheckResult is the dependency structure used in the returned results.
type DependencyCheckResult struct {
	// ChangeType indicates whether the dependency is added, updated, or removed.
	ChangeType *ChangeType `json:"changeType"`

	// Package URL is a short link for a package.
	PackageURL *string `json:"packageUrl"`

	// SourceRepository is the source repository URL of the dependency.
	SourceRepository *string `json:"sourceRepository"`

	// ManifestPath is the path of the manifest file of the dependency, such as go.mod for Go.
	ManifestPath *string `json:"manifestPath"`

	// Ecosystem is the name of the package management system, such as NPM, GO, PYPI.
	Ecosystem *string `json:"ecosystem"`

	// Version is the package version of the dependency.
	Version *string `json:"version"`

	// ScorecardResultsWithError is the scorecard checking results of the dependency.
	ScorecardResultsWithError ScorecardResultsWithError `json:"scorecardResultsWithError"`

	// Name is the name of the dependency.
	Name string `json:"name"`
}

// AsJSON for DependencyCheckResult exports the DependencyCheckResult as a JSON object.
func (dr *DependencyCheckResult) AsJSON(writer io.Writer) error {
	if err := json.NewEncoder(writer).Encode(*dr); err != nil {
		return sce.WithMessage(sce.ErrScorecardInternal, fmt.Sprintf("encoder.Encode: %v", err))
	}
	return nil
}
