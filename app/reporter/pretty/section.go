package pretty

import (
	"github.com/pterm/pterm"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
)

type section struct{}

func (s *section) AddCheck(name string) reporter.Check {
	spinner, _ := pterm.DefaultSpinner.Start(name)
	return &check{spinner: spinner}
}

func (s *section) Info(message string) {
	pterm.Println()
	pterm.Println(message)
}
