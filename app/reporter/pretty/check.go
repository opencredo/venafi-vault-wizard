package pretty

import (
	"fmt"

	"github.com/pterm/pterm"
)

type check struct {
	spinner *pterm.SpinnerPrinter
}

func (c *check) UpdateStatus(status string) {
	c.spinner.UpdateText(status)
}

func (c *check) UpdateStatusf(status string, a ...interface{}) {
	c.UpdateStatus(fmt.Sprintf(status, a...))
}

func (c *check) Error(message string) {
	c.spinner.Fail(message)
}

func (c *check) Errorf(status string, a ...interface{}) {
	c.Error(fmt.Sprintf(status, a...))
}

func (c *check) Success(message string) {
	c.spinner.Success(message)
}

func (c *check) Successf(status string, a ...interface{}) {
	c.Success(fmt.Sprintf(status, a...))
}

func (c *check) Warning(message string) {
	c.spinner.Warning(message)
}

func (c *check) Warningf(status string, a ...interface{}) {
	c.Warning(fmt.Sprintf(status, a...))
}
