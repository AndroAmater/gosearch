package main

import (
	"flag"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateMatchNodeSelectHandler(
	path string,
	lineNumber string,
	preview *tview.Flex,
	previewText *tview.TextView,
	sidebarOnly bool,
) func() {
	return func() {
		if !sidebarOnly {
			output := bat(path, lineNumber)
			previewText.Clear()
			previewText.SetDynamicColors(true)
			previewText.SetText(output)
			preview.SetTitle(path)
		}
	}
}

func main() {
	sidebarOnly := flag.Bool("sidebar-only", false, "Only show the sidebar")
	flag.Parse()

	app := tview.NewApplication()
	searchButton := tview.NewButton("Search")

	main := tview.NewFlex()

	sidebar := tview.NewFlex().
		SetDirection(tview.FlexRow)

	var preview *tview.Flex
	var previewText *tview.TextView

	if !*sidebarOnly {
		preview = tview.NewFlex()
		previewText = tview.NewTextView()
		preview.AddItem(previewText, 0, 1, false)
	}

	searchInput := tview.NewInputField()

	searchResultsRootNode := tview.NewTreeNode("")
	searchResults := tview.NewTreeView().
		SetRoot(searchResultsRootNode).
		SetCurrentNode(searchResultsRootNode)

	Search := func() {
		folders := rg(searchInput.GetText())

		searchResultsRootNode.ClearChildren()
		var firstNode *tview.TreeNode
		var firstNodeSelectHandler func()

		for i, folder := range folders {
			folderNode := tview.NewTreeNode(folder.Path)
			searchResultsRootNode.
				AddChild(folderNode)

			for k, match := range folder.Matches {
				matchNode := tview.NewTreeNode(match.Text).
					SetSelectedFunc(
						CreateMatchNodeSelectHandler(
							folder.Path,
							strconv.Itoa(match.LineNumber),
							preview,
							previewText,
							*sidebarOnly,
						))

				folderNode.AddChild(matchNode)

				if i == 0 && k == 0 {
					firstNode = matchNode
					firstNodeSelectHandler = CreateMatchNodeSelectHandler(
						folder.Path,
						strconv.Itoa(match.LineNumber),
						preview,
						previewText,
						*sidebarOnly,
					)
				}
			}
		}

		if firstNode != nil {
			app.SetFocus(searchResults)
			searchResults.SetCurrentNode(firstNode)
			firstNodeSelectHandler()
		}
	}

	searchButton.SetSelectedFunc(Search)

	searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			Search()
		}
		return event
	})

	searchInput.SetBorder(true).SetRect(0, 0, 5, 1)
	searchButton.SetBorder(true).SetRect(0, 0, 5, 3)
	searchResults.SetBorder(true).SetRect(0, 0, 0, 0)

	main.AddItem(sidebar, 50, 0, false)
	if !*sidebarOnly {
		main.AddItem(preview, 0, 1, false)
	}
	main.SetDirection(tview.FlexColumn)

	sidebar.
		AddItem(searchInput, 3, 0, false).
		AddItem(searchButton, 3, 0, false).
		AddItem(searchResults, 0, 1, false).
		SetBorder(true).
		SetTitle("Search")

	if !*sidebarOnly {
		preview.
			SetBorder(true).
			SetTitle("Preview")
	}

	if err := app.SetRoot(main, true).EnableMouse(true).SetFocus(searchInput).Run(); err != nil {
		panic(err)
	}
}
