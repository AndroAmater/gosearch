package main

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
)

type RgSearchResult struct {
	Type string `json:"type"`
	Data struct {
		Path struct {
			Text string `json:"text"`
		} `json:"path"`
		Lines struct {
			Text string `json:"text"`
		} `json:"lines"`
		LineNumber     int `json:"line_number"`
		AbsoluteOffset int `json:"absolute_offset"`
		Submatches     []struct {
			Match struct {
				Text string `json:"text"`
			} `json:"match"`
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"submatches"`
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

func rg(pattern string) []Folder {
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
		pattern,
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

	output = "[" + strings.TrimSuffix(strings.ReplaceAll(output, "\n", ","), ",") + "]"
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
				LineNumber:    result.Data.LineNumber,
				Text:          result.Data.Lines.Text,
				ContextAfter:  contextAfter,
				Submatches:    result.Data.Submatches,
			})
		}
	}

	return folders
}
