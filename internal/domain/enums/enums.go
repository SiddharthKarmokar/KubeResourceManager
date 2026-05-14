package enums

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityModerate Severity = "moderate"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
	SeverityDegraded Severity = "degraded"
)

// String returns the string representation of the severity.
func (s Severity) String() string {
	return string(s)
}
