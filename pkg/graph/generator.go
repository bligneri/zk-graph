package graph

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/bligneri/zk-graph/pkg/assets"
)

type Link struct {
	SourcePath string `json:"sourcePath"`
	TargetPath string `json:"targetPath"`
}

type LinkData struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  string `json:"value"`
}

type Note struct {
	AbsPath  string `json:"absPath"`
	Path     string `json:"path"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
}

type Node struct {
	ID      string `json:"id"`
	Group   int    `json:"group"`
	AbsPath string `json:"absPath"`
	Title   string `json:"title"`
}

type GraphData struct {
	Nodes []Node     `json:"nodes"`
	Links []LinkData `json:"links"`
}

func GenerateForceGraph(notes []Note, links []Link, highlight []string, outputFileName string, templateName string) error {
	nodes := []Node{}
	var templateContent string
	var tmpl *template.Template

	highlightMap := make(map[string]bool)
	for _, h := range highlight {
		highlightMap[h] = true
	}

	for _, note := range notes {
		title := note.Title
		if title == "" {
			title = note.Filename
		}

		group := 1
		if highlightMap[note.Filename] {
			group = 2
		}

		nodes = append(nodes, Node{
			ID:      title,
			Group:   group,
			AbsPath: note.AbsPath,
			Title:   title,
		})
	}

	pathToTitle := make(map[string]string)
	for _, note := range notes {
		if note.Path != "" {
			pathToTitle[note.Path] = note.Title
			if note.Title == "" {
				pathToTitle[note.Path] = note.Filename // Fallback to filename if title is empty
			}
		}
	}

	fmt.Println(pathToTitle)

	// Create the list of links
	var linkList []LinkData
	for _, link := range links {
		sourceID, sourceOK := pathToTitle[link.SourcePath]
		targetID, targetOK := pathToTitle[link.TargetPath]
		if sourceOK && targetOK {
			linkList = append(linkList, LinkData{
				Source: sourceID,
				Target: targetID,
				Value:  "2",
			})
		} else {
			// Improved error logging for missing links
			if !sourceOK {
				fmt.Printf("Warning: Missing source ID for link from %s\n", link.SourcePath)
			}
			if !targetOK {
				fmt.Printf("Warning: Missing target ID for link to %s\n", link.TargetPath)
			}
		}
	}

	data := GraphData{Nodes: nodes, Links: linkList}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling graph data: %v", err)
	}

	if templateName != "" {
		content, err := os.ReadFile(templateName)
		if err == nil {
			templateContent = string(content)
		} else {
			fmt.Printf("Warning: Unable to read specified template file '%s', using default template. Error: %v\n", templateName, err)
			templateContent = assets.GetForceGraphTemplate()
		}
	} else {
		templateContent = assets.GetForceGraphTemplate()
	}

	tmpl, err = template.New("graph").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	outputDir := filepath.Dir(outputFileName)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	outFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, map[string]string{"Data": string(dataJSON)})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fmt.Printf("Generated graph saved to %s\n", outputFileName)
	return nil
}
