package report

import (
	"fmt"
	"time"
)

// Status Enum defined for PolicyReport
type Status = string

// Severity Enum defined for PolicyReport
type Severity = string

// Enums for predefined values from the PolicyReport spec
const (
	Fail  Status = "fail"
	Warn  Status = "warn"
	Error Status = "error"
	Pass  Status = "pass"
	Skip  Status = "skip"

	Low    Severity = "low"
	Medium Severity = "medium"
	Heigh  Severity = "heigh"

	defaultString = ""
	debugString   = "debug"
	infoString    = "info"
	warningString = "warning"
	errorString   = "error"
)

// Internal Priority definitions and weighting
const (
	DefaultPriority = iota
	DebugPriority
	InfoPriority
	WarningPriority
	ErrorPriority
)

// Priority Enum for internal Result weighting
type Priority int

// String maps the internal weighting of Priorities to a String representation
func (p Priority) String() string {
	switch p {
	case DebugPriority:
		return debugString
	case InfoPriority:
		return infoString
	case WarningPriority:
		return warningString
	case ErrorPriority:
		return errorString
	default:
		return defaultString
	}
}

// PriorityFromStatus creates a Priority based on a Status
func PriorityFromStatus(p Status) Priority {
	switch p {
	case Fail:
		return ErrorPriority
	case Error:
		return ErrorPriority
	case Warn:
		return WarningPriority
	case Pass:
		return InfoPriority
	default:
		return DefaultPriority
	}
}

// NewPriority creates a new Priority based an its string representation
func NewPriority(p string) Priority {
	switch p {
	case debugString:
		return DebugPriority
	case infoString:
		return InfoPriority
	case warningString:
		return WarningPriority
	case errorString:
		return ErrorPriority
	default:
		return DefaultPriority
	}
}

// Resource from the Kubernetes spec k8s.io/api/core/v1.ObjectReference
type Resource struct {
	APIVersion string
	Kind       string
	Name       string
	Namespace  string
	UID        string
}

// Result from the PolicyReport spec wgpolicyk8s.io/v1alpha1.PolicyReportResult
type Result struct {
	Message   string
	Policy    string
	Rule      string
	Priority  Priority
	Status    Status
	Severity  Severity
	Category  string
	Scored    bool
	Resources []Resource
}

// GetIdentifier returns a global unique Result identifier
func (r Result) GetIdentifier() string {
	suffix := ""
	if len(r.Resources) > 0 {
		suffix = "__" + r.Resources[0].UID
	}

	return fmt.Sprintf("%s__%s__%s%s", r.Policy, r.Rule, r.Status, suffix)
}

// Summary from the PolicyReport spec wgpolicyk8s.io/v1alpha1.PolicyReportSummary
type Summary struct {
	Pass  int
	Skip  int
	Warn  int
	Error int
	Fail  int
}

// PolicyReport from the PolicyReport spec wgpolicyk8s.io/v1alpha1.PolicyReport
type PolicyReport struct {
	Name              string
	Namespace         string
	Results           map[string]Result
	Summary           Summary
	CreationTimestamp time.Time
}

// GetIdentifier returns a global unique PolicyReport identifier
func (pr PolicyReport) GetIdentifier() string {
	return fmt.Sprintf("%s__%s", pr.Namespace, pr.Name)
}

// GetNewResults filters already existing Results from the old PolicyReport and returns only the diff with new Results
func (pr PolicyReport) GetNewResults(or PolicyReport) []Result {
	diff := make([]Result, 0)

	for _, r := range pr.Results {
		if _, ok := or.Results[r.GetIdentifier()]; ok {
			continue
		}

		diff = append(diff, r)
	}

	return diff
}

// ClusterPolicyReport from the PolicyReport spec wgpolicyk8s.io/v1alpha1.ClusterPolicyReport
type ClusterPolicyReport struct {
	Name              string
	Results           map[string]Result
	Summary           Summary
	CreationTimestamp time.Time
}

// GetIdentifier returns a global unique ClusterPolicyReport identifier
func (cr ClusterPolicyReport) GetIdentifier() string {
	return cr.Name
}

// GetNewResults filters already existing Results from the old PolicyReport and returns only the diff with new Results
func (cr ClusterPolicyReport) GetNewResults(cor ClusterPolicyReport) []Result {
	diff := make([]Result, 0)

	for _, r := range cr.Results {
		if _, ok := cor.Results[r.GetIdentifier()]; ok {
			continue
		}

		diff = append(diff, r)
	}

	return diff
}
