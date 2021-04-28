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
	// UpdateStatusf is like UpdateStatus but takes a format string and optional args in Printf style
	UpdateStatusf(status string, a ...interface{})
	// Error finishes the Check with an unsuccessful message
	Error(message string)
	// Errorf is like Error but takes a format string and optional args in Printf style
	Errorf(status string, a ...interface{})
	// Success finishes the Check with a successful message
	Success(message string)
	// Successf is like Success but takes a format string and optional args in Printf style
	Successf(status string, a ...interface{})
	// Warning finishes the Check with a warning message
	Warning(message string)
	// Warningf is like Warning but takes a format string and optional args in Printf style
	Warningf(status string, a ...interface{})
}
