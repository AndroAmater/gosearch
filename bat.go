package main

import (
	"os/exec"

	"github.com/rivo/tview"
)

func bat(path string, lineNumber string) string {
	cmd := exec.Command(
		"bat",
		"--color=always",
		"--highlight-line",
		lineNumber,
		path,
	)
	stdout, err := cmd.Output()

	output := ""

	if err != nil {
		if err.Error() != "exit status 1" {
			panic(err)
		}
		output = tview.TranslateANSI("No results")
	} else {
		output = tview.TranslateANSI(string(stdout))
	}

	return output
}
