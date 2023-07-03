package main

import (
	"os/exec"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	searchButton := tview.NewButton("Search")

	main := tview.NewFlex()

	sidebar := tview.NewFlex().
		SetDirection(tview.FlexRow)

	preview := tview.NewFlex()
	previewText := tview.NewTextView()
	preview.AddItem(previewText, 0, 1, false)

	searchInput := tview.NewInputField()

	searchResults := tview.NewTable()

	searchButton.SetSelectedFunc(func() {
		cmd := exec.Command("rg", "-A", "2", "-B", "2", "--color", "ansi", searchInput.GetText())
		out, err := cmd.Output()
		if err != nil {
			panic(err)
		}

		previewText.Clear()
		previewText.SetDynamicColors(true)
		previewText.SetText(tview.TranslateANSI(string(out)))
	})

	searchInput.SetBorder(true).SetRect(0, 0, 5, 1)
	searchButton.SetBorder(true).SetRect(0, 0, 5, 3)
	searchResults.SetBorder(true).SetRect(0, 0, 0, 0)

	main.
		AddItem(sidebar, 50, 0, false).
		AddItem(preview, 0, 1, false).
		SetDirection(tview.FlexColumn)

	sidebar.
		AddItem(searchInput, 3, 0, false).
		AddItem(searchButton, 3, 0, false).
		AddItem(searchResults, 0, 1, false).
		SetBorder(true).
		SetTitle("Search")

	preview.
		SetBorder(true).
		SetTitle("Preview")

	app.SetFocus(searchInput)

	if err := app.SetRoot(main, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
