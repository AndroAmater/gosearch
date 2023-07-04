package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RgSearchResult struct {
	Type string `json:"type"`
	Data struct {
		Path struct {
			Text string `json:"text"`
		} `json:"path"`
		Lines struct {
			Text           string `json:"text"`
			LineNumber     int    `json:"line_number"`
			AbsoluteOffset int    `json:"absolute_offset"`
			Submatches     []struct {
				Match struct {
					Text string `json:"text"`
				} `json:"match"`
				Start int `json:"start"`
				End   int `json:"end"`
			} `json:"submatches"`
		} `json:"lines"`
	} `json:"data"`
}

type Match struct {
	ContextBefore []string
	LineNumber    int
	Text          string
	ContextAfter  []string
	Submatches    []struct {
		Match struct {
			Text string `json:"text"`
		} `json:"match"`
		Start int `json:"start"`
		End   int `json:"end"`
	}
}

type Folder struct {
	Path    string
	Matches []Match
}

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

	searchResultsRootNode := tview.NewTreeNode("")
	searchResults := tview.NewTreeView().
		SetRoot(searchResultsRootNode).
		SetCurrentNode(searchResultsRootNode)

	Search := func() {
		cmd := exec.Command(
			"rg",
			"-A",
			"2",
			"-B",
			"2",
			"--color",
			"ansi",
			"--line-number",
			"--json",
			searchInput.GetText(),
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

		previewText.Clear()
		previewText.SetDynamicColors(true)
		previewText.SetText(output)

		output = "[" + strings.TrimSuffix(strings.ReplaceAll(output, "\n", ","), ",") + "]"
		previewText.SetText(fmt.Sprintf("%+v\n", output))
		results := []RgSearchResult{}
		err = json.Unmarshal([]byte(output), &results)

		if err != nil {
			panic(err)
		}

		folders := []Folder{}

		for i, result := range results {
			if result.Type == "context" {
				continue
			}
			if result.Type == "begin" {
				folders = append(folders, Folder{
					Path:    result.Data.Path.Text,
					Matches: []Match{},
				})
			}
			if result.Type == "match" {
				contextBefore := []string{}
				if i > 1 {
					contextBefore = append(contextBefore, results[i-2].Data.Lines.Text)
				}
				if i > 0 {
					contextBefore = append(contextBefore, results[i-1].Data.Lines.Text)
				}

				contextAfter := []string{}
				if i < len(results)-2 {
					contextAfter = append(contextAfter, results[i+1].Data.Lines.Text)
				}
				if i < len(results)-1 {
					contextAfter = append(contextAfter, results[i+2].Data.Lines.Text)
				}

				if len(folders) == 0 {
					panic("Missing folder in rg result")
				}

				folders[len(folders)-1].Matches = append(folders[len(folders)-1].Matches, Match{
					ContextBefore: contextBefore,
					LineNumber:    result.Data.Lines.LineNumber,
					Text:          result.Data.Lines.Text,
					ContextAfter:  contextAfter,
					Submatches:    result.Data.Lines.Submatches,
				})
			}
		}

		searchResultsRootNode.ClearChildren()

		for _, folder := range folders {
			folderNode := tview.NewTreeNode(folder.Path)
			searchResultsRootNode.
				AddChild(folderNode)

			for _, match := range folder.Matches {
				matchNode := tview.NewTreeNode(match.Text)
				folderNode.AddChild(matchNode)
			}
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

	if err := app.SetRoot(main, true).EnableMouse(true).SetFocus(searchInput).Run(); err != nil {
		panic(err)
	}
}
