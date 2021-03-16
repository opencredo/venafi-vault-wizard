package pretty

import (
	"github.com/pterm/pterm"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
)

type report struct{}

func NewReport() reporter.Report {
	pterm.Error.ShowLineNumber = false
	return &report{}
}

func (r *report) AddSection(name string) reporter.Section {
	pterm.DefaultSection.Println(name)
	return &section{}
}

func (r *report) Finish(summary string, message string) {
	pterm.Println()
	pterm.DefaultBasicText.WithStyle(&pterm.Style{pterm.FgGreen}).Println(summary)
	pterm.DefaultHeader.Println(message)
}
