package pretty

import "github.com/pterm/pterm"

type check struct {
	spinner *pterm.SpinnerPrinter
}

func (c *check) UpdateStatus(status string) {
	c.spinner.UpdateText(status)
}

func (c *check) Error(message string) {
	c.spinner.Fail(message)
}

func (c *check) Success(message string) {
	c.spinner.Success(message)
}

func (c *check) Warning(message string) {
	c.spinner.Warning(message)
}
