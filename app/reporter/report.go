package reporter

// Report represents all of the output to be reported to the user from a command
type Report interface {
	// AddSection adds a section to the report
	AddSection(name string) Section
	// Finish performs any final logging or saving of output at the end of the command
	Finish(summary string, message string)
}

// Section represents a logical part of a report
type Section interface {
	// AddCheck adds a check to the report section
	AddCheck(name string) Check
	// Info records some miscellaneous info discovered during a section
	Info(message string)
}

// Check represents an atomic check or task performed as part of a command
type Check interface {
	// UpdateStatus adds an intermediate update to a report Check
	UpdateStatus(status string)
	// Error finishes the Check with an unsuccessful message
	Error(message string)
	// Success finishes the Check with a successful message
	Success(message string)
	// Warning finishes the Check with a warning message
	Warning(message string)
}
